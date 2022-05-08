// Code Written Manually

package response

import (
	"strconv"
)

// NewMeta returns a new Meta object
func NewMeta() *Meta {
	return &Meta{}
}

func (x *Meta) SetRequestID(id string) *Meta {
	x.RequestId = id

	return x
}

func (x *Meta) SetTotal(total int) *Meta {
	var t = int32(total)
	x.Total = &t

	return x
}

func (x *Meta) SetPaging(page, size int) *Meta {
	if x.Pagination == nil {
		x.Pagination = &Pagination{}
	}

	x.Pagination.Page = int32(page)
	x.Pagination.PageSize = int32(size)

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

func (x *Meta) Merge(src *Meta) *Meta {
	if src == nil {
		return x
	}

	if src.RequestId != "" {
		x.RequestId = src.RequestId
	}

	if src.Total != nil {
		x.Total = src.Total
	}

	if src.Pagination != nil {
		x.Pagination = src.Pagination
	}

	if src.Extra != nil {
		x.SetExtraMap(src.Extra)
	}

	return x
}
