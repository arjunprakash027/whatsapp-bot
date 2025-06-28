package utils

import (
    "context"

    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types"
)

func FetchHistroyRequest(ctx context.Context, c *whatsmeow.Client, jid types.JID, count int) error {
	
    req := c.BuildHistoryRequest(&types.MessageInfo{})
}