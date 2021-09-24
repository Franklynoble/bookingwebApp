package dbrepo

import (
	"github.com/Franlky01/bookingwebApp/internal/Models"
	"time"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

//InsertReservation insertes reservation into database
func (m *testDBRepo) InsertReservation(res Models.Reservation) (int, error) {

	return 1, nil
}

//InsertRoomRestriction inserts a room restriction into the Database
func (m *testDBRepo) InsertRoomRestriction(r Models.RoomRestriction) error {

	return nil
}

//SearchAvailabilityByDates returns true of Availability exist for roomID, false if no  Availability does not exist
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	return false, nil

}

//SearchAvailabilityForAllRooms returns  a slice of available rooms, if any, for given Date Range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]Models.Room, error) {

	var rooms []Models.Room
	return rooms, nil
}

func (m *testDBRepo) GetRoomByID(id int) (Models.Room, error) {
	var room Models.Room
	return room, nil
}
