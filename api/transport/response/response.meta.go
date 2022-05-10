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
	var (
		p = int32(page)
		s = int32(size)
	)

	x.Page = &p
	x.PageSize = &s

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

	if src.Page != nil {
		// page and size are set in the same time
		x.Page = src.Page
		x.PageSize = src.PageSize
	}

	if src.Extra != nil {
		x.SetExtraMap(src.Extra)
	}

	return x
}
