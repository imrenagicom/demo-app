package catalog

import (
	"database/sql"
	"time"

	v1 "github.com/imrenagicom/demo-app/pkg/apiclient/course/v1"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CourseStatus int

const (
	CourseStatusDraft CourseStatus = iota
	CourseStatusPublished
	CourseStatusArchived
)

type Course struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime
	Name        string
	Slug        string
	Description string
	PublishedAt sql.NullTime
	Batches     []Batch
	Status      CourseStatus
}

func (c Course) ApiV1() *v1.Course {
	var publishedAt *timestamppb.Timestamp
	if c.PublishedAt.Valid {
		publishedAt = timestamppb.New(c.PublishedAt.Time)
	}

	return &v1.Course{
		Name:        c.Slug,
		CourseId:    c.ID.String(),
		DisplayName: c.Name,
		Description: c.Description,
		PublishedAt: publishedAt,
		Batches:     c.batchesPkg(),
	}
}

func (c Course) batchesPkg() []*v1.Batch {
	var bs []*v1.Batch
	for _, b := range c.Batches {
		bs = append(bs, b.ApiV1())
	}
	return bs
}
