package dcubeerrs

type Error interface {
	StatusCode() int
	Message() string
}

type DcubeError struct {
	statusCode int
	message    string
}

func New(statusCode int, message string) *DcubeError {
	return &DcubeError{statusCode: statusCode, message: message}
}

func (e *DcubeError) StatusCode() int {
	return e.statusCode
}

func (e *DcubeError) Message() string {
	return e.message
}
