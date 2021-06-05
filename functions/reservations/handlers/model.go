package handlers

import "github.com/coke-day/functions/reservations/model"

// Reservation model
type Reservation struct {
	RoomName string `json:"room_name" validate:"required,min=2,max=3,roomCustom"`
	Time     string `json:"time" validate:"required,numeric,min=0,max=23"`
}

func buildModel(r Reservation, userEmail string) model.Reservation {
	return model.Reservation{
		RoomName:  r.RoomName,
		Time:      r.Time,
		UserEmail: userEmail,
	}
}

func mapFromModel(u model.Reservation) Reservation {
	return Reservation{
		RoomName: u.RoomName,
		Time:     u.Time,
	}
}
