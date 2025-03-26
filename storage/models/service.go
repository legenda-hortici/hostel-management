package models

type Service struct {
	ID          int
	Name        string
	Type        string
	Amount      int
	Description string
	Is_date     bool
	Is_hostel   bool
	Is_phone    bool
}
