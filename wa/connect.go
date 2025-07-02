package wa

import (
	"context"
	"fmt"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received message:", v.Message.GetConversation())
	}
}

func Connect(ctx context.Context, handler whatsmeow.EventHandler) (*whatsmeow.Client, error) {

	dbLog := waLog.Stdout("DB", "INFO", true)

	container, err := sqlstore.New(ctx, "sqlite3", "file:whatsapp.db?_foreign_keys=on", dbLog)
	if err != nil {
		return nil, fmt.Errorf("failed to connect client: %w", err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect client: %w", err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(handler)

	if client.Store.ID == nil {
		fmt.Println("No device ID found, please scan the QR code.")
		qrChan, _ := client.GetQRChannel(context.Background())

		if err = client.Connect(); err != nil {
			return nil, fmt.Errorf("failed to connect client: %w", err)
		}

		fmt.Println("Waiting for QR code...")
		for evt := range qrChan {
			fmt.Println("Event received:", evt.Event)
			switch evt.Event {
			case "code":
				fmt.Println("QR Code received, please scan it:", evt.Code)
				output, err := exec.Command("qrencode", "-o", "-", "-t", "UTF8", evt.Code).Output()

				if err != nil {
					fmt.Println("Error generating QR code:", err)
				} else {
					fmt.Println("QR Code:")
					fmt.Println(string(output))
				}
			default:
				fmt.Println("Event received:", evt.Event)
			}
		}
	}

	if err = client.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect client: %w", err)
	}

	return client, nil
}
