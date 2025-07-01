package agents

import (
	"context"
	"log"
	//"time"
	"whatsapp-bot/db"
	//"whatsapp-bot/wa"
	//"whatsapp-bot/wa/handlers"
)

func ProcessMessageByAI (ctx context.Context) error {
	
	var resp *AgentHouseResponse
	msgs, err := db.GetConvoMessagesUnProcessed(ctx)

	if err != nil {
		log.Printf("failed to get message by ID: %v", err)
	}

	for _, msg := range msgs {
		log.Printf("Incoming message read from db: [%s] %s", msg.ChatJID, msg.Text)
		resp, err = AIProcessHouseMessage(msg.Text)

		if err != nil {
			log.Printf("failed to process house message: %v", err)
			continue
		}

		log.Printf("Processed message: %s", resp.AiMessage)
		err = db.UpdateConvoMessageAIPRocessedByID(
			ctx,
			msg.ID,
		)

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
		 	0, //no message is sent yet
		)

		if err != nil {
			log.Printf("failed to save processed message: %v", err)
		}
	}

	return err
	
}