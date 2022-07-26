package web

type Paging struct {
	Page     int32 `json:"page" form:"page" validate:"omitempty,min=1" example:"1"`
	PageSize int32 `json:"pageSize" form:"pageSize" validate:"omitempty,min=1" example:"10"`
}
