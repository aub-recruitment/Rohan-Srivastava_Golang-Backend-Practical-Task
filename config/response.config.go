package config

type Response struct {
	Success           bool   `json:"success"`
	Status            string `json:"status"`
	Data              any    `json:"data,omitempty"`
	Error             error  `json:"error,omitempty"`
	CustomErrorNumber int    `json:"int,omitempty"`
}
