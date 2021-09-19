package repository

import "github.com/Franlky01/bookingwebApp/internal/Models"

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res Models.Reservation) error
}
