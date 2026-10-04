package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"go-rlm/client"
	"go-rlm/internal/payload"
	"go-rlm/internal/protocol"
)

const serverAddress = "127.0.0.1:8000"

func main() {
	conn, err := client.DialServer(serverAddress)
	if err != nil {
		log.Fatalf("failed to connect to %s: %v", serverAddress, err)
	}
	defer conn.Close()

	if err := sendRequest(conn); err != nil {
		log.Fatalf("request to %s failed: %v", serverAddress, err)
	}

	log.Printf("request to %s completed", serverAddress)
}

func sendRequest(conn net.Conn) error {
	requestType := "ping"
	if len(os.Args) > 1 {
		requestType = os.Args[1]
	}

	switch requestType {
	case "ping":
		return sendPing(conn)
	case "rate-limit":
		return errors.New("rate-limit requests are not implemented yet")
	default:
		return fmt.Errorf("unknown request type %q: use ping or rate-limit", requestType)
	}
}

func sendPing(conn net.Conn) error {
	request := protocol.Frame{
		Version:   protocol.Version,
		Type:      protocol.MessageTypePing,
		RequestID: 1,
	}

	encoded, err := protocol.EncodeFrame(request)
	if err != nil {
		return fmt.Errorf("encode frame: %w", err)
	}
	if _, err := conn.Write(encoded); err != nil {
		return fmt.Errorf("write frame: %w", err)
	}

	response, err := protocol.DecodeFrame(conn)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if response.RequestID != request.RequestID {
		return fmt.Errorf("expected request ID %d, got %d", request.RequestID, response.RequestID)
	}
	if response.Type == protocol.MessageTypeError {
		serverError, err := payload.DecodeErrorResponse(response.Payload)
		if err != nil {
			return fmt.Errorf("decode server error: %w", err)
		}
		return fmt.Errorf("server error %d: %s", serverError.Code, serverError.Message)
	}
	if response.Type != protocol.MessageTypePong {
		return fmt.Errorf("expected Pong, got message type %d", response.Type)
	}
	if len(response.Payload) != 0 {
		return fmt.Errorf("expected empty Pong payload, got %d bytes", len(response.Payload))
	}

	return nil
}
