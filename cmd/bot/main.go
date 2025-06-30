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

	"whatsapp-bot/utils" //This is a utility package that has some helper functions, like config reader and so on
)


func main() {
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initializing the database before connecting to wa
	if err := db.InitDB(); err != nil { log.Fatalf("failed to initialize database: %v", err) }
	if err := db.PrepareConvoInsertStatement(ctx); err != nil {log.Fatalf("failed to prepare conversation insert statement: %v", err) }

	// Reading config file to get the configurations

	config, err := utils.ReadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}
	log.Printf("whitelisted chats: %v", config.Whatsapp.WhiteListedChats)

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