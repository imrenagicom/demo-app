package db

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNoRowUpdated = errors.New("no row is updated")
)

type ErrResourceNotFound struct {
	Message string
}

func (e ErrResourceNotFound) Error() string {
	return e.Message
}

func (e ErrResourceNotFound) GRPCStatus() *status.Status {
	return status.New(codes.NotFound, e.Error())
}
