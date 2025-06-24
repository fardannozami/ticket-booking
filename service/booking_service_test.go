package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/fardannozami/ticket-booking/helper"
	"github.com/fardannozami/ticket-booking/repository"
	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/ticket-booking")
	helper.PanicIfError(err)

	insertEvent(1, 100, "Event 1")
	insertSeat(2, 1, "A1", "AVAILABLE")
	inserUser(3, "AJitama")

	db.Exec("DELETE FROM bookings")

	exitCode := m.Run()
	db.Close()
	os.Exit(exitCode)
}

func insertEvent(id, quota int, name string) {
	SQL := "INSERT INTO events (id, quota, name) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE quota = VALUES(quota)"
	_, err := db.Exec(SQL, id, quota, name)
	helper.PanicIfError(err)
}

func insertSeat(id, event_id int, seat_number, status string) {
	SQL := "INSERT INTO seats (id, event_id, seat_number, status) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE status = VALUES(status)"
	_, err := db.Exec(SQL, id, event_id, seat_number, status)
	helper.PanicIfError(err)
}

func inserUser(id int, name string) {
	SQL := "INSERT INTO users (id, name) VALUES (?, ?) ON DUPLICATE KEY UPDATE name = VALUES(name)"
	_, err := db.Exec(SQL, id, name)
	helper.PanicIfError(err)
}

func TestBookSeat(t *testing.T) {
	repo := repository.NewBookingRepository()
	service := NewBookingService(db, repo)
	ctx := context.Background()

	id, err := service.BookSeat(ctx, 1, 2, 3)
	assert.NoError(t, err)

	var eventId, seatId, userId int
	SQL := "SELECT event_id, seat_id, user_id FROM bookings WHERE id = ?"
	err = db.QueryRowContext(ctx, SQL, id).Scan(&eventId, &seatId, &userId)
	assert.NoError(t, err)
	assert.Equal(t, 1, eventId)
	assert.Equal(t, 2, seatId)
	assert.Equal(t, 3, userId)
}

func TestBookSeatRaceCondition(t *testing.T) {
	repo := repository.NewBookingRepository()
	service := NewBookingService(db, repo)
	ctx := context.Background()

	totalUsers := 50
	for i := 1; i <= totalUsers; i++ {
		inserUser(i, fmt.Sprintf("user %d", i))
	}

	var wg sync.WaitGroup
	type bookingResult struct {
		userId int
		error  error
	}

	result := make(chan bookingResult, totalUsers)
	wg.Add(totalUsers)

	for i := 1; i <= totalUsers; i++ {
		go func(userId int) {
			defer wg.Done()

			_, err := service.BookSeat(ctx, 1, 2, userId)
			result <- bookingResult{userId: userId, error: err}
		}(i)
	}

	wg.Wait()
	close(result)

	var succesUsers []int

	for res := range result {
		if res.error == nil {
			succesUsers = append(succesUsers, res.userId)
			fmt.Printf("[SUCCESS] user %d berhasil booking\n", res.userId)
		} else {
			fmt.Printf("[FAILED] user %d gagal booking dengan error %v\n", res.userId, res.error)
		}
	}

	assert.Equal(t, 1, len(succesUsers))
	fmt.Printf("jumlah user yang berhasil booking: %d", len(succesUsers))
}
