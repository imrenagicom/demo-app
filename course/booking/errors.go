package booking

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrReservationMaxRetryExceeded = errors.New("reservation max retry exceeded")
	ErrReleaseMaxRetryExceeded     = errors.New("booking release max retry exceeded")

	ErrBookingAlreadyExpired   = errors.New("booking already expired")
	ErrBookingAlreadyCompleted = ErrInvalidStateChange{Message: "booking already completed"}
)

type ErrInvalidStateChange struct {
	Message string
}

func (e ErrInvalidStateChange) Error() string {
	return e.Message
}

func (e ErrInvalidStateChange) GRPCStatus() *status.Status {
	return status.New(codes.FailedPrecondition, e.Error())
}
