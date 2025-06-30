package handlers

import (
	"log"
	"context"
	"whatsapp-bot/db"
	"whatsapp-bot/utils"
	"go.mau.fi/whatsmeow/types/events"
	"time"
)

var ctx = context.Background()

// This is a event handler that performs certain action on every event recived by the WhatsApp client.
func HandleEvent(
	evt interface{},
	cfg *utils.Config, // This is the config that is read from the config.yaml file
) {

	switch v := evt.(type) {
	case *events.Message:
		if txt := v.Message.GetConversation(); txt != "" {

			log.Printf("[%s] [%s] %s",
				v.Info.Timestamp.Format(time.RFC3339),
				v.Info.Chat.String(),
				txt,
			)
			
			// save the convo to database
			// After serious thought, I decided linear scan is the best way, everything else is just overkill
			// Check if there are any whitelisted chats in the config, if yes only save them, else save all messages - this is primarily done to avoid unncessary database writes and space
			if len(cfg.Whatsapp.WhiteListedChats) > 0 {
				for _, chat := range cfg.Whatsapp.WhiteListedChats {
					if chat == v.Info.Chat.String() {
						StoreConvo(v, txt, "live-message")
						log.Printf("Message saved for whitelisted chat: %s", chat)
					}
				}
			} else {
				StoreConvo(v,txt,"live-message")
			}
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

	const read_by_ai = 0
	err := db.SaveConvoMessage(
		ctx,
		v.Info.ID,
		v.Info.Chat.String(),
		v.Info.Sender.String(),
		txt,
		channel,
		v.Info.Timestamp,
		read_by_ai,
	)

	if err != nil {
		log.Printf("Error saving message: %v", err)
	} else {
		log.Printf("Message saved: %s", v.Info.ID)
	}
}