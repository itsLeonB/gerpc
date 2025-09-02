package gerpc

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func NullableTimeToProto(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}
