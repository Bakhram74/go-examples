package httpclient

const (
	ApplicationJson ContentType = "application/json"
)

type (
	ContentType string

	IHttpClient interface {
		Get(
			url string,
			body *SwClientRequest,
		) (*SwClientResponse, error)
		Post(
			url string,
			body *SwClientRequest,
		) (*SwClientResponse, error)
	}

	SwClientResponse struct {
		Status  int // Должен быть только из пакета net/http
		Payload []byte
	}

	SwClientRequestBody struct {
		ContentType ContentType
		Payload     []byte
	}

	SwClientRequest struct {
		Body    *SwClientRequestBody
		Params  map[string]string
		Headers map[string]string
	}

	RequestOpts func(r *SwClientRequest)
)

func NewRequest(opts ...RequestOpts) *SwClientRequest {
	req := &SwClientRequest{
		Body:    nil,
		Params:  nil,
		Headers: nil,
	}

	for _, opt := range opts {
		opt(req)
	}

	return req
}

func WithJsonRequestBody(p []byte) RequestOpts {
	return func(r *SwClientRequest) {
		r.Body = &SwClientRequestBody{
			Payload:     p,
			ContentType: ApplicationJson,
		}
	}
}

func WithRequestParams(params map[string]string) RequestOpts {
	return func(r *SwClientRequest) {
		r.Params = params
	}
}

func WithRequestHeaders(headers map[string]string) RequestOpts {
	return func(r *SwClientRequest) {
		r.Headers = headers
	}
}
