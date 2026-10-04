package payload

type RateLimitRequest struct {
	ClientID string
	Tier     string
	Resource string
	Subject  string
}

type RateLimitResponse struct {
	Allowed      bool
	Limit        uint32
	Remaining    uint32
	RetryAfterMS uint64
}

const (
	ErrorCodeInvalidPayload         uint8 = 1
	ErrorCodeUnsupportedMessageType uint8 = 2
)

type ErrorResponse struct {
	Code    uint8
	Message string
}
