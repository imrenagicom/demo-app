package booking

import (
	"context"
	"math/rand"
	"time"

	"github.com/imrenagicom/demo-app/course/catalog"
	"github.com/imrenagicom/demo-app/internal/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

var (
	bookingTTL = 10 * time.Minute
)

func NewStore(db *sqlx.DB, redis redis.UniversalClient) *Store {
	return &Store{
		db:      db,
		dbCache: sq.NewStmtCache(db),
		redis:   redis,
	}
}

type Store struct {
	db      *sqlx.DB
	dbCache *sq.StmtCache
	redis   redis.UniversalClient
}

func (s *Store) Clear() error {
	return s.dbCache.Clear()
}

func (s *Store) CreateBooking(ctx context.Context, booking *Booking, opts ...CreateOption) error {
	options := &CreateOptions{}
	for _, o := range opts {
		o(options)
	}

	sb := sq.StatementBuilder.RunWith(s.dbCache)
	if options.Tx != nil {
		sb = sb.RunWith(options.Tx)
	}
	insertBooking := sb.Insert("bookings").
		Columns("id", "course_id", "course_batch_id", "price", "currency", "status", "created_at", "updated_at", "cust_name", "cust_email", "cust_phone").
		Values(booking.ID, booking.Course.ID, booking.Batch.ID,
			booking.Price, booking.Currency, booking.Status,
			booking.CreatedAt, booking.UpdatedAt, booking.Customer.Name, booking.Customer.Email, booking.Customer.Phone).
		PlaceholderFormat(sq.Dollar)

	_, err := insertBooking.ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) FindBookingByID(ctx context.Context, ID string, opts ...FindOption) (*Booking, error) {
	options := &FindOptions{}
	for _, o := range opts {
		o(options)
	}

	var b Booking = Booking{
		Course:   &catalog.Course{},
		Batch:    &catalog.Batch{},
		Customer: Customer{},
	}

	sb := sq.StatementBuilder.RunWith(s.dbCache)
	if options.Tx != nil {
		sb = sb.RunWith(options.Tx)
	}
	query := sb.Select("b.id", "c.id", "cb.id", "b.price", "b.currency", "b.status",
		"b.reserved_at", "b.expired_at", "b.paid_at", "b.created_at", "b.updated_at", "b.version",
		"b.cust_name", "b.cust_email", "b.cust_phone", "b.invoice_number", "b.payment_type",
		"c.name", "c.slug", "cb.name", "cb.start_date", "cb.end_date").
		From("bookings b").
		LeftJoin("courses c ON b.course_id = c.id").
		LeftJoin("course_batches cb ON b.course_batch_id = cb.id").
		Where(sq.Eq{"b.id": ID, "b.deleted_at": nil}).
		PlaceholderFormat(sq.Dollar)

	err := query.QueryRowContext(ctx).
		Scan(&b.ID, &b.Course.ID, &b.Batch.ID, &b.Price, &b.Currency, &b.Status,
			&b.ReservedAt, &b.ExpiredAt, &b.PaidAt, &b.CreatedAt, &b.UpdatedAt, &b.Version,
			&b.Customer.Name, &b.Customer.Email, &b.Customer.Phone, &b.InvoiceNumber, &b.PaymentType,
			&b.Course.Name, &b.Course.Slug, &b.Batch.Name, &b.Batch.StartDate, &b.Batch.EndDate)
	if err != nil {
		return nil, err
	}

	if rand.Intn(5)+1 == 3 {
		<-time.After(time.Duration(rand.Intn(300)) * time.Millisecond)
	}

	return &b, nil
}

func (s *Store) UpdateBookingStatus(ctx context.Context, booking *Booking, opts ...UpdateOption) error {
	options := &UpdateOptions{}
	for _, o := range opts {
		o(options)
	}

	sb := sq.StatementBuilder.RunWith(s.dbCache)
	if options.Tx != nil {
		sb = sb.RunWith(options.Tx)
	}
	updateBooking := sb.Update("bookings").
		Set("reserved_at", booking.ReservedAt).
		Set("expired_at", booking.ExpiredAt).
		Set("paid_at", booking.PaidAt).
		Set("status", booking.Status).
		Set("invoice_number", booking.InvoiceNumber).
		Set("version", booking.Version+1).
		Where(sq.Eq{"id": booking.ID, "version": booking.Version}).
		PlaceholderFormat(sq.Dollar)
	res, err := updateBooking.ExecContext(ctx)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return db.ErrNoRowUpdated
	}
	return nil
}

func (s *Store) UpdateBookingPayment(ctx context.Context, booking *Booking, opts ...UpdateOption) error {
	options := &UpdateOptions{}
	for _, o := range opts {
		o(options)
	}
	sb := sq.StatementBuilder.RunWith(s.dbCache)
	if options.Tx != nil {
		sb = sb.RunWith(options.Tx)
	}
	updateBooking := sb.Update("bookings").
		Set("paid_at", booking.PaidAt).
		Set("invoice_number", booking.InvoiceNumber).
		Set("payment_type", booking.PaymentType).
		Set("version", booking.Version+1).
		Where(sq.Eq{"id": booking.ID, "version": booking.Version}).
		PlaceholderFormat(sq.Dollar)
	res, err := updateBooking.ExecContext(ctx)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return nil
	}
	return nil
}

func (s *Store) FindAllBookings(ctx context.Context, opts ...ListOption) ([]Booking, string, error) {
	options := &ListOptions{
		Limit: 5,
	}
	for _, o := range opts {
		o(options)
	}

	sb := sq.StatementBuilder.RunWith(s.dbCache)
	if options.Tx != nil {
		sb = sb.RunWith(options.Tx)
	}

	var filter map[string]interface{} = map[string]interface{}{
		"b.deleted_at": nil,
	}
	if options.Status != 0 {
		filter["b.status"] = options.Status
	}
	if options.InvoiceNumber != "" {
		filter["b.invoice_number"] = options.InvoiceNumber
	}
	query := sb.Select("b.id", "c.id", "cb.id", "b.price", "b.currency", "b.status",
		"b.reserved_at", "b.expired_at", "b.paid_at", "b.created_at", "b.updated_at", "b.version",
		"b.cust_name", "b.cust_email", "b.cust_phone", "b.invoice_number", "b.payment_type",
		"c.name", "c.slug", "cb.name", "cb.start_date", "cb.end_date").
		From("bookings b").
		LeftJoin("courses c ON b.course_id = c.id").
		LeftJoin("course_batches cb ON b.course_batch_id = cb.id").
		Where(filter).
		Offset(uint64(options.GetOffset())).
		Limit(uint64(options.Limit)).
		PlaceholderFormat(sq.Dollar)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, "", err
	}

	var bookings []Booking
	for rows.Next() {
		var b Booking = Booking{
			Course:   &catalog.Course{},
			Batch:    &catalog.Batch{},
			Customer: Customer{},
		}
		if err := rows.
			Scan(&b.ID, &b.Course.ID, &b.Batch.ID, &b.Price, &b.Currency, &b.Status,
				&b.ReservedAt, &b.ExpiredAt, &b.PaidAt, &b.CreatedAt, &b.UpdatedAt, &b.Version,
				&b.Customer.Name, &b.Customer.Email, &b.Customer.Phone, &b.InvoiceNumber, &b.PaymentType,
				&b.Course.Name, &b.Course.Slug, &b.Batch.Name, &b.Batch.StartDate, &b.Batch.EndDate); err != nil {
			return nil, "", err
		}
		bookings = append(bookings, b)
	}
	return bookings, "", nil
}

func bookingCacheKey(id string) string {
	return "booking:" + id
}

func ttl(dur time.Duration) time.Duration {
	return dur + time.Duration(rand.Intn(5)+1)*time.Second // add jitter
}
