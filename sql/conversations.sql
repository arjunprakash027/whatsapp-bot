-- conversations.sql
CREATE TABLE IF NOT EXISTS messages (
    id          TEXT    PRIMARY KEY,          
    chat_jid    TEXT    NOT NULL,
    sender_jid  TEXT    NOT NULL,
    timestamp   INTEGER NOT NULL,             -- Unix seconds
    text        TEXT    NOT NULL,
    channel TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_chat_time ON messages(chat_jid, timestamp);
