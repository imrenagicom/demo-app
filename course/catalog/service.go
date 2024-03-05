package catalog

import (
	"context"
	"database/sql"
	"math/rand"
	"time"

	v1 "github.com/imrenagicom/demo-app/pkg/apiclient/course/v1"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func NewService(store *Store, db *sqlx.DB) *Service {
	return &Service{
		db:    db,
		store: store,
	}
}

type Service struct {
	db    *sqlx.DB
	store *Store
}

func (s Service) ListCourse(ctx context.Context, req *v1.ListCoursesRequest) ([]Course, string, error) {
	opts := []ListOption{
		WithMaxResults(req.GetPageSize()),
		WithNextPage(req.GetPageToken()),
	}
	for _, f := range req.GetListMask().GetPaths() {
		switch f {
		case "courses.batches":
			opts = append(opts, WithPreload())
		}
	}
	return s.store.FindAllCourse(ctx, opts...)
}

func (s Service) GetCourse(ctx context.Context, req *v1.GetCourseRequest) (*Course, error) {
	return s.store.FindCourseByID(ctx, req.GetCourse())
}

func (s Service) Seed(ctx context.Context) error {
	for i := 0; i < 1000; i++ {
		start := time.Now()
		end := start.AddDate(0, 2, 0)

		// generate random number of batches
		numBatches := rand.Intn(5) + 1
		var batches []Batch
		for j := 0; j < numBatches; j++ {
			tc := Batch{
				ID:             uuid.New(),
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
				Name:           faker.Name(),
				MaxSeats:       50,
				AvailableSeats: int32(rand.Intn(50)),
				Status:         BatchStatusPublished,
				Price:          float64(rand.Intn(100000)) + 100000,
				Currency:       "IDR",
				StartDate:      sql.NullTime{Time: start, Valid: true},
				EndDate:        sql.NullTime{Time: end, Valid: true},
			}
			batches = append(batches, tc)
		}

		c := &Course{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Name:        faker.Name(),
			Slug:        faker.Username(),
			Description: faker.Paragraph(),
			PublishedAt: sql.NullTime{
				Time:  start,
				Valid: true,
			},
			Status:  CourseStatusPublished,
			Batches: batches,
		}

		err := s.store.CreateCourse(ctx, c)
		if err != nil {
			return err
		}
	}
	return nil
}
