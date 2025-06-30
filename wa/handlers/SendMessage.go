package handlers

import (
	"context"
	"log"

	"go.mau.fi/whatsmeow"
	"google.golang.org/protobuf/proto"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func SendText(
	ctx context.Context, 
	client *whatsmeow.Client, 
	toJIDStr, body string,
) error {
	toJID, err := types.ParseJID(toJIDStr)
	if err != nil {
		return err
	}

	msg := &waProto.Message{
		Conversation: proto.String(body),
	}

	_, err = client.SendMessage(ctx, toJID, msg)

	log.Printf("Sent message to %s: %s", toJIDStr, body)
	
	return err
}


