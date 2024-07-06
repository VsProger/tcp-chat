# TCP-chat (Final Project)

## Description

This is a simple chat application that uses TCP sockets to communicate between a server and multiple clients. The server is able to handle multiple clients at once and can broadcast messages to all connected clients. The server is also able to send private messages to specific clients. The clients are able to send messages to the server and receive messages from the server. 

## Features

- Notifications when a new client connects or disconnects
- Commands such as `/help`, `/ban`, `/kick`, `/users`
- Private messages
- Manage users

## How to run

1. Clone the repository
2. Run the server
```go run .\cmd\server\main.go```
3. Run the client
```ncat localhost 8080``` or you can use any other tool to connect to the server
4. Follow the instructions on the client

## Developers

- Azat Abdirashituly
- Baurzhan Saliyev
- Dias Imakanov
