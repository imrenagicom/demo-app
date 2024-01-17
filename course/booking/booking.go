package booking

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/imrenagicom/demo-app/course/catalog"
	pu "github.com/imrenagicom/demo-app/internal/proto"
	v1 "github.com/imrenagicom/demo-app/pkg/apiclient/course/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Status int

func (s Status) ApiV1() v1.Status {
	switch s {
	case StatusCreated:
		return v1.Status_CREATED
	case StatusReserved:
		return v1.Status_RESERVED
	case StatusCompleted:
		return v1.Status_COMPLETED
	case StatusFailed:
		return v1.Status_FAILED
	case StatusExpired:
		return v1.Status_EXPIRED
	default:
		return v1.Status_BOOKING_UNSPECIFIED
	}
}

const (
	StatusUnknown Status = iota
	StatusCreated
	StatusReserved
	StatusCompleted
	StatusFailed
	StatusExpired
)

type builder struct {
	b *Booking
}

func (b *builder) WithCustomer(name string, email string, phone string) *builder {
	b.b.Customer = Customer{
		Name:  name,
		Email: email,
		Phone: sql.NullString{Valid: true, String: phone},
	}
	return b
}

func (b *builder) Build() *Booking {
	return b.b
}

func For(c *catalog.Course, b *catalog.Batch) *builder {
	booking := &Booking{
		ID:        uuid.New(),
		Course:    c,
		Batch:     b,
		Price:     b.Price,
		Currency:  b.Currency,
		Status:    StatusCreated,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return &builder{
		b: booking,
	}
}

type Booking struct {
	ID            uuid.UUID
	Course        *catalog.Course
	Batch         *catalog.Batch
	NumTickets    int64
	Price         float64
	Currency      string
	Status        Status
	ReservedAt    sql.NullTime
	ExpiredAt     sql.NullTime
	PaidAt        sql.NullTime
	FailedAt      sql.NullTime
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     sql.NullTime
	PaymentType   sql.NullString
	InvoiceNumber sql.NullString
	Version       int64
	Customer      Customer
}

func (b *Booking) CompletePayment(ctx context.Context, paidAt time.Time) error {
	b.Status = StatusCompleted
	b.PaidAt = sql.NullTime{
		Time:  paidAt,
		Valid: true,
	}
	b.FailedAt = sql.NullTime{}
	return nil
}

func (b *Booking) FailPayment(ctx context.Context, failedAt time.Time) error {
	b.Status = StatusFailed
	b.FailedAt = sql.NullTime{
		Time:  failedAt,
		Valid: true,
	}
	b.PaidAt = sql.NullTime{}
	return nil
}

// UpdatePayment will update the payment type and reset the payment id created before
func (b *Booking) UpdatePayment(ctx context.Context, paymentType string) error {
	b.PaymentType = sql.NullString{Valid: true, String: paymentType}
	b.InvoiceNumber = sql.NullString{}
	return nil
}

const (
	bookingHoldDuration = 10 * time.Minute
)

func (b *Booking) Reserve(ctx context.Context, batch *catalog.Batch) error {
	if err := batch.Available(ctx); err != nil {
		return err
	}
	if err := batch.Reserve(ctx); err != nil {
		return err
	}
	now := time.Now()
	b.Status = StatusReserved
	b.ReservedAt = sql.NullTime{
		Time:  now,
		Valid: true,
	}
	b.ExpiredAt = sql.NullTime{
		Time:  now.Add(bookingHoldDuration),
		Valid: true,
	}
	return nil
}

func (b *Booking) Expire(ctx context.Context) error {
	if b.Status == StatusExpired {
		return ErrBookingAlreadyExpired
	}
	if b.Status == StatusCompleted || b.Status == StatusFailed {
		return ErrBookingAlreadyCompleted
	}
	b.Status = StatusExpired
	b.UpdatedAt = time.Now()
	return nil
}

func (b Booking) ApiV1() *v1.Booking {
	var course *v1.Course
	if b.Course != nil {
		course = b.Course.ApiV1()
	}
	var batch *v1.Batch
	if b.Batch != nil {
		batch = b.Batch.ApiV1()
	}

	return &v1.Booking{
		Number:     b.ID.String(),
		Course:     course.GetCourseId(),
		Batch:      batch.GetBatchId(),
		Price:      b.Price,
		Currency:   b.Currency,
		Status:     b.Status.ApiV1(),
		CreatedAt:  timestamppb.New(b.CreatedAt),
		ReservedAt: pu.FromSQLNullTime(b.ReservedAt),
		PaidAt:     pu.FromSQLNullTime(b.PaidAt),
		ExpiredAt:  pu.FromSQLNullTime(b.ExpiredAt),
		FailedAt:   pu.FromSQLNullTime(b.FailedAt),
		Customer: &v1.Customer{
			Name:        b.Customer.Name,
			Email:       b.Customer.Email,
			PhoneNumber: b.Customer.Phone.String,
		},
		Payment: &v1.Payment{
			InvoiceNumber: b.InvoiceNumber.String,
		},
	}
}

type Customer struct {
	Name  string
	Email string
	Phone sql.NullString
}
