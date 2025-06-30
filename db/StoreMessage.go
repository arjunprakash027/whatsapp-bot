package db

import (

	"database/sql"
	"context"
	"time"
)

var ConvoInsertStmt *sql.Stmt

func PrepareConvoInsertStatement(ctx context.Context) error {
	
	var err error
	
	ConvoInsertStmt, err = Conn.PrepareContext (ctx, `
		INSERT OR IGNORE INTO messages 
			(id, chat_jid, sender_jid, timestamp, text, channel)
		VALUES (?, ?, ?, ?, ?, ?);
	`)

	return err
}

func SaveConvoMessage(
	ctx context.Context,
	id, chat, sender, text, channel string,
	ts time.Time,
) error {

	_, err := ConvoInsertStmt.ExecContext(ctx,
		id, chat, sender, ts.Unix(),
		text, channel,
	)

	return err
}