package labubu

// Labubu represents the domain entity
type Labubu struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

// CreateLabubuRequest represents the request to create a labubu
type CreateLabubuRequest struct {
	Text string `json:"text" validate:"required"`
}