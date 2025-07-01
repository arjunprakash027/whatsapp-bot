package db

import (
    "database/sql"
    "context"
)

// Struct matching the processed_messages table
type ProcessedMessage struct {
    AiMessageID         int64  // Auto-incremented primary key
    ID                  string 
    ChatJID             string
    SenderJID           string
    Text                string
    AIAddress           string
    AIPrimaryContact    string
    AISecondaryContact  string
    AIMessage           string
    MessageSent         int64 // 0 or 1
}

// Prepared statement for inserting processed messages
var ProcessedInsertStmt *sql.Stmt

func PrepareProcessedInsertStatement(ctx context.Context) error {
    var err error
    ProcessedInsertStmt, err = Conn.PrepareContext(ctx, `
        INSERT INTO processed_messages 
            (id, chat_jid, sender_jid, text, ai_address, ai_primary_contact, ai_secondary_contact, ai_message, message_sent)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
    `)
    return err
}

// Save a processed message
func SaveProcessedMessage(
    ctx context.Context,
    id, chatJID, senderJID, text, aiAddress, aiPrimaryContact, aiSecondaryContact, aiMessage string,
    messageSent int64,
) error {
    _, err := ProcessedInsertStmt.ExecContext(ctx,
        id, chatJID, senderJID, text, aiAddress, aiPrimaryContact, aiSecondaryContact, aiMessage, messageSent,
    )
    return err
}

// Get all processed messages where message_sent = 0 (AI has not sent a reply)
func GetUnsentProcessedMessages(ctx context.Context) ([]ProcessedMessage, error) {
    const q = `
        SELECT ai_message_id, id, chat_jid, sender_jid, text, ai_address, ai_primary_contact, ai_secondary_contact, ai_message, message_sent
        FROM processed_messages
        WHERE message_sent = 0;
    `
    rows, err := Conn.QueryContext(ctx, q)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var msgs []ProcessedMessage
    for rows.Next() {
        var msg ProcessedMessage
        err := rows.Scan(
            &msg.AiMessageID, &msg.ID, &msg.ChatJID, &msg.SenderJID, &msg.Text,
            &msg.AIAddress, &msg.AIPrimaryContact, &msg.AISecondaryContact, &msg.AIMessage, &msg.MessageSent,
        )
        if err != nil {
            return nil, err
        }
        msgs = append(msgs, msg)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return msgs, nil
}

// Mark a processed message as sent by ai_message_id
func MarkProcessedMessageSent(ctx context.Context, aiMessageID int64) error {
    const q = `
        UPDATE processed_messages
        SET message_sent = 1
        WHERE ai_message_id = ?;
    `
    _, err := Conn.ExecContext(ctx, q, aiMessageID)
    return err
}