package model

type UpdateSource struct {
	Id   int64  `json:"id" validate:"required"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}
