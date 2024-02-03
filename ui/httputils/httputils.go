package httputils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

const (
	NilAccessToken = "access token is nill"

	HeaderContentType                = "Content-Type"
	HeaderApplicationJsonContentType = "application/json"
	HeaderAuthorization              = "Authorization"

	RequestTypePost = "POST"
	RequestTypePut  = "PUT"
	RequestTypeGet  = "GET"
)

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: make(map[string]string),
	}
}

func (headers *Headers) Add(key string, val string) *Headers {
	if len(key) > 0 {
		headers.headers[key] = val
	}
	return headers
}

func (headers *Headers) WithJsonContentTypeHeader() *Headers {
	headers.headers[HeaderContentType] = HeaderApplicationJsonContentType
	return headers
}

func (headers *Headers) WithAuthorizationHeader(authValue string) *Headers {
	if len(authValue) > 0 {
		headers.headers[HeaderAuthorization] = authValue
	}

	return headers
}

type Queries struct {
	queries map[string]string
}

func NewQueries() *Queries {
	return &Queries{
		queries: make(map[string]string),
	}
}

func (queries *Queries) Add(key string, val string) *Queries {
	if len(key) > 0 {
		queries.queries[key] = val
	}
	return queries
}

func SendGETRequest(
	url string,
	queries *Queries,
	headers *Headers) ([]byte, int, error) {

	return SendRequest(RequestTypeGet, url, headers, queries, nil)
}

func SendPOSTRequest(
	url string,
	headers *Headers,
	postData []byte) (response []byte, statusCode int, err error) {

	return SendRequest(RequestTypePost, url, headers, nil, postData)
}

func SendRequest(
	requestType string,
	url string,
	headers *Headers,
	queries *Queries,
	data []byte) (response []byte, statusCode int, err error) {

	if queries != nil && len(queries.queries) > 0 {
		url += "?"
		for qp, qv := range queries.queries {
			url += fmt.Sprintf("%s=%s", qp, qv)
		}
	}

	httpRq, err := http.NewRequest(requestType, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("error creating post request: %w", err)
	}

	if headers != nil {
		for hk, hv := range headers.headers {
			httpRq.Header.Set(hk, hv)
		}
	}
	client := &http.Client{}
	httpRs, err := client.Do(httpRq)
	if err != nil {
		if httpRs != nil {
			return nil, httpRs.StatusCode, fmt.Errorf("error in %s request: %w", requestType, err)
		} else {
			return nil, http.StatusInternalServerError, fmt.Errorf("error in %s request: %w", requestType, err)
		}
	}
	if httpRs.StatusCode != http.StatusOK {
		return nil, httpRs.StatusCode, fmt.Errorf("returned status not OK: %d", httpRs.StatusCode)
	}
	defer httpRs.Body.Close()

	response, err = io.ReadAll(httpRs.Body)
	if err != nil {
		return nil, httpRs.StatusCode, fmt.Errorf("error reading response from body: %w", err)
	}

	return response, httpRs.StatusCode, nil
}
