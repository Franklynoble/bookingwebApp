package dbrepo

import (
	"github.com/Franlky01/bookingwebApp/internal/models"
	"time"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

//InsertReservation insertes reservation into database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	return 1, nil
}

//InsertRoomRestriction inserts a room restriction into the Database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	return nil
}

//SearchAvailabilityByDates returns true of Availability exist for roomID, false if no  Availability does not exist
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	return false, nil

}

//SearchAvailabilityForAllRooms returns  a slice of available rooms, if any, for given Date Range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room
	return rooms, nil
}

func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	return room, nil
}
