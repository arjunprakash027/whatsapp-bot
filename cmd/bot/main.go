package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mau.fi/whatsmeow/types/events"

	"whatsapp-bot/wa"
)

// This is a event handler that performs certain action on every event recived by the WhatsApp client.
func handleEvent(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if txt := v.Message.GetConversation(); txt != "" {
			log.Printf("[%s] %s",
				v.Info.Timestamp.Format(time.RFC3339),
				txt,
			)
		}
	}
}


func main() {
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := wa.Connect(ctx, handleEvent) 

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer client.Disconnect()

	log.Println("connected to WhatsApp")

	// Block until the user hits Ctrl‑C.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Println("shutting down")
}