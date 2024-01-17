package proto

import (
	"errors"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var ErrInvalidPayload = errors.New("invalid payload")

func Unmarshal(b []byte, v proto.Message) error {
	data := &anypb.Any{}
	err := protojson.Unmarshal(b, data)
	if err != nil {
		return ErrInvalidPayload
	}
	err = anypb.UnmarshalTo(data, v, proto.UnmarshalOptions{})
	if err != nil {
		return ErrInvalidPayload
	}
	return nil
}

func Marshal(v proto.Message) ([]byte, error) {
	data, err := anypb.New(v)
	if err != nil {
		return nil, err
	}
	return protojson.Marshal(data)
}
