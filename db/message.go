package db

import (
	"context"
	"database/sql"
	"time"
)

type Message struct {
	ID, ChatJID, SenderJID, Text, Channel string
	Timestamp                             int64
	ReadByAI                              int64
}

var ConvoInsertStmt *sql.Stmt

func PrepareConvoInsertStatement(ctx context.Context) error {

	var err error

	ConvoInsertStmt, err = Conn.PrepareContext(ctx, `
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

func GetConvoMessagesUnProcessed(ctx context.Context) ([]Message, error) {

	const q = `
	    SELECT id, chat_jid, sender_jid, timestamp, text, channel, read_by_ai
	    FROM   messages
	    WHERE  read_by_ai = 0;
	`

	rows, err := Conn.QueryContext(ctx, q)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []Message

	for rows.Next() {
		var msg Message
		err := rows.Scan(
			&msg.ID, &msg.ChatJID, &msg.SenderJID, &msg.Timestamp,
			&msg.Text, &msg.Channel, &msg.ReadByAI,
		)

		if err != nil {
			return nil, err
		}

		msgs = append(msgs, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return msgs, err
}

func UpdateConvoMessageAIPRocessedByID(
	ctx context.Context,
	id string,
	readByAI int64, // 0 for false, 1 for true
) error {
	const q = `
        UPDATE messages
        SET read_by_ai = ?
        WHERE id = ?;
    `

	_, err := Conn.ExecContext(ctx, q, readByAI, id)
	return err
}
