package types

type WSResponse struct {
	Method string      `json:"method"`
	ID     int         `json:"id"`
	Code   int         `json:"code"`
	Result interface{} `json:"result"`
}