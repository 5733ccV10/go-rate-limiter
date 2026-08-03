package protocol

const (
	Magic          uint16 = 0x524C
	Version        uint8  = 1
	MaxPayloadSize uint32 = (1 << 17)
)

type Size int

const (
	MagicSize        Size = 2
	VersionSize      Size = 1
	MessageTypeSize  Size = 1
	RequestIDSize    Size = 8
	LengthPrefixSize Size = 4
	FixedBodySize    Size = MagicSize + VersionSize + MessageTypeSize + RequestIDSize
)

type MessageType uint8

const (
	MessageTypeRequest  MessageType = 1
	MessageTypeResponse MessageType = 2
	MessageTypeError    MessageType = 3
	MessageTypePing     MessageType = 4
	MessageTypePong     MessageType = 5
)

type Frame struct {
	Version   uint8
	Type      MessageType
	RequestID uint64
	Payload   []byte
}
