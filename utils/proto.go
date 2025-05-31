package utils

import (
	"google.golang.org/protobuf/proto"
)

func DeepCopyProto[T proto.Message](msg T) T {
	return proto.Clone(msg).(T)
}
