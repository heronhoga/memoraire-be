package requests

type RequestCreateMemo struct {
	Date string `json:"date" validate:"required"`
	Note string `json:"note" validate:"required"` 
}