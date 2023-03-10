# Chat Server

This is a chat server built in Golang to demonstrate how to design scalable chat services using Redis and the publisher-subscriber model.

## Overview

I designed the chat server to be horizontally scalable by using Redis to store user states, including subscribed channels and login information. Additionally, the server includes authentication and a middleware for the WebSocket endpoint.

## Running the Server

To run the server, you need to run the `docker-compose.yaml` file with the following command:

`docker-compose up`

## Endpoints

The server supports the following endpoints:

### /signup

This endpoint requires a JSON with a username and password. It returns an OTP that must be added as a query parameter when connecting to the WebSocket endpoint.

### /login

This endpoint requires a JSON with a username and password. It returns an OTP that must be added as a query parameter when connecting to the WebSocket endpoint.

### /community

This is a WebSocket endpoint that should be accessed via `ws://localhost:8080?OTP="otp"`. The message to the WebSocket is in JSON format and has the following structure:

{"command":"", "content":"", "channel":"", "error":""}


The `command` field can be one of "SendMessage", "Subscribe" or "Unsubscribe". The `channel` field must be provided. By default, all users are subscribed to the "eje" channel.

## Testing Horizontal Scalability

To test the horizontal scalability, create another session on your websocket client, signup and login using `localhost:8081` and establish a websocket connection on `ws://localhost:8081?OTP="otp"`.
Keeping the two websocket sessions active, try sending messages.


## Conclusion

With its support for authentication and scalable design, this chat server provides a solid foundation for building real-world chat applications. 
