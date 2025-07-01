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

	"whatsapp-bot/agents" // This package contains the AI processing logic, like ProcessHouseMessage and ProcessMessageByAI
)


func main() {
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initializing the database before connecting to wa
	if err := db.InitDB(ctx); err != nil { log.Fatalf("failed to initialize database: %v", err) }

	// Reading config file to get the configurations
	config, err := utils.ReadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	//Initilazing the connection to whatsapp
	client, err := wa.Connect(ctx, func(evt interface{}) {
		handlers.HandleEvent(evt, config) // A closure that handles whatsapp events as defined by config
	})

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer client.Disconnect()

	log.Println("connected to WhatsApp")

	// This is how we send a message to a whatsapp chat, will be used by the bot to send messages
	// err = handlers.SendText(
	// 	ctx,
	// 	client,
	// 	config.Whatsapp.WhiteListedChats[0],
	// 	"Hello from WhatsApp Bot!",
	// )

	// if err != nil {
	// 	log.Println("failed to send message: %v", err)
	// }

	//Process message by AI
	if err := agents.ProcessMessageByAI(ctx); err != nil {
		log.Printf("failed to process message by AI: %v", err)
	}

	
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Println("shutting down")
}