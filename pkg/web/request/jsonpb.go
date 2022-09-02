package request

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin/binding"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	responsepb "github.com/go-goim/api/transport/response"

	"github.com/go-goim/core/pkg/mid"
)

type PbJSONBinding struct {
	validate bool
}

var (
	_ binding.Binding = PbJSONBinding{}
)

var (
	// NonValidatePbJSONBinding is a PbJSONBinding without validation.
	NonValidatePbJSONBinding = PbJSONBinding{
		validate: false,
	}
	// ValidatePbJSONBinding is a PbJSONBinding with validation.
	ValidatePbJSONBinding = PbJSONBinding{
		validate: true,
	}
)

func (b PbJSONBinding) Name() string {
	return "protobuf/json"
}

func (b PbJSONBinding) Bind(req *http.Request, obj interface{}) error {
	buf, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return b.BindBody(buf, obj)
}

func (b PbJSONBinding) BindBody(body []byte, obj interface{}) error {
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

	if !b.validate {
		return nil
	}

	validate, ok := pbObj.(mid.Validator)
	if !ok {
		return nil
	}

	err := validate.Validate()
	if err != nil {
		return responsepb.NewBaseResponseWithMessage(responsepb.Code_InvalidParams, err.Error())
	}

	return nil
}

func MarshallPb(m proto.Message) ([]byte, error) {
	o := protojson.MarshalOptions{EmitUnpopulated: true, UseEnumNumbers: true}
	return o.Marshal(m)
}
