package repository

import (
	"github.com/Franlky01/bookingwebApp/internal/Models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res Models.Reservation) (int, error)

	InsertRoomRestriction(r Models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]Models.Room, error)
	GetRoomByID(id int) (Models.Room, error)
}
