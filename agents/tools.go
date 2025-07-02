package agents

import (
	"context"
	"log"
	//"time"
	"whatsapp-bot/db"
	"sync"
	"runtime"
	//"whatsapp-bot/wa"
	//"whatsapp-bot/wa/handlers"
)

func ProcessMessageByAI (ctx context.Context) error {
	
	var resp *AgentHouseResponse
	msgs, err := db.GetConvoMessagesUnProcessed(ctx)

	if err != nil {
		log.Printf("failed to get message by ID: %v", err)
	}

	if len(msgs) == 0{
		log.Println("No unprocessed messages found")
		return nil
	}

	workerN := runtime.NumCPU()
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
					log.Printf("AI worker %d panic: %v", id, r)
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
					log.Printf("AI worker %d processing message: %s", id, msg.Text)
					resp, err = AIProcessHouseMessage(msg.Text)
					if err != nil {
						log.Printf("failed to process house message for msg ID %d: %v", msg.ID, err)
						continue
					}

					log.Printf("Processed message: %s", resp.AiMessage)
					err = db.UpdateConvoMessageAIPRocessedByID(ctx, msg.ID)
					if err != nil {
						log.Printf("failed to update convo message for msg ID %d: %v", msg.ID, err)
						continue
					}

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
			}
		}(id)
	}

	wg.Wait()
	return err

}