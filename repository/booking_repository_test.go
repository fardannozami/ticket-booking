package repository

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/fardannozami/ticket-booking/helper"
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

func TestGetSeatStatus(t *testing.T) {
	repo := NewBookingRepository()
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	assert.NoError(t, err)

	defer func() {
		r := recover()
		if r != nil {
			tx.Rollback()
			t.Errorf("panic: %v", r)
		} else {
			tx.Commit()
		}
	}()

	status := repo.GetSeatStatus(ctx, tx, 2)
	assert.Equal(t, "AVAILABLE", status)
}
