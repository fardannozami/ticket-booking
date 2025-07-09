package request

type BookingCreateRequest struct {
	EventId int `json:"event_id" example:"1"`
	SeatId  int `json:"seat_id" example:"2"`
	UserId  int `json:"user_id" example:"14"`
}
