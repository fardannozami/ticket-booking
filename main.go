package main

import (
	"net/http"

	"github.com/fardannozami/ticket-booking/app"
	"github.com/fardannozami/ticket-booking/controller"
	"github.com/fardannozami/ticket-booking/exception"
	"github.com/fardannozami/ticket-booking/helper"
	"github.com/fardannozami/ticket-booking/repository"
	"github.com/fardannozami/ticket-booking/service"
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/fardannozami/ticket-booking/docs"
	_ "github.com/go-sql-driver/mysql"
)

func adaptHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h.ServeHTTP(w, r)
	}
}

func main() {
	db := app.NewSqlDb()

	bookingRepo := repository.NewBookingRepository()
	bookingService := service.NewBookingService(db, bookingRepo)
	bookingController := controller.NewBookingController(bookingService)

	router := httprouter.New()

	router.GET("/documentation/*any", adaptHandler(httpSwagger.WrapHandler))
	router.POST("/api/book-seat", bookingController.BookSeat)

	router.PanicHandler = exception.PanicHandler

	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
