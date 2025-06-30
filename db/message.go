package db

import (

	"database/sql"
	"context"
	"time"
)

type Message struct {
	ID, ChatJID, SenderJID, Text, Channel string
	Timestamp int64
	ReadByAI int64
}


var ConvoInsertStmt *sql.Stmt

func PrepareConvoInsertStatement(ctx context.Context) error {
	
	var err error
	
	ConvoInsertStmt, err = Conn.PrepareContext (ctx, `
		INSERT OR IGNORE INTO messages 
			(id, chat_jid, sender_jid, timestamp, text, channel, read_by_ai)
		VALUES (?, ?, ?, ?, ?, ?, ?);
	`)

	return err
}

func SaveConvoMessage(
	ctx context.Context,
	id, chat, sender, text, channel string,
	ts time.Time,
	read_by_ai int64, // 0 for false, 1 for true
) error {

	_, err := ConvoInsertStmt.ExecContext(ctx,
		id, chat, sender, ts.Unix(),
		text, channel, read_by_ai, 
	)

	return err
}

func GetConvoMessageByID (ctx context.Context, id string) (*Message, error) {
	
	const q = `
	    SELECT id, chat_jid, sender_jid, timestamp, text, channel, read_by_ai
	    FROM   messages
	    WHERE  id = ?;
	`

	var msg Message

	err := Conn.QueryRowContext(ctx, q, id).Scan(
		&msg.ID, &msg.ChatJID, &msg.SenderJID, &msg.Timestamp,
		&msg.Text, &msg.Channel, &msg.ReadByAI,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &msg, err
}

