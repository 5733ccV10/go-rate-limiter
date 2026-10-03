package main

import (
	"errors"
	"go-rlm/internal/payload"
	"go-rlm/internal/protocol"
	"io"
	"log"
	"net"
)

func handler(conn net.Conn) {
	defer conn.Close()
	log.Printf("new client connected from %v", conn.RemoteAddr())
	for {
		frame, err := protocol.DecodeFrame(conn)

		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("client disconnected: %v", conn.RemoteAddr())
			} else {
				log.Printf("failed to decode frame from %v: %v", conn.RemoteAddr(), err)
			}
			return
		}

		log.Printf("received frame type %d, request ID %d from %v", frame.Type, frame.RequestID, conn.RemoteAddr())
		responseFrame, err := handleFrame(frame)

		if err != nil {
			log.Printf("failed to handle frame from %v: %v", conn.RemoteAddr(), err)
			return
		}

		encodedResponseFrame, err := protocol.EncodeFrame(responseFrame)

		if err != nil {
			log.Printf("failed to encode response frame for %v: %v", conn.RemoteAddr(), err)
			return
		}

		if _, err := conn.Write(encodedResponseFrame); err != nil {
			log.Printf("failed to write response to %v: %v", conn.RemoteAddr(), err)
			return
		}
	}
}

func handleFrame(frame protocol.Frame) (protocol.Frame, error) {
	switch frame.Type {
	case protocol.MessageTypePing:
		return handleFramePing(frame)
	default:
		return makeErrorFrame(frame.RequestID, payload.ErrorCodeUnsupportedMessageType, "message type is not supported by this server")
	}
}

func handleFramePing(frame protocol.Frame) (protocol.Frame, error) {
	if len(frame.Payload) > 0 {
		return makeErrorFrame(frame.RequestID, payload.ErrorCodeInvalidPayload, "ping payload must be empty")
	}
	pongFrame := protocol.Frame{
		Version:   protocol.Version,
		Type:      protocol.MessageTypePong,
		RequestID: frame.RequestID,
	}
	return pongFrame, nil
}

func makeErrorFrame(requestID uint64, code uint8, message string) (protocol.Frame, error) {
	encodedPayload, err := payload.EncodeErrorResponse(payload.ErrorResponse{
		Code:    code,
		Message: message,
	})
	if err != nil {
		return protocol.Frame{}, err
	}

	return protocol.Frame{
		Version:   protocol.Version,
		Type:      protocol.MessageTypeError,
		RequestID: requestID,
		Payload:   encodedPayload,
	}, nil
}
