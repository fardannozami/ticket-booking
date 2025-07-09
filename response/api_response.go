package response

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type CreatedResponse struct {
	Code    int    `json:"code" example:"201"`
	Message string `json:"message" example:"created"`
}
