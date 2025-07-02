package handlers

import (
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"log"
)

func HistoryHandler(history *events.HistorySync) {

	for _, conv := range history.Data.Conversations {

		if conv.ID == nil {
			log.Printf("Skipping history sync for conversation with no ID")
			continue
		}

		jid, err := types.ParseJID(*conv.ID)
		if err != nil {
			log.Printf("Error parsing JID for conversation id %s: %v", conv.ID, err)
			continue
		}

		log.Printf("History sync for conversation id %s",
			jid)
	}
}
