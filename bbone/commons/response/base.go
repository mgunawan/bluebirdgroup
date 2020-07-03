package response

//BaseResponse ...
type BaseResponse struct {
	Status int         `json:"code"`
	Result interface{} `json:"result"`
}
