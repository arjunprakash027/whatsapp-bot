# WhatsApp Bot

This is a WhatsApp bot built with Go that can connect to WhatsApp, receive messages, and process them with an AI.

## Features

- Connects to WhatsApp using a QR code.
- Receives messages and saves them to a SQLite database.
- Processes messages with an AI (placeholder function).
- Whitelists chats to only process messages from specific groups or users.
- Uses a worker pool to handle events concurrently.
- Dockerized for easy deployment.

## Getting Started

### Prerequisites

- Go 1.23 or higher
- Docker (optional)
- `qrencode` library (`sudo apt-get install qrencode` on Debian/Ubuntu)

### Installation

1.  Clone the repository:
    ```bash
    git clone https://github.com/arjunprakash027/whatsapp-bot
    ```
2.  Install dependencies:
    ```bash
    go mod tidy
    ```
3.  Configure the bot by editing `config.yaml`:
    ```yaml
    Whatsapp:
      WhiteListedChats:
        - "120363027892302133@g.us" # Add your chat JIDs here

    AI:
      Controls:
        PollingInterval: 300 # in seconds
        WorkerCount: 3
    ```
4.  Run the bot:
    ```bash
    go run cmd/bot/main.go
    ```
    You will be prompted to scan a QR code with your phone to connect to WhatsApp.

### Docker

You can also run the bot using Docker:

1.  Build the image:
    ```bash
    docker-compose build
    ```
2.  Run the container:
    ```bash
    docker-compose up
    ```

## How it works

The bot uses the `whatsmeow` library to connect to WhatsApp. When a message is received, it is saved to a SQLite database. A poller then fetches unprocessed messages and sends them to an AI for processing. The AI's response is then saved to the database.

The `main` function in `cmd/bot/main.go` initializes the database, reads the configuration, and connects to WhatsApp. It also starts a pool of workers to handle incoming events from WhatsApp.

The `agents` package contains the AI processing logic. The `db` package handles all database operations. The `wa` package contains the WhatsApp connection and event handling logic.
