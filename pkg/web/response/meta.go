package response

import (
	"strconv"

	"github.com/go-goim/core/pkg/web"
)

type Meta struct {
	*web.Paging `json:",inline"`
	Total       int32             `json:"total"`
	Extra       map[string]string `json:"extra,omitempty"`
}

func (x *Meta) SetTotal(total int32) *Meta {
	x.Total = total

	return x
}

func (x *Meta) SetPaging(page, size int32) *Meta {
	x.Page = page
	x.PageSize = size

	return x
}

func (x *Meta) SetExtra(key, value string) *Meta {
	if x.Extra == nil {
		x.Extra = make(map[string]string)
	}

	x.Extra[key] = value

	return x
}

func (x *Meta) SetExtraInt(key string, value int) *Meta {
	if x.Extra == nil {
		x.Extra = make(map[string]string)
	}

	x.Extra[key] = strconv.Itoa(value)

	return x
}

func (x *Meta) SetExtraMap(m map[string]string) *Meta {
	if x.Extra == nil {
		x.Extra = make(map[string]string)
	}

	for k, v := range m {
		x.Extra[k] = v
	}

	return x
}
