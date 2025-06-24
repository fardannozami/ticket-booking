package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/fardannozami/ticket-booking/helper"
	"github.com/fardannozami/ticket-booking/request"
	"github.com/fardannozami/ticket-booking/response"
	"github.com/fardannozami/ticket-booking/service"
	"github.com/julienschmidt/httprouter"
)

type BookingController interface {
	BookSeat(writer http.ResponseWriter, req *http.Request, params httprouter.Params)
}

type bookingController struct {
	bookingService service.BookingService
}

func NewBookingController(service service.BookingService) BookingController {
	return &bookingController{
		bookingService: service,
	}
}

func (c *bookingController) BookSeat(writer http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var bookingRequest request.BookingCreateRequest
	err := json.NewDecoder(req.Body).Decode(&bookingRequest)
	helper.PanicIfError(err)

	_, err = c.bookingService.BookSeat(context.Background(), bookingRequest.EventId, bookingRequest.SeatId, bookingRequest.UserId)
	helper.PanicIfError(err)

	apiResponse := response.ApiResponse{
		Code:    http.StatusCreated,
		Message: "created",
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(writer).Encode(&apiResponse)
	helper.PanicIfError(err)
}
