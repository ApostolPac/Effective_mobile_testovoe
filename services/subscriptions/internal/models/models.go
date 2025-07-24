package models


type Subscription struct {
	Id          int       `json:"id,omitempty"`
	ServiceName string    `json:"service_name,omitempty"`
	Price       int       `json:"price,omitempty"`
	UserId      string    `json:"user_id,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string    `json:"end_date,omitempty"`
	TotalSum    int       `json:"total_sum,omitempty"`
}

type ShowSubscSum struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
