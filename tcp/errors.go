package tcp

const (
	TARGET_SERVER TCPErrorType = "TARGET_SERVER"
	INTERNAL                   = "INTERNAL"
)

type TCPErrorType string

type TCPError struct {
	msg     string
	errType TCPErrorType
}

func NewTCPError(msg string) TCPError {
	return TCPError{msg: msg}
}

func (e *TCPError) Error() string {
	return e.msg
}
