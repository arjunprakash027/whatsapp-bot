package agents

import (
	"context"
	"log"
	"sync"
	"time"
	"whatsapp-bot/db"
	"whatsapp-bot/utils"
	"whatsapp-bot/wa/handlers"
	"go.mau.fi/whatsmeow"
	//"whatsapp-bot/wa"
	//"whatsapp-bot/wa/handlers"
)

func ProcessMessageByAIPoller(
	ctx context.Context,
	Config *utils.Config,
) {

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("AI Poller stopped")
			return
		case <-ticker.C:
			if err := ProcessBatchAI(ctx, Config.AI.Controls.WorkerCount, Config.AI.BenchmarkMessage); err != nil {
				log.Printf("failed to process batch AI: %v", err)
			}
		}
	}
}

func ProcessBatchAI(ctx context.Context, workerN int, BenchmarkMessage string) error {

	var resp *AgentHouseResponse
	msgs, err := db.GetConvoMessagesUnProcessed(ctx)

	if err != nil {
		log.Printf("failed to get message by ID: %v", err)
	}

	if len(msgs) == 0 {
		log.Println("No unprocessed messages found by polling")
		return nil
	}

	jobsChan := make(chan db.Message, len(msgs))
	var wg sync.WaitGroup
	wg.Add(workerN)

	for _, msg := range msgs {
		jobsChan <- msg
	}
	close(jobsChan)

	for i := 0; i < workerN; i++ {
		id := i

		go func(id int) {
			defer wg.Done()

			defer func() {
				if r := recover(); r != nil {
					log.Printf("AI worker %d panic: %v", r)
				}
			}()

			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-jobsChan:
					if !ok {
						log.Printf("AI worker %d: no more jobs", id)
						return
					}

					err = db.UpdateConvoMessageAIPRocessedByID(ctx, msg.ID, 1)
					if err != nil {
						log.Printf("failed to update convo message for msg ID %d: %v", msg.ID, err)
						continue
					}

					log.Printf("AI worker %d processing message: %s", id, msg.Text)
					resp, err = AIProcessHouseMessage(msg.Text, BenchmarkMessage)

					if resp.AgreedToProcess {
						log.Printf("AI worker %d response: %+v", id, resp)

						if err != nil {
							log.Printf("failed to process house message for msg ID %d: %v", msg.ID, err)
							err = db.UpdateConvoMessageAIPRocessedByID(ctx, msg.ID, 0)
							if err != nil {
								log.Printf("failed to update convo message for msg ID %d: %v", msg.ID, err)
							}
							continue
						} else {
							log.Printf("AI worker %d processed message: %s", id, resp.AiMessage)
							err = db.SaveProcessedMessage(
								ctx,
								msg.ID,
								msg.ChatJID,
								msg.SenderJID,
								msg.Text,
								resp.AiAddress,
								resp.AiPrimaryContact,
								resp.AiSecondaryContact,
								resp.AiMessage,
								0, // no message is sent yet
							)
							if err != nil {
								log.Printf("failed to save processed message for msg ID %d: %v", msg.ID, err)
							}
						}
					} else {
						log.Printf("AI worker %d refused to process the message due to high edit distance", id)
					}
				}
			}
		}(id)
	}

	wg.Wait()
	return err
}

func DispatcherTool(client *whatsmeow.Client ,ctx context.Context, config *utils.Config) {
	
	batchComplete := make(chan struct{}, 1)

	batchComplete <- struct{}{}

	for {
		select {
		case <-ctx.Done():
			log.Println("Dispatcher Poller stopped")
			return
		
		case <-batchComplete:
			unsentMessages, err := db.GetUnsentProcessedMessages(ctx)
			if err != nil {
				log.Fatalf("Error while retreiving unsent processed messages")
			}

			if len(unsentMessages) == 0 {
				time.Sleep(20 * time.Second)
				batchComplete <- struct{}{}
				continue
			}

			go ProcessBatchDispatch(
				client,
				ctx,
				config,
				unsentMessages,
				batchComplete,
			)
		}

	}
}

func ProcessBatchDispatch(
	client *whatsmeow.Client, 
	ctx context.Context, 
	config *utils.Config, 
	messages []db.ProcessedMessage, 
	completion chan<- struct{},
) {
	defer func() {
		completion <- struct{}{}
	}()
	
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for _, msg := range messages {

		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			log.Println(msg)
			err := handlers.SendText(
				ctx,
				client,
				config.Whatsapp.MessageReceiver,
				msg.AIMessage,
			)

			if err != nil {
				log.Println("failed to send message:", err)
			}
		}
	}
}
