package entity

type NewsResponse struct {
	Status       string `json:"status"`
	TotalResults uint   `json:"totalResults"`
}
