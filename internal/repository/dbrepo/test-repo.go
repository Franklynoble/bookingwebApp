package dbrepo

import (
	"errors"
	"github.com/Franlky01/bookingwebApp/internal/models"
	"time"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

//InsertReservation insertes reservation into database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	//if the room id is 2, then fail, other wise , pass
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}
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

// GetRoomByID gets a room by id
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room
	if id > 2 {
		return room, errors.New("Some error")
	}

	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	return u, nil
}
func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}
func (m *testDBRepo) Authenticate(email, testPassowrd string) (int, string, error) {

	return 1, "", nil
}
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}

//AllNewReservations returns a slice of all repository
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}

//GetReservationByID resturns one reservation by ID
func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {

	//varaibale to be scanned to
	var res models.Reservation

	return res, nil
}

func (m *testDBRepo) UpdatReservation(u models.Reservation) error {

	return nil
}

// UpdateProcessedForReservation updates process for a  reservation by id
func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {

	return nil

}

//DeleteReservation deletes reservation from the reservation Table
func (m *testDBRepo) DeleteReservation(id int) error {

	return nil

}
func (m *testDBRepo) AllRooms() ([]models.Room, error) {
	var rm []models.Room

	return rm, nil
}
func (m testDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	var restrictions []models.RoomRestriction
	return restrictions, nil
}

//  InsertBlockForRoom inserts a room restriction
func (m *testDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {

	return nil

}

//  DeleteBlockByID deletes a room restriction
func (m *testDBRepo) DeleteBlockByID(id int) error {

	return nil

}
