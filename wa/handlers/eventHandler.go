package handlers

import (
	"log"
	"go.mau.fi/whatsmeow/types/events"
	"time"
)

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
		}
		
	case *events.HistorySync:
		log.Printf("History sync event started")
		HistoryHandler(v)
	}
}
