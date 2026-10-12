package models

type APIRequest struct {
	Headers               map[string]string
	QueryStringParameters map[string]string
	Body                  string
	IsBase64Encoded       bool
}

type APIResponse struct {
	StatusCode      int
	Body            string
	Headers         map[string]string
	IsBase64Encoded bool
}
