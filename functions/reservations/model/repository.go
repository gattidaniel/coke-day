package model

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/coke-day/pkg/datastore"
	"github.com/coke-day/pkg/utils"
	"github.com/coke-day/pkg/validators"
	"strings"
)

// ItemRepository interfaces with items table
type ItemRepository struct {
	datastore datastore.Datastore
}

// NewItemRepository returns a new instance of item repository
func NewItemRepository(ds datastore.Datastore) *ItemRepository {
	return &ItemRepository{datastore: ds}
}

// Delete a single item
func (r *ItemRepository) Delete(item *Reservation) error {
	reservationDB := item.toRoomPersistence()
	if err := r.datastore.Delete(reservationDB.PK, reservationDB.SK); err != nil {
		return err
	}
	return nil
}

// Search items based on filters
func (r *ItemRepository) Search(room, time, userEmail string) ([]Reservation, error) {
	var filt = expression.ConditionBuilder{}
	if room != "" {
		filt = expression.And(
			expression.Name("pk").Equal(expression.Value(getPKStartingPart()+room)),
			expression.Name("sk").BeginsWith(getSKStartingPart()+time),
		)
	} else {
		domain, err := utils.GetDomain(userEmail)
		if err != nil {
			return []Reservation{}, fmt.Errorf("fail to get domain: %v", err)
		}
		if strings.Contains(domain, validators.CokeDomain) {
			filt = expression.And(
				expression.Name("pk").BeginsWith(getPKStartingPart()+"C"),
				expression.Name("sk").BeginsWith(getSKStartingPart()+time),
			)
		}
		if strings.Contains(domain, validators.PepsiDomain) {
			filt = expression.And(
				expression.Name("pk").BeginsWith(getPKStartingPart()+"P"),
				expression.Name("sk").BeginsWith(getSKStartingPart()+time),
			)
		}
	}

	var items []ReservationPersistence
	if err := r.datastore.Scan(filt, &items); err != nil {
		return nil, err
	}

	var itemsParsed []Reservation
	for _, e := range items {
		itemsParsed = append(itemsParsed, e.toRoom())
	}

	return itemsParsed, nil
}

// Store a new reservation
func (r *ItemRepository) Store(item *Reservation) error {
	itemDB := item.toRoomPersistence()
	return r.datastore.Store(itemDB)
}

// Get a reservation
func (r *ItemRepository) Get(item *Reservation) (*ReservationPersistence, error) {
	storedReservation := &ReservationPersistence{}
	searchItem := item.toRoomPersistence()
	if err := r.datastore.Get(searchItem.PK, searchItem.SK, &storedReservation); err != nil {
		return nil, err
	}

	return storedReservation, nil
}
