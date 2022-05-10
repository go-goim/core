package request

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	responsepb "github.com/yusank/goim/api/transport/response"
	"github.com/yusank/goim/pkg/mid"
)

type PbJSONBinding struct{}

func (PbJSONBinding) Name() string {
	return "protobuf/json"
}

func (b PbJSONBinding) Bind(req *http.Request, obj interface{}) error {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return b.BindBody(buf, obj)
}

func (PbJSONBinding) BindBody(body []byte, obj interface{}) error {
	o := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	pbObj, ok := obj.(proto.Message)
	if !ok {
		return fmt.Errorf("%T is not a protobuf message", obj)
	}

	if err := o.Unmarshal(body, pbObj); err != nil {
		return err
	}

	validate, ok := pbObj.(mid.Validator)
	if !ok {
		return nil
	}

	err := validate.Validate()
	if err != nil {
		return responsepb.NewBaseResponse(responsepb.Code_InvalidParams, err.Error())
	}

	return nil
}

func MarshallPb(m proto.Message) ([]byte, error) {
	o := protojson.MarshalOptions{EmitUnpopulated: true, UseEnumNumbers: true}
	return o.Marshal(m)
}
