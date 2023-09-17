package utils

type APIResponse struct {
	Payload    interface{} `json:"payload"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
}
