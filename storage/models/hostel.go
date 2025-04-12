package models

type Hostel struct {
	HostelID       int
	HostelNumber   int
	ResidentsCount int
	RoomCount      int
	OccupiedRooms  int
	AvailableRooms int
	HostelLocation string
	HostelContacts string
	HeadmanName    *string
	HeadmanSurname *string
	HeadmanEmail   *string

	OccupiedPercent  int // <= добавим эти два поля
	AvailablePercent int
}
