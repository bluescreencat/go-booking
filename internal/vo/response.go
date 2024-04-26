package vo

type Response struct {
	Pagination   Pagination  `json:"pagination,omitempty"`
	Success      bool        `json:"success"`
	ErrorMessage string      `json:"errorMessage,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	Errors       []any       `json:"errors,omitempty"`
}

func (r *Response) SetData(data interface{}) {
	r.Success = true
	r.Data = data
}

func (r *Response) SetDataWithPagination(data interface{}, page uint, size uint, total uint) {
	r.Success = true
	r.Data = data
	r.Pagination.Page = page
	r.Pagination.Size = size
	r.Pagination.Total = total
}

func (r *Response) SetErrorMessage(errorMessage string) {
	r.Success = false
	r.ErrorMessage = errorMessage
}

func (r *Response) AppendErrors(err interface{}) {
	r.Success = false
	r.Errors = append(r.Errors, err)
}

func (r *Response) SetErrorsWithDefaultData(data any, errs *[]interface{}) {
	r.Success = false
	r.Data = data
	r.Errors = *errs
}
