package agents

import (
	"context"
	"log"
	"sync"
	"time"
	"whatsapp-bot/db"
	"whatsapp-bot/utils"
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
			if err := ProcessBatchAI(ctx, Config.AI.Controls.WorkerCount); err != nil {
				log.Printf("failed to process batch AI: %v", err)
			}
		}
	}
}

func ProcessBatchAI(ctx context.Context, workerN int) error {

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

					err = db.UpdateConvoMessageAIPRocessedByID(ctx, msg.ID, 0)
					if err != nil {
						log.Printf("failed to update convo message for msg ID %d: %v", msg.ID, err)
						continue
					}

					log.Printf("AI worker %d processing message: %s", id, msg.Text)
					resp, err = AIProcessHouseMessage(msg.Text)

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
					}					

					// err = db.SaveProcessedMessage(
					// 	ctx,
					// 	msg.ID,
					// 	msg.ChatJID,
					// 	msg.SenderJID,
					// 	msg.Text,
					// 	resp.AiAddress,
					// 	resp.AiPrimaryContact,
					// 	resp.AiSecondaryContact,
					// 	resp.AiMessage,
					// 	0, // no message is sent yet
					// )
					// if err != nil {
					// 	log.Printf("failed to save processed message for msg ID %d: %v", msg.ID, err)
					// }
				}
			}
		}(id)
	}

	wg.Wait()
	return err
}
