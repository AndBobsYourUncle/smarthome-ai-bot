package entities

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}
