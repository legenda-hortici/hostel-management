package models

type Service struct {
	Point       int
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
	Is_date     bool   `json:"is_date"`
	Is_hostel   bool   `json:"is_hostel"`
	Is_phone    bool   `json:"is_phone"`
}
