package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fardannozami/ticket-booking/helper"
	"github.com/fardannozami/ticket-booking/repository"
)

type BookingService interface {
	BookSeat(ctx context.Context, eventId, seatId, userId int) (int, error)
}

type bookingService struct {
	Db                *sql.DB
	BookingRepository repository.BookingRepository
}

func NewBookingService(db *sql.DB, repo repository.BookingRepository) BookingService {
	return &bookingService{
		Db:                db,
		BookingRepository: repo,
	}
}

func (s *bookingService) BookSeat(ctx context.Context, eventId int, seatId int, userId int) (int, error) {
	tx, err := s.Db.BeginTx(ctx, nil)
	helper.PanicIfError(err)

	defer func() {
		r := recover()
		if r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	status := s.BookingRepository.GetSeatStatus(ctx, tx, seatId)
	if status != "AVAILABLE" {
		tx.Rollback()
		return 0, errors.New("seat already booked")
	}

	s.BookingRepository.MarkSeatAsBooked(ctx, tx, seatId)
	s.BookingRepository.DecrementEventQuota(ctx, tx, eventId)
	id := s.BookingRepository.InsertBooking(ctx, tx, eventId, seatId, userId)

	return id, tx.Commit()
}
