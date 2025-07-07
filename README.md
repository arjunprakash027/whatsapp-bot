# WhatsApp House Hunter Bot

This project was born out of a personal struggle: the daunting task of finding a house in Dublin. I was a member of numerous WhatsApp groups, all flooded with housing ads. Manually sifting through each message, trying to find the right fit, was a time-consuming and frustrating experience.

I created this bot for me to swift through these groups and send messages on my behalf. It automates the process of monitoring these WhatsApp groups, filtering messages, and (eventually) using AI to identify the most promising housing opportunities. While the AI component is a work in progress, the bot is fully functional as a WhatsApp message processor.

## The Vision

The ultimate goal for this bot is to be a fully autonomous house-hunting assistant. Here's the plan:

1.  **Monitor WhatsApp Groups:** The bot will be a member of various Dublin housing groups.
2.  **Filter and Store Messages:** It will listen for new messages, saving them to a database.
3.  **AI-Powered Analysis:** An AI model will analyze the messages, extracting key details like price, location, and number of rooms. It will then use this information to determine if the property is a good match for my needs.
4.  **Automated Responses:** If the AI finds a good match, it will automatically send a message to the poster, expressing interest.

## How it Works (The Current State)

Currently, the bot is a powerful message-processing engine. Here's how it works:

- **Connection:** The bot uses the `whatsmeow` library to connect to WhatsApp.
- **Message Handling:** Incoming messages are saved to a SQLite database.
- **Filtering:** The bot uses a `BenchmarkMessage` (configured in `config.yaml`) to calculate the Levenshtein distance to incoming messages. This allows it to filter out irrelevant messages and only focus on those that are likely to be housing ads.
- **Dispatching:** The bot can be configured to send messages to a specific chat, which can be used to notify me of potential matches.

## Modes of Operation

The bot can be run in four different modes:

-   **`collect`:** In this mode, the bot connects to WhatsApp and collects messages from the whitelisted chats, saving them to the database.
-   **`process`:** This mode processes the collected messages. It calculates the Levenshtein distance between each message and the `BenchmarkMessage` and saves the results.
-   **`dispatch`:** This mode sends the processed messages to the `MessageReceiver` specified in the `config.yaml` file.
-   **`all`:** This is the default mode, which runs all three modes (`collect`, `process`, and `dispatch`) concurrently.

You can specify the mode using the `--mode` flag when running the bot:

```bash
go run cmd/bot/main.go --mode=collect
```

## What does it not have yet

- **AI:** The bot has all the groundwork to integrate AI, but I have not integrated AI yet and will do it while I find time (I found a house before I could integrate AI..luckly or sadly)
- **Bad way to find similarity:** Edit distance is not the right way to find similairty in my case where the key and query are both very large sentences, a better way would have been embedding them using doc to vectors and calculating their dot product. There are lots of other ways and edit distance is not one of them. I did it just for my learning

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
      MessageReceiver: "120363027892302133@g.us" # The chat to send notifications to

    AI:
      Controls:
        PollingInterval: 300 # in seconds
        WorkerCount: 3
      BenchmarkMessage: |
        # This is the message I use edit distance to compare with other messages to decide wether to use AI to process a message or not
        Immediate permanent /Temporary accommodations available
        LOOKING for **TIDY,  CLEAN , QUITE,Non smoker,Non Drinker*
        *Permanent accommodation for *:
        1) Double occupancy available in a *Triple* sharing bedroom: 700 per person including bills. (Male only)
        Deposit same as rent
        2)Single room 900 including bills ( Male/Female)
        3) Double room (Male/Female)
        700 per person including bills
        Deposit same as rent
        It's 4 bhk house Two toilets and one bathroom
        📍 Prime Location: Drumkeen manor, Dun Laoghaire
        Bus stop 1 min walk 7 7a 111 45a
        Supervalue and shopping center 1 min walk
        Lidl ,aldi ,tesco ,centra, all groceries electronics  10 min walk
    ```
4.  Run the bot:
    ```bash
    go run cmd/bot/main.go --mode=all
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