package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
)

func EncodeFrame(frame Frame) ([]byte, error) {
	if frame.Version != Version {
		return nil, fmt.Errorf("version mismatch: Current frame's version %d doesn't match expected version: %d", frame.Version, Version)
	}

	if frame.Type >= 6 || frame.Type <= 0 {
		return nil, errors.New("message Type mismatch: The specified message type doesn't exist in our configuration")
	}

	if len(frame.Payload) > int(MaxPayloadSize) {
		return nil, errors.New("message too long: The message length is too long")
	}

	bodyLength := int(FixedBodySize) + len(frame.Payload)
	frameSize := bodyLength + int(LengthPrefixSize)

	encodedFrame := make([]byte, frameSize)
	offset := 0
	binary.BigEndian.PutUint32(
		encodedFrame[offset:offset+int(LengthPrefixSize)],
		uint32(bodyLength),
	)
	offset += int(LengthPrefixSize)
	binary.BigEndian.PutUint16(
		encodedFrame[offset:offset+int(MagicSize)],
		Magic,
	)
	offset += int(MagicSize)

	encodedFrame[offset] = frame.Version
	offset += int(VersionSize)

	encodedFrame[offset] = uint8(frame.Type)
	offset += int(MessageTypeSize)

	binary.BigEndian.PutUint64(
		encodedFrame[offset:offset+int(RequestIDSize)],
		frame.RequestID,
	)
	offset += int(RequestIDSize)

	offset += copy(encodedFrame[offset:], frame.Payload)

	if offset != len(encodedFrame) {
		return nil, fmt.Errorf(
			"protocol: encoded frame size mismatch: wrote %d bytes, allocated %d",
			offset,
			len(encodedFrame),
		)
	}

	return encodedFrame, nil
}
