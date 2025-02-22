# Chat App

This is a simple chat application that allows users to join chat rooms and send messages. It also includes a bot that can process commands such as fetching stock prices.

## Features

- Join chat rooms
- Send messages
- List available rooms
- Support multiple rooms
- Easy to change room
- Loads previous room messages
- Bot commands:
  - `/help`: shows the help menu
  - `/stock=SYMBOL`: fetches the value of a given stock

### Technical features
- depends exclusively of docker and git to run
- depends exclusively of docker and git to run the tests

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Running the Application

1. Clone the repository:

    ```sh
    git clone https://github.com/ozzono/chat-app
    cd chat-app
    ```

2. Build and run the application using Docker Compose:

    ```sh
    make start
    ```

3. Open your browser and navigate to `http://localhost:8080` to access the chat application.

### API Endpoints

- **Join Room**: `ws /api/v1/rooms/{room}/bind`
- **Send Message**: `ws /api/v1/rooms/{room}/{nickname}/send?content={message}`
- **List Rooms**: `GET /api/v1/rooms`
- These can be tested using [open api](http://localhost:8080/swagger/index.html)

### Bot Commands

- **Help**: `/help`
- **Stock**: `/stock=SYMBOL`

### Example Usage

1. Open the chat application in your browser.
2. Enter a room ID and nickname, then click "Log into Room".
3. Type a message and click "Send".
4. Use bot commands like `/help` and `/stock=AAPL.US` in the message input.

### Running Tests

To run the tests, use the following command:

```sh
make test
```

## Notes
- This project approached the log in as simple as possible without session management;
- To use minimal resources, I chose to use in-memory sqLite as database;
- To use minimal resources, I chose to use a runtime queue and worker system;