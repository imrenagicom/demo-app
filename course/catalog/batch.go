package catalog

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	v1 "github.com/imrenagicom/demo-app/pkg/apiclient/course/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BatchStatus int

const (
	BatchStatusDraft BatchStatus = iota
	BatchStatusPublished
	BatchStatusArchived
)

type Batch struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      sql.NullTime
	Name           string
	MaxSeats       int32
	AvailableSeats int32
	Price          float64
	Currency       string
	Status         BatchStatus
	StartDate      sql.NullTime
	EndDate        sql.NullTime
	Version        int64
}

func (b Batch) ApiV1() *v1.Batch {
	var startDate, endDate *timestamppb.Timestamp
	if b.StartDate.Valid {
		startDate = timestamppb.New(b.StartDate.Time)
	}
	if b.EndDate.Valid {
		endDate = timestamppb.New(b.EndDate.Time)
	}

	return &v1.Batch{
		DisplayName: b.Name,
		Name:        string(b.ID.String()), // TODO change with slug
		BatchId:     b.ID.String(),
		Price: &v1.Price{
			Value:    b.Price,
			Currency: b.Currency,
		},
		MaxSeats:       b.MaxSeats,
		AvailableSeats: b.AvailableSeats,
		StartDate:      startDate,
		EndDate:        endDate,
	}
}

var (
	ErrNotEnoughSeats           = errors.New("no seat available")
	ErrClassSoldOut             = errors.New("class is sold out")
	ErrClassNotAvailableForSale = errors.New("class is not available for sale")
)

func (b *Batch) Reserve(ctx context.Context) error {
	if err := b.Available(ctx); err != nil {
		return ErrClassNotAvailableForSale
	}
	if b.AvailableSeats < 1 {
		return ErrNotEnoughSeats
	}
	if b.MaxSeats > 0 {
		b.AvailableSeats -= 1
	}
	return nil
}

func (b *Batch) Available(ctx context.Context) error {
	if b.MaxSeats <= 0 {
		return nil
	}
	if b.AvailableSeats == 0 {
		return ErrClassSoldOut
	}
	if b.EndDate.Valid && time.Now().After(b.EndDate.Time) {
		return ErrClassNotAvailableForSale
	}
	return nil
}

// Allocate increases number of available seats. Only applicable for batch with limited seats.
func (b *Batch) Allocate(ctx context.Context, numSeat int) error {
	if b.MaxSeats > 0 {
		b.AvailableSeats += int32(numSeat)
	}
	return nil
}
