package model

import "strings"

// Reservation model
type Reservation struct {
	RoomName  string `json:"room_name" validate:"required,min=2,max=3,roomCustom"`
	Time      string `json:"time" validate:"required,numeric,min=0,max=23"`
	UserEmail string `json:"-"`
}

// ReservationPersistence describes dynamodb representation
type ReservationPersistence struct {
	PK        string `dynamodbav:"pk"`
	SK        string `dynamodbav:"sk"`
	UserEmail string `dynamodbav:"user_email"`
}

func (r Reservation) toRoomPersistence() ReservationPersistence {
	return ReservationPersistence{
		PK:        getPKStartingPart() + r.RoomName,
		SK:        getSKStartingPart() + r.Time,
		UserEmail: r.UserEmail,
	}
}

func (r ReservationPersistence) toRoom() Reservation {
	return Reservation{
		RoomName:  strings.TrimLeft(r.PK, getPKStartingPart()),
		Time:      strings.TrimLeft(r.SK, getSKStartingPart()),
		UserEmail: r.UserEmail,
	}
}

func getPKStartingPart() string {
	return "room#"
}

func getSKStartingPart() string {
	return "time#"
}
