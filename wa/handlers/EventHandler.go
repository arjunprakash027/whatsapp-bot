package handlers

import (
	"log"
	"context"
	"whatsapp-bot/db"
	"go.mau.fi/whatsmeow/types/events"
	"time"
)

var ctx = context.Background()

// This is a event handler that performs certain action on every event recived by the WhatsApp client.
func HandleEvent(evt interface{}) {

	switch v := evt.(type) {
	case *events.Message:
		if txt := v.Message.GetConversation(); txt != "" {
			log.Printf("[%s] [%s] %s",
				v.Info.Timestamp.Format(time.RFC3339),
				v.Info.Chat.String(),
				txt,
			)
			
			// save the convo to database
			StoreConvo(v,txt,"live-message")
		}
		
	case *events.HistorySync:
		log.Printf("History sync event started")
		HistoryHandler(v)
	}
}

func StoreConvo(
	v *events.Message,
	txt, channel string,
) {
	err := db.SaveConvoMessage(
		ctx,
		v.Info.ID,
		v.Info.Chat.String(),
		v.Info.Sender.String(),
		txt,
		channel,
		v.Info.Timestamp,
	)

	if err != nil {
		log.Printf("Error saving message: %v", err)
	} else {
		log.Printf("Message saved: %s", v.Info.ID)
	}
}