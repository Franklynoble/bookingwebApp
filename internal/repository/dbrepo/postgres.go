package dbrepo

import (
	"context"
	"github.com/Franlky01/bookingwebApp/internal/Models"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

//InsertReservation insertes reservation into database
func (m *postgresDBRepo) InsertReservation(res Models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var newID int
	stmt := `insert into reservations(first_name,last_name,email,phone,
		start_date,end_date, room_id, created_at, updated_at)
      	values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`
	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID) // this function  returns a  row using the id

	if err != nil {
		return 0, err
	}
	return newID, nil
}

//InsertRoomRestriction inserts a room restriction into the Database
func (m *postgresDBRepo) InsertRoomRestriction(r Models.RoomRestriction) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `insert into room_restriction(start_date, end_date,room_id,reservation_id,
   			created_at,updated_at,restriction_id)
			values
    			($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)
	if err != nil {
		return err
	}

	return nil
}

//SearchAvailabilityByDates returns true of Availability exist for roomID, false if no  Availability does not exist
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	/*
		---existing reservation from 2021-02-03 to 2021-02-05
			---search date is exactly the same as existing reservation
	*/
	var numRows int
	query := `	
	select
		count(id)
		from
		room_restriction
		where
		      room_id = $1 
		  and $2 < end_date and $3 > start_date`

	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)

	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}
	return false, nil

}

//SearchAvailabilityForAllRooms returns  a slice of available rooms, if any, for given Date Range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]Models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []Models.Room
	query := `select 
r.id, r.room_name
from
 rooms r 
where   r.id not in
( select rr.room_id from room_restriction rr where $1 <rr.end_date and $2 > rr.start_date)`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room Models.Room
		err := rows.Scan(&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		//append to slice rooms
		rooms = append(rooms, room)
	}
	if err != nil {
		return rooms, err

	}
	return rooms, nil
}

//GetsRoomByID gets a roomby id
func (m *postgresDBRepo) GetRoomByID(id int) (Models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var room Models.Room

	query := `select id, room_name, created_at, updated_at from rooms where id=$1`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}
	return room, nil
}
