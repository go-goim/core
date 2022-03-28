package util

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func MarshallPb(m proto.Message) ([]byte, error) {
	o := protojson.MarshalOptions{EmitUnpopulated: true}
	return o.Marshal(m)
}
