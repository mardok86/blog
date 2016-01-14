package article

type BlogError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
