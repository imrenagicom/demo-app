package catalog

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/imrenagicom/demo-app/internal/db"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

var (
	courseBatchKeyFmt = "course_batch:%s"
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

func (s *Store) FindAllCourse(ctx context.Context, opts ...ListOption) ([]Course, string, error) {
	options := &ListOptions{
		Limit: 10,
	}
	for _, o := range opts {
		o(options)
	}

	nextPage := pageToken{page: options.Page + 1}.encode()
	var courses []Course

	sb := sq.StatementBuilder.RunWith(s.dbCache)
	selectCourses := sb.
		Select("c.id", "c.name", "c.slug", "c.description", "c.status", "c.published_at").
		From("courses c").
		Where(sq.Eq{"c.deleted_at": nil, "c.status": CourseStatusPublished}).
		OrderBy("c.published_at DESC").
		Offset(uint64(options.GetOffset())).
		Limit(uint64(options.Limit)).
		PlaceholderFormat(sq.Dollar)

	rows, err := selectCourses.QueryContext(ctx)
	if err != nil {
		return nil, "", err
	}

	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.Status, &c.PublishedAt); err != nil {
			return nil, "", err
		}

		if options.Preload {
			batches, _, err := s.FindAllBatchesByCourseID(ctx, c.ID.String())
			if err != nil {
				return nil, "", err
			}
			c.Batches = batches
		}
		courses = append(courses, c)
	}
	return courses, nextPage, nil
}

func (s *Store) FindCourseByID(ctx context.Context, id string) (*Course, error) {
	c := Course{}
	sb := sq.StatementBuilder.RunWith(s.dbCache)
	getConcert := sb.
		Select("c.id", "c.name", "c.slug", "c.description", "c.status", "c.published_at").
		From("courses c").
		Where(sq.Eq{"c.deleted_at": nil, "c.id": id, "c.status": CourseStatusPublished}).
		PlaceholderFormat(sq.Dollar)
	if err := getConcert.QueryRowContext(ctx).Scan(
		&c.ID, &c.Name, &c.Slug, &c.Description, &c.Status, &c.PublishedAt,
	); err != nil {
		// if errors.Is(err, sql.ErrNoRows) {
		// 	return nil, db.ErrResourceNotFound{Message: fmt.Sprintf("course with id %s not found", id)}
		// }
		return nil, err
	}

	var batches []Batch
	selectBatches := sb.
		Select("id", "name", "max_seats", "available_seats", "price", "currency", "start_date", "end_date", "version").
		From("course_batches").
		Where(sq.Eq{"course_id": c.ID, "deleted_at": nil, "status": BatchStatusPublished}).
		PlaceholderFormat(sq.Dollar)
	rows, err := selectBatches.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var b Batch
		if err := rows.Scan(
			&b.ID, &b.Name, &b.MaxSeats, &b.AvailableSeats, &b.Price, &b.Currency, &b.StartDate, &b.EndDate, &b.Version,
		); err != nil {
			return nil, err
		}
		batches = append(batches, b)
	}
	c.Batches = batches
	return &c, nil
}

func (c *Store) CreateCourse(ctx context.Context, course *Course) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	sb := sq.StatementBuilder.RunWith(tx)
	insertCourse := sb.
		Insert("courses").
		Columns("id", "name", "slug", "description", "status", "published_at", "created_at", "updated_at").
		Values(course.ID, course.Name, course.Slug, course.Description, course.Status, course.PublishedAt, course.CreatedAt, course.UpdatedAt).
		PlaceholderFormat(sq.Dollar)

	insertBatches := sb.
		Insert("course_batches").
		Columns("id", "name", "max_seats", "available_seats", "price", "currency", "start_date", "end_date", "course_id", "created_at", "updated_at", "status").
		PlaceholderFormat(sq.Dollar)
	for _, b := range course.Batches {
		insertBatches = insertBatches.Values(b.ID, b.Name, b.MaxSeats, b.AvailableSeats, b.Price, b.Currency, b.StartDate, b.EndDate, course.ID, b.CreatedAt, b.UpdatedAt, b.Status)
	}

	_, err = insertCourse.ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = insertBatches.ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return err
}

func (c *Store) FindCourseBatchByID(ctx context.Context, id string, opts ...FindOption) (*Batch, error) {
	options := &FindOptions{}
	for _, o := range opts {
		o(options)
	}

	var b Batch
	sb := sq.StatementBuilder
	if options.Tx != nil {
		sb = sb.RunWith(options.Tx)
	} else {
		sb = sb.RunWith(c.dbCache)
	}

	selectBatch := sb.
		Select("id", "name", "max_seats", "available_seats", "price", "currency", "start_date", "end_date", "version", "status").
		From("course_batches").
		Where(sq.Eq{"id": id, "deleted_at": nil}).
		PlaceholderFormat(sq.Dollar)

	err := selectBatch.QueryRowContext(ctx).
		Scan(&b.ID, &b.Name, &b.MaxSeats, &b.AvailableSeats, &b.Price, &b.Currency, &b.StartDate, &b.EndDate, &b.Version, &b.Status)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (c *Store) FindCourseBatchByIDAndCourseID(ctx context.Context, batchID, courseID string, opts ...FindOption) (*Batch, error) {
	options := &FindOptions{}
	for _, o := range opts {
		o(options)
	}

	sb := sq.StatementBuilder
	if options.Tx != nil {
		sb = sb.RunWith(options.Tx)
	} else {
		sb = sb.RunWith(c.dbCache)
	}

	selectBatch := sb.
		Select("cb.id", "cb.name", "cb.max_seats", "cb.available_seats", "cb.price", "cb.currency", "cb.start_date", "cb.end_date", "cb.version", "cb.status").
		From("course_batches cb").
		Where(sq.Eq{"cb.id": batchID, "cb.course_id": courseID}).
		PlaceholderFormat(sq.Dollar)

	var b Batch
	err := selectBatch.QueryRowContext(ctx).
		Scan(&b.ID, &b.Name, &b.MaxSeats, &b.AvailableSeats, &b.Price, &b.Currency, &b.StartDate, &b.EndDate, &b.Version, &b.Status)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (c *Store) UpdateBatchAvailableSeats(ctx context.Context, b *Batch, opts ...UpdateOption) error {
	options := &UpdateOptions{}
	for _, o := range opts {
		o(options)
	}

	sb := sq.StatementBuilder
	if options.Tx != nil {
		sb = sb.RunWith(options.Tx)
	} else {
		sb = sb.RunWith(c.dbCache)
	}

	updateSeat := sb.
		Update("course_batches").
		Set("available_seats", b.AvailableSeats).
		Set("version", b.Version+1).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": b.ID, "version": b.Version}).
		PlaceholderFormat(sq.Dollar)

	res, err := updateSeat.ExecContext(ctx)
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

func (c *Store) FindAllBatchesByCourseID(ctx context.Context, courseID string, opts ...ListOption) ([]Batch, string, error) {	
	options := &ListOptions{
		Limit: 10,
	}
	for _, o := range opts {
		o(options)
	}

	nextPage := pageToken{page: options.Page + 1}.encode()
	var batches []Batch
	sb := sq.StatementBuilder.RunWith(c.dbCache)
	selectBatches := sb.
		Select("id", "name", "max_seats", "available_seats", "price", "currency", "start_date", "end_date", "version").
		From("course_batches").
		Where(sq.Eq{"course_id": courseID, "deleted_at": nil, "status": BatchStatusPublished}).
		OrderBy("created_at DESC").
		Offset(uint64(options.GetOffset())).
		Limit(uint64(options.Limit)).
		PlaceholderFormat(sq.Dollar)

	rows, err := selectBatches.QueryContext(ctx)
	if err != nil {
		return nil, "", err
	}

	for rows.Next() {
		var b Batch
		if err := rows.Scan(
			&b.ID, &b.Name, &b.MaxSeats, &b.AvailableSeats, &b.Price, &b.Currency, &b.StartDate, &b.EndDate, &b.Version,
		); err != nil {
			return nil, "", err
		}
		batches = append(batches, b)
	}
	return batches, nextPage, nil
}
