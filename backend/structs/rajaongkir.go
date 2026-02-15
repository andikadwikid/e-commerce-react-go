package structs

type RajaOngkirResponseWrapper struct {
	Meta struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"meta"`
	Data any `json:"data"`
}
