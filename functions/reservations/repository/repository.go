package repository

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/coke-day/functions/reservations/model"
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
func (r *ItemRepository) Delete(item model.Reservation) error {
	reservationDB := buildInstanceItem(item)
	if err := r.datastore.Delete(reservationDB.PK, reservationDB.SK); err != nil {
		return err
	}
	return nil
}

// Search items based on filters
func (r *ItemRepository) Search(room, time, userEmail string) ([]model.Reservation, error) {
	var filt = expression.ConditionBuilder{}
	if room != "" {
		filt = expression.And(
			expression.Name("pk").Equal(expression.Value(getPKStartingPart()+room)),
			expression.Name("sk").BeginsWith(getSKStartingPart()+time),
		)
	} else {
		domain, err := utils.GetDomain(userEmail)
		if err != nil {
			return []model.Reservation{}, fmt.Errorf("fail to get domain: %v", err)
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

	var items []instanceItem
	if err := r.datastore.Scan(filt, &items); err != nil {
		return nil, err
	}

	var itemsParsed []model.Reservation
	for _, e := range items {
		itemsParsed = append(itemsParsed, *mapInstanceItemToReservation(e))
	}

	return itemsParsed, nil
}

// Store a new reservation
func (r *ItemRepository) Store(item model.Reservation) error {
	itemDB := buildInstanceItem(item)
	return r.datastore.Store(itemDB)
}

// Get a reservation
func (r *ItemRepository) Get(item model.Reservation) (*model.Reservation, error) {
	storedReservation := &instanceItem{}
	searchItem := buildInstanceItem(item)
	if err := r.datastore.Get(searchItem.PK, searchItem.SK, &storedReservation); err != nil {
		return nil, err
	}

	return mapInstanceItemToReservation(*storedReservation), nil
}
