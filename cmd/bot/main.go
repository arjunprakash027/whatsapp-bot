package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"runtime"
	"sync"
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

	//Setup a event channel to unbloack whatsapp events 
	eventChan := make(chan interface{}, 500)

	//Initilazing the connection to whatsapp
	client, err := wa.Connect(ctx, func(evt interface{}) {
		select {
		case eventChan <- evt:
		default:handlers.HandleEvent(evt, config)
			log.Println("event channel is full, dropping event to avoid blocking")	
		} // A closure that handles whatsapp events as defined by config
	})


	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer client.Disconnect()

	workerN := runtime.NumCPU()
	wg := startWorkers(ctx, workerN, eventChan, config)

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

	close(eventChan)
	wg.Wait() // Wait for all workers to finish
	log.Println("shut down complete")
}

func startWorkers(
	ctx context.Context,
	workerN int,
	eventChan <-chan interface{},
	cfg *utils.Config,
) *sync.WaitGroup {

	var wg sync.WaitGroup
	wg.Add(workerN)

	for i := 0; i < workerN; i++ {
		id := i

		go func(id int) {
			defer wg.Done()

			defer func() {
				if r := recover(); r != nil {
					log.Printf("worker %d panic: %v", id, r)
				}
			}()

			for {
				select {
				case <-ctx.Done():
					return
				case evt, ok := <-eventChan:
					if !ok {
						log.Printf("worker %d: event channel closed", id)
						return
					}
					handlers.HandleEvent(evt, cfg)
				}
			}
		}(id)

	}

	return &wg
}