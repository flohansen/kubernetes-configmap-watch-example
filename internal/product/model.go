package product

import "time"

type Model struct {
	ID           string         `json:"id"`
	Data         map[string]any `json:"data"`
	LastModified time.Time      `json:"lastModified"`
}
