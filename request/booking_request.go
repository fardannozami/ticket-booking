package request

type BookingCreateRequest struct {
	EventId int `json:"event_id"`
	SeatId  int `json:"seat_id"`
	UserId  int `json:"user_id"`
}
