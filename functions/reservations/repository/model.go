package repository

import (
	"github.com/coke-day/functions/reservations/model"
	"strings"
)

// instanceItem describes dynamodb representation
type instanceItem struct {
	PK        string `dynamodbav:"pk"`
	SK        string `dynamodbav:"sk"`
	UserEmail string `dynamodbav:"user_email"`
}

func buildInstanceItem(r model.Reservation) instanceItem {
	return instanceItem{
		PK:        getPKStartingPart() + r.RoomName,
		SK:        getSKStartingPart() + r.Time,
		UserEmail: r.UserEmail,
	}
}

func mapInstanceItemToReservation(instanceItem instanceItem) *model.Reservation {
	return &model.Reservation{
		RoomName:  strings.TrimLeft(instanceItem.PK, getPKStartingPart()),
		Time:      strings.TrimLeft(instanceItem.SK, getSKStartingPart()),
		UserEmail: instanceItem.UserEmail,
	}
}

func getPKStartingPart() string {
	return "room#"
}

func getSKStartingPart() string {
	return "time#"
}
