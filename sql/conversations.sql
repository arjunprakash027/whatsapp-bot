-- conversations.sql
CREATE TABLE IF NOT EXISTS messages (
    id          TEXT    PRIMARY KEY,          
    chat_jid    TEXT    NOT NULL,
    sender_jid  TEXT    NOT NULL,
    timestamp   INTEGER NOT NULL,             -- Unix seconds
    text        TEXT    NOT NULL,
    channel     TEXT    NOT NULL,
    read_by_ai  INTEGER NOT NULL CHECK (read_by_ai IN (0, 1))
);

CREATE TABLE IF NOT EXISTS processed_messages (
    ai_message_id INTEGER PRIMARY KEY AUTOINCREMENT,  -- Auto-incremented ID for processed message
    id           TEXT    NOT NULL REFERENCES messages(id),
    chat_jid     TEXT    NOT NULL,
    sender_jid   TEXT    NOT NULL,
    text         TEXT    NOT NULL,
    -- AI processed rows

    ai_address TEXT NOT NULL,  -- Address of the AI that processed the message
    ai_primary_contact TEXT NOT NULL,  -- Primary contact of the AI that processed the message
    ai_secondary_contact TEXT NOT NULL,  -- Secondary contact of the AI that processed the message
    ai_message TEXT NOT NULL,  -- The message that was processed by the AI

    message_sent INTEGER NOT NULL CHECK (message_sent IN (0,1)) --check if AI has sent a reply or not

);


CREATE INDEX IF NOT EXISTS idx_chat_time ON messages(chat_jid, timestamp);
