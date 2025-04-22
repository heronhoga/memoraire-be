package requests

type RequestCreateMemo struct {
	Date string `json:"date" validate:"required"`
	Note string `json:"note" validate:"required"` 
}

type RequestUpdateMemo struct {
	MemoId	string `json:"memo_id" validate:"required"`
	Date 	string `json:"date" validate:"required"`
	Note 	string `json:"note" validate:"required"` 
}

type RequestDeleteMemo struct {
	MemoId string `json:"memo_id" validate:"required"`
}