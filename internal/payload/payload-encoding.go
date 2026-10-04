package payload

import (
	"encoding/binary"
	"errors"
	"go-rlm/internal/protocol"
)

func EncodeRateLimitRequest(rateLimitRequest RateLimitRequest) ([]byte, error) {
	messageLength := len(rateLimitRequest.Tier) + len(rateLimitRequest.Resource) + len(rateLimitRequest.ClientID) + len(rateLimitRequest.Subject)
	messageLength += 8

	if messageLength > int(protocol.MaxPayloadSize) {
		return nil, errors.New("message too large to send")
	}

	encodedRateLimitRequest := make([]byte, messageLength)
	offset := 0
	if len(rateLimitRequest.ClientID) >= (1 << 16) {
		return nil, errors.New("client id too large")
	}

	binary.BigEndian.PutUint16(encodedRateLimitRequest[offset:offset+2], uint16(len(rateLimitRequest.ClientID)))
	offset += 2
	copy(encodedRateLimitRequest[offset:], rateLimitRequest.ClientID)
	offset += len(rateLimitRequest.ClientID)

	if len(rateLimitRequest.Tier) >= (1 << 16) {
		return nil, errors.New("tier too large")
	}

	binary.BigEndian.PutUint16(encodedRateLimitRequest[offset:offset+2], uint16(len(rateLimitRequest.Tier)))
	offset += 2
	copy(encodedRateLimitRequest[offset:], rateLimitRequest.Tier)
	offset += len(rateLimitRequest.Tier)

	if len(rateLimitRequest.Resource) >= (1 << 16) {
		return nil, errors.New("resource too large")
	}

	binary.BigEndian.PutUint16(encodedRateLimitRequest[offset:offset+2], uint16(len(rateLimitRequest.Resource)))
	offset += 2
	copy(encodedRateLimitRequest[offset:], rateLimitRequest.Resource)
	offset += len(rateLimitRequest.Resource)

	if len(rateLimitRequest.Subject) >= (1 << 16) {
		return nil, errors.New("subject too large")
	}

	binary.BigEndian.PutUint16(encodedRateLimitRequest[offset:offset+2], uint16(len(rateLimitRequest.Subject)))
	offset += 2
	copy(encodedRateLimitRequest[offset:], rateLimitRequest.Subject)
	offset += len(rateLimitRequest.Subject)

	if offset != len(encodedRateLimitRequest) {
		return nil, errors.New("something went wrong while encoding the message")
	}

	return encodedRateLimitRequest, nil
}

func EncodeRateLimitResponse(rateLimitResponse RateLimitResponse) ([]byte, error) {
	const responseSize = 1 + 4 + 4 + 8

	encodedRateLimitResponse := make([]byte, responseSize)
	offset := 0

	if rateLimitResponse.Allowed {
		encodedRateLimitResponse[offset] = 1
	}
	offset++

	binary.BigEndian.PutUint32(
		encodedRateLimitResponse[offset:offset+4],
		rateLimitResponse.Limit,
	)
	offset += 4

	binary.BigEndian.PutUint32(
		encodedRateLimitResponse[offset:offset+4],
		rateLimitResponse.Remaining,
	)
	offset += 4

	binary.BigEndian.PutUint64(
		encodedRateLimitResponse[offset:offset+8],
		rateLimitResponse.RetryAfterMS,
	)
	offset += 8

	if offset != len(encodedRateLimitResponse) {
		return nil, errors.New("rate-limit response was not fully encoded")
	}

	return encodedRateLimitResponse, nil
}

func EncodeErrorResponse(errorResponse ErrorResponse) ([]byte, error) {
	if len(errorResponse.Message) >= 1<<16 {
		return nil, errors.New("error message too large")
	}

	messageLength := 1 + 2 + len(errorResponse.Message)
	if messageLength > int(protocol.MaxPayloadSize) {
		return nil, errors.New("error response too large")
	}

	encodedErrorResponse := make([]byte, messageLength)
	offset := 0

	encodedErrorResponse[offset] = errorResponse.Code
	offset++

	binary.BigEndian.PutUint16(
		encodedErrorResponse[offset:offset+2],
		uint16(len(errorResponse.Message)),
	)
	offset += 2

	copy(encodedErrorResponse[offset:], errorResponse.Message)
	offset += len(errorResponse.Message)

	if offset != len(encodedErrorResponse) {
		return nil, errors.New("error response was not fully encoded")
	}

	return encodedErrorResponse, nil
}
