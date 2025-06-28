package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"os/exec"
	"syscall"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	_ "github.com/mattn/go-sqlite3"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func eventHandler(evt interface{}) {
	switch v:= evt.(type) {
	case *events.Message:
		fmt.Println("Received message:", v.Message.GetConversation())
	}
}

func main() {
	dbLog := waLog.Stdout("DB","DEBUG",true)
	ctx := context.Background()

	container, err := sqlstore.New(ctx, "sqlite3", "file:whatsapp.db?_foreign_keys=on", dbLog)

	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)

	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		fmt.Println("No device ID found, please scan the QR code.")
		
		qrChan, _ := client.GetQRChannel(context.Background())
		err := client.Connect()

		fmt.Println("QR Channel created:", qrChan)

		if err != nil {
			panic(err)
		}

		for evt := range qrChan {

			if evt.Event == "code" {
				fmt.Println("QR Code received, please scan it:", evt.Code) //manually echo the code into the terminal and login to WhatsApp Web to let the bot control your WhatsApp account
				
				cmd := exec.Command("qrencode", "-o", "-", "-t", "UTF8", evt.Code)

				output, err := cmd.Output()

				if err != nil {
					fmt.Println("Error generating QR code:", err)
				} else {
					fmt.Println("QR Code:")
					fmt.Println(string(output))
				}

			} else {
				fmt.Println("Event received:", evt.Event)
			}

		}

	} else {
		err := client.Connect()
		if err != nil {
			panic(err)
		}

		fmt.Println("Using existing device ID:", client.Store.ID)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
	fmt.Println("Client disconnected gracefully.")
}