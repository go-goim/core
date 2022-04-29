package util

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type PbJsonBinding struct{}

func (PbJsonBinding) Name() string {
	return "protobuf/json"
}

func (b PbJsonBinding) Bind(req *http.Request, obj interface{}) error {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return b.BindBody(buf, obj)
}

func (PbJsonBinding) BindBody(body []byte, obj interface{}) error {
	o := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	pbObj, ok := obj.(proto.Message)
	if !ok {
		return fmt.Errorf("%T is not a protobuf message", obj)
	}

	return o.Unmarshal(body, pbObj)
}

func MarshallPb(m proto.Message) ([]byte, error) {
	o := protojson.MarshalOptions{EmitUnpopulated: true}
	return o.Marshal(m)
}
