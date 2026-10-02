package payload

import (
	"encoding/binary"
	"errors"
)

func DecodeRateLimitRequest(payloadBytes []byte) (RateLimitRequest, error) {
	offset := 0
	var request RateLimitRequest

	if len(payloadBytes)-offset < 2 {
		return RateLimitRequest{}, errors.New("payload too short to read client id length")
	}

	clientIDLength := int(binary.BigEndian.Uint16(payloadBytes[offset : offset+2]))
	offset += 2

	if len(payloadBytes)-offset < clientIDLength {
		return RateLimitRequest{}, errors.New("payload too short to read client id")
	}

	request.ClientID = string(payloadBytes[offset : offset+clientIDLength])
	offset += clientIDLength

	if len(payloadBytes)-offset < 2 {
		return RateLimitRequest{}, errors.New("payload too short to read tier length")
	}

	tierLength := int(binary.BigEndian.Uint16(payloadBytes[offset : offset+2]))
	offset += 2

	if len(payloadBytes)-offset < tierLength {
		return RateLimitRequest{}, errors.New("payload too short to read tier")
	}

	request.Tier = string(payloadBytes[offset : offset+tierLength])
	offset += tierLength

	if len(payloadBytes)-offset < 2 {
		return RateLimitRequest{}, errors.New("payload too short to read resource length")
	}

	resourceLength := int(binary.BigEndian.Uint16(payloadBytes[offset : offset+2]))
	offset += 2

	if len(payloadBytes)-offset < resourceLength {
		return RateLimitRequest{}, errors.New("payload too short to read resource")
	}

	request.Resource = string(payloadBytes[offset : offset+resourceLength])
	offset += resourceLength

	if len(payloadBytes)-offset < 2 {
		return RateLimitRequest{}, errors.New("payload too short to read subject length")
	}

	subjectLength := int(binary.BigEndian.Uint16(payloadBytes[offset : offset+2]))
	offset += 2

	if len(payloadBytes)-offset < subjectLength {
		return RateLimitRequest{}, errors.New("payload too short to read subject")
	}

	request.Subject = string(payloadBytes[offset : offset+subjectLength])
	offset += subjectLength

	if offset != len(payloadBytes) {
		return RateLimitRequest{}, errors.New("payload size didn't match")
	}

	return request, nil
}

func DecodeRateLimitResponse(payloadBytes []byte) (RateLimitResponse, error) {
	const responseSize = 1 + 4 + 4 + 8
	if len(payloadBytes) != responseSize {
		return RateLimitResponse{}, errors.New("invalid rate limit response payload size")
	}

	offset := 0
	var response RateLimitResponse

	allowed := payloadBytes[offset]
	offset++

	switch allowed {
	case 0:
		response.Allowed = false
	case 1:
		response.Allowed = true
	default:
		return RateLimitResponse{}, errors.New("invalid allowed value in rate limit response")
	}

	response.Limit = binary.BigEndian.Uint32(payloadBytes[offset : offset+4])
	offset += 4
	response.Remaining = binary.BigEndian.Uint32(payloadBytes[offset : offset+4])
	offset += 4
	response.RetryAfterMS = binary.BigEndian.Uint64(payloadBytes[offset : offset+8])
	offset += 8

	if offset != len(payloadBytes) {
		return RateLimitResponse{}, errors.New("rate limit response payload length didn't match")
	}

	return response, nil
}

func DecodeErrorResponse(payloadBytes []byte) (ErrorResponse, error) {
	offset := 0
	var response ErrorResponse

	if len(payloadBytes)-offset < 1 {
		return ErrorResponse{}, errors.New("payload too short to read error code")
	}
	response.Code = payloadBytes[offset]
	offset++

	if len(payloadBytes)-offset < 2 {
		return ErrorResponse{}, errors.New("payload too short to read error message length")
	}
	messageLength := int(binary.BigEndian.Uint16(payloadBytes[offset : offset+2]))
	offset += 2

	if len(payloadBytes)-offset < messageLength {
		return ErrorResponse{}, errors.New("payload too short to read error message")
	}
	response.Message = string(payloadBytes[offset : offset+messageLength])
	offset += messageLength

	if offset != len(payloadBytes) {
		return ErrorResponse{}, errors.New("error response payload length did not match")
	}

	return response, nil
}
