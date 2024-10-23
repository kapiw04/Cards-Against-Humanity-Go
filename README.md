# Cards Against Humanity Online Game (Go Implementation)

This repository contains the source code for a Cards Against Humanity-like game server, implemented in Go. It features WebSocket-based multiplayer support and RESTful API endpoints to handle game logic and player interactions.

## Features

- WebSocket support for real-time communication.
- REST API for game management.
- Dynamic player management.
- Full game flow: starting a game, playing cards, and evaluating rounds.
- Hand management for each player.
- SQLite database for card storage.
- Simple CORS handling for development.
- Swagger documentation for the API.

## Installation

### Prerequisites

- [Go 1.23.2](https://golang.org/doc/install)
- SQLite database with Cards Against Humanity cards.
- Docker (if using the containerised version).

### Local Development Setup

1. Clone the repository:

    ```bash
    git clone https://github.com/kapiw04/go-app.git
    cd go-app
    ```

2. Install dependencies:

    ```bash
    go mod download
    ```

3. Make sure you have an SQLite database with the necessary card data. This project expects a `cah/cah_cards.db` file in the root directory with the `cards` table populated.

4. Run the server:

    ```bash
    go run .
    ```

5. Swagger documentation will be available at [http://localhost:8080/swagger/](http://localhost:8080/swagger/).

6. WebSocket connection can be established at `ws://localhost:8080/ws`, and the API is available at the listed routes.

### Docker Setup

You can build and run the project using Docker.

    ```bash
    docker compose up
    ```

The game server will be available at `http://localhost:8080`.

## API Endpoints

- **POST** `/start`: Start the game.
- **GET** `/hand`: Retrieve a player's hand of cards.
- **GET** `/black-card`: Get the current black card.
- **POST** `/play-card`: Play a card by sending its index.
- **GET** `/played-cards`: Retrieve all played cards.

For detailed API documentation, please visit the [Swagger UI](http://localhost:8080/swagger/).

## WebSocket

Players connect via WebSocket to the `ws://localhost/ws` endpoint. 

## Joining the Game on the Same Network (Frontend Pending)

1. **Host the Server**: Start the server and find the host machine's local IP (e.g., `192.168.x.x`).

2. **API and WebSocket Access**: Players on the same network can:
   - Make API requests to `http://<host-ip>:8080` for game interactions.
   - Connect to WebSocket at `ws://<host-ip>:8080/ws` for real-time gameplay.

A frontend will be developed later, so for now, players must interact directly via API calls and WebSocket clients.

### CORS Handling

For development, CORS is temporarily configured to allow all origins. Modify the `corsMiddleware` function in `main.go` to restrict origins in production.

### Error Handling

Some basic error handling is included for incorrect card indices and invalid requests. Improve as needed for production.

## Contributing

Feel free to fork this repository, submit issues, or make pull requests. Contributions are always welcome.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.