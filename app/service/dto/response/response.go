package response

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

// ErrorResponseData -
type ErrorResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse -
type ErrorResponse struct {
	Success bool              `json:"success" default:"false"`
	Error   ErrorResponseData `json:"data"`
}

type Response struct {
	Success bool `json:"success"`
	Data    Data `json:"data"`
}

type ResponseV2 struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Request interface{} `json:"request,omitempty"`
	// List    interface{} `json:"list,omitempty"`
}

type ResponseV3 struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Request interface{} `json:"request,omitempty"`
}

type Pagination struct {
	Page       int `json:"page" default:"1"`
	PerPage    int `json:"per_page" default:"10"`
	TotalPages int `json:"total_pages"`
	TotalCount int `json:"total_count"`
}
type Data struct {
	Message string `json:"message"`
}
