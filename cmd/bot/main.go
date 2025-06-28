package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	//"time"

	//"go.mau.fi/whatsmeow/types/events"

	"whatsapp-bot/wa"

	"whatsapp-bot/wa/handlers"
)


func main() {
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := wa.Connect(ctx, handlers.HandleEvent) 

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer client.Disconnect()

	log.Println("connected to WhatsApp")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Println("shutting down")
}