package catalog

import (
	"context"

	"github.com/imrenagicom/demo-app/course/catalog"
	v1 "github.com/imrenagicom/demo-app/pkg/apiclient/course/v1"
)

type Service interface {
	ListCourse(ctx context.Context, req *v1.ListCoursesRequest) ([]catalog.Course, string, error)
	GetCourse(ctx context.Context, req *v1.GetCourseRequest) (*catalog.Course, error)
}

func New(s Service) *Server {
	return &Server{
		service: s,
	}
}

type Server struct {
	v1.UnimplementedCatalogServiceServer

	service Service
}

func (s Server) ListCourses(ctx context.Context, req *v1.ListCoursesRequest) (*v1.ListCoursesResponse, error) {
	courses, nextPage, err := s.service.ListCourse(ctx, req)
	if err != nil {
		return nil, err
	}

	var data []*v1.Course
	for _, c := range courses {
		data = append(data, c.ApiV1())
	}

	res := &v1.ListCoursesResponse{
		Courses:       data,
		NextPageToken: nextPage,
	}
	return res, nil
}

func (s Server) GetCourse(ctx context.Context, req *v1.GetCourseRequest) (*v1.Course, error) {
	course, err := s.service.GetCourse(ctx, req)
	if err != nil {
		return nil, err
	}
	return course.ApiV1(), nil
}
