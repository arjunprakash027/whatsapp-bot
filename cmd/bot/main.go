package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	//"time"

	//"go.mau.fi/whatsmeow/types/events"

	"whatsapp-bot/wa" //Has the connection functionality to WhatsApp

	"whatsapp-bot/wa/handlers" //HandleEvent handler is imported from here

	"whatsapp-bot/db"
)


func main() {
	
	

	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initializing the database before connecting to wa
	if err := db.InitDB(); err != nil { log.Fatalf("failed to initialize database: %v", err) }
	if err := db.PrepareConvoInsertStatement(ctx); err != nil {log.Fatalf("failed to prepare conversation insert statement: %v", err) }

	//Initilazing the connection to whatsapp
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