package model

type VerifyRequest struct {
	Email string `json:"email"`
}

type VerifyResponse struct {
	Email  string `json:"email"`
	Syntax bool   `json:"syntax"`
	Domain bool   `json:"domain"`
	MX     bool   `json:"mx"`
	SMTP   string `json:"smtp"`
	Status string `json:"status"`
}
