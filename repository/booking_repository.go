package repository

import (
	"context"
	"database/sql"

	"github.com/fardannozami/ticket-booking/helper"
)

type BookingRepository interface {
	GetSeatStatus(ctx context.Context, tx *sql.Tx, seatId int) string
	MarkSeatAsBooked(ctx context.Context, tx *sql.Tx, seatId int)
	DecrementEventQuota(ctx context.Context, tx *sql.Tx, eventId int)
	InsertBooking(ctx context.Context, tx *sql.Tx, eventId, seatId, userId int) int
}

type bookingRepository struct{}

func NewBookingRepository() BookingRepository {
	return &bookingRepository{}
}

func (repo *bookingRepository) GetSeatStatus(ctx context.Context, tx *sql.Tx, seatId int) string {
	var status string
	SQL := "SELECT status FROM seats WHERE id = ? FOR UPDATE"
	err := tx.QueryRowContext(ctx, SQL, seatId).Scan(&status)
	helper.PanicIfError(err)

	return status
}

func (repo *bookingRepository) MarkSeatAsBooked(ctx context.Context, tx *sql.Tx, seatId int) {
	SQL := "UPDATE seats SET status = 'BOOKED' WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, seatId)
	helper.PanicIfError(err)
}

func (repo *bookingRepository) DecrementEventQuota(ctx context.Context, tx *sql.Tx, eventId int) {
	SQL := "UPDATE events SET quota = quota - 1 WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, eventId)
	helper.PanicIfError(err)
}

func (repo *bookingRepository) InsertBooking(ctx context.Context, tx *sql.Tx, eventId int, seatId int, userId int) int {
	SQL := "INSERT INTO bookings (event_id, seat_id, user_id) VALUES (?, ?, ?)"
	result, err := tx.ExecContext(ctx, SQL, eventId, seatId, userId)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	return int(id)
}
