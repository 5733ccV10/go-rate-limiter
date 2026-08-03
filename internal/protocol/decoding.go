package protocol

import (
	"encoding/binary"
	"errors"
	"io"
)

func DecodeFrame(reader io.Reader) (Frame, error) {
	lengthInBytes := make([]byte, LengthPrefixSize)
	if _, err := io.ReadFull(reader, lengthInBytes); err != nil {
		return Frame{}, errors.New("something went wrong while calculating the length of the payload")
	}

	bodyLength := binary.BigEndian.Uint32(lengthInBytes)
	if bodyLength < uint32(FixedBodySize) {
		return Frame{}, errors.New("body length too small to parse")
	}
	if bodyLength > uint32(FixedBodySize)+MaxPayloadSize {
		return Frame{}, errors.New("body frame too large to decode")
	}

	bodyBytes := make([]byte, int(bodyLength))
	if _, err := io.ReadFull(reader, bodyBytes); err != nil {
		return Frame{}, errors.New("something went wrong while extracting the payload bytes")
	}

	offset := 0
	magic := binary.BigEndian.Uint16(bodyBytes[offset : offset+int(MagicSize)])
	offset += int(MagicSize)

	if magic != Magic {
		return Frame{}, errors.New("something went wrong while decoding magic")
	}

	version := bodyBytes[offset]
	offset += int(VersionSize)

	if version != Version {
		return Frame{}, errors.New("version mismatch while decoding")
	}

	messageType := MessageType(bodyBytes[offset])
	offset += int(MessageTypeSize)

	if messageType < 1 || messageType > 5 {
		return Frame{}, errors.New("message type mismatch while decoding")
	}

	requestId := binary.BigEndian.Uint64(bodyBytes[offset : offset+int(RequestIDSize)])
	offset += int(RequestIDSize)

	payload := bodyBytes[offset:]
	offset += len(payload)

	if offset != len(bodyBytes) {
		return Frame{}, errors.New("frame body was not fully decoded")
	}

	return Frame{
		Version:   version,
		Type:      messageType,
		RequestID: requestId,
		Payload:   payload,
	}, nil
}
