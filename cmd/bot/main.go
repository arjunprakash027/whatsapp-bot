package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"flag"
	//"time"

	//"go.mau.fi/whatsmeow/types/events"
	"go.mau.fi/whatsmeow"

	"whatsapp-bot/wa" //Has the connection functionality to WhatsApp

	"whatsapp-bot/wa/handlers" //HandleEvent handler is imported from here

	"whatsapp-bot/db"

	"whatsapp-bot/utils" //This is a utility package that has some helper functions, like config reader and so on

	"whatsapp-bot/agents" // This package contains the AI processing logic, like ProcessHouseMessage and ProcessMessageByAI
)

func main() {

	mode := flag.String("mode", "all", "Execution mode: collect, process, dispatch, all")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initializing the database before connecting to wa
	if err := db.InitDB(ctx); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// Reading config file to get the configurations
	config, err := utils.ReadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	switch *mode {
	case "collect":
		go func () {
			<-sig
			log.Println("Shutdown Signal Received!, cancelling context")
			cancel()
		} ()

		client, eventChan, _ := setupWhatsmeowClient(ctx, config)
		runCollectMode(client, eventChan, ctx, config)

	case "process":
		go func () {
			<-sig
			log.Println("Shutdown Signal Received!, cancelling context")
			cancel()
		} ()
		runProcessMode(ctx, config)
	
	case "dispatch":
		go func() {
			<-sig
			log.Println("Shutdown Signal Received!, cancelling context")
			cancel()
		} ()
		client ,_ ,_ := setupWhatsmeowClient(ctx, config)
		runDispatchMode(client, ctx, config)

	case "all":
		var wg sync.WaitGroup
		wg.Add(3)

		client, eventChan, _ := setupWhatsmeowClient(ctx, config)

		go func() {
			defer wg.Done()
			runCollectMode(client, eventChan, ctx, config)
		}()
		
		go func() {
			defer wg.Done()
			runProcessMode(ctx, config)
		}()

		go func() {
			defer wg.Done()
			runDispatchMode(client ,ctx, config)
		}()

		go func () {
			<-sig
			log.Println("Shutdown Signal Received!, cancelling context")
			cancel()
		}()

		wg.Wait()

	default:
		log.Printf("Invalid mode: %s. Valid options: collect, process, dispatch, all\n", *mode)
		os.Exit(1)
	}

	log.Println("shut down complete")
}

func setupWhatsmeowClient(ctx context.Context, config *utils.Config) (*whatsmeow.Client, chan interface{},error) {
	//Setup a event channel to unbloack whatsapp events
	eventChan := make(chan interface{}, 500)

	//Initilazing the connection to whatsapp
	client, err := wa.Connect(ctx, func(evt interface{}) {
		select {
		case eventChan <- evt:
		default:
			handlers.HandleEvent(evt, config)
			log.Println("event channel is full, dropping event to avoid blocking")
		} // A closure that handles whatsapp events as defined by config
	})

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer client.Disconnect()

	return client, eventChan, err
}

func runCollectMode(
	client *whatsmeow.Client,
	eventChan chan interface{}, 
	ctx context.Context, 
	config *utils.Config,
) {

	log.Println("Starting Collect Mode")

	workerN := runtime.NumCPU()
	wg := startWorkers(ctx, workerN, eventChan, config)

	log.Println("connected to WhatsApp")

	<-ctx.Done()
	log.Println("Collect mode shutting down!")

	close(eventChan)
	wg.Wait()
	log.Println("Collect mode shutdown complete")
}

func runProcessMode(ctx context.Context, config *utils.Config) {
	log.Println("Starting Process Mode")
	//Process message by AI in a non blocking background goroutine
	agents.ProcessMessageByAIPoller(ctx, config)
	log.Println("Process Mode Shutdown")
}

func runDispatchMode(client *whatsmeow.Client, ctx context.Context, Config *utils.Config) {
	log.Println("Running Dispatch Mode")

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

