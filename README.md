# Chat App

This is a simple chat application that allows users to join chat rooms and send messages. It also includes a bot that can process commands such as fetching stock prices.

## Features

- Join chat rooms
- Send messages
- List available rooms
- Bot commands:
  - `/help`: shows the help menu
  - `/stock=SYMBOL`: fetches the value of a given stock

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
    docker-compose up --build
    ```

3. Open your browser and navigate to `http://localhost:8080` to access the chat application.

### API Endpoints

- **Join Room**: `GET /api/v1/rooms/{room}/bind`
- **Send Message**: `GET /api/v1/rooms/{room}/{nickname}/send?content={message}`
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
go test ./...
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
````

</file>

This README file includes:
1. A brief description of the chat application and its features.
2. Prerequisites for running the application.
3. Instructions on how to build and run the application using Docker Compose.
4. API endpoints and bot commands.
5. Example usage of the application.
6. Instructions for running tests.
7. License information.