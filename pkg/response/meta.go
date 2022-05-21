package response

import (
	"strconv"
)

type Meta struct {
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Extra    map[string]string `json:"extra,omitempty"`
}

func (x *Meta) SetTotal(total int) *Meta {
	x.Total = total

	return x
}

func (x *Meta) SetPaging(page, size int) *Meta {
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
