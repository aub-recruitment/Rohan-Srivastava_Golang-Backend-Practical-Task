package config

type APIError struct {
	Code              int    `json:"code"`
	Message           string `json:"message"`
	Error             error  `json:"error,omitempty"`
	CustomErrorNumber int    `json:"int,omitempty"`
}
