// Helper package to handle drudgery of REST request handling
package rest

import (
	"io"
	"net/http"
)

type Request struct {
	*http.Request
	*http.Client
	*http.Transport
}

func wrapRequest(method, urlStr string, body io.Reader, options []func(*Request)) (*Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}

	request := &Request{
		Request:   req,
		Client:    &http.Client{},
		Transport: &http.Transport{},
	}

	// Apply options in the parameters to request.
	for _, option := range options {
		option(request)
	}

	return request, nil
}

// Get sends a HTTP GET request to the provided url.
// options are run prior to the request being sent
//
//     addAuth := func(r *Request) {
//             r.SetBasicAuth("username", "password")
//     }
//
//     resp, err := requests.Get("http://ilo.com/redfish", addAuth)
//     if err != nil {
//             panic(err)
//     }
//     fmt.Println(resp.StatusCode)
//
func Get(urlStr string, options ...func(*Request)) (*Response, error) {
	return do("GET", urlStr, nil, options...)
}

// Post sends a HTTP POST request to the provided url.
// and option should be include to set the correct content-type header
//
//     addMimeType := func(r *Request) {
//             r.Header.Add("content-type", "application/json")
//     }
// resp, err := requests.Post("https://ilo.com/login", buf, redirect)
// if err != nil {
//         panic(err)
// }
// fmt.Println(resp.JSON())
//
func Post(urlStr string, body io.Reader, options ...func(*Request)) (*Response, error) {
	return do("POST", urlStr, body, options...)
}

func do(method, urlStr string, body io.Reader, options ...func(*Request)) (*Response, error) {
	request, err := wrapRequest("POST", urlStr, body, options)
	if err != nil {
		return nil, err
	}

	//nolint:bodyclose // body is closed in the response helper calls
	resp, err := request.Client.Do(request.Request)
	if err != nil {
		return nil, err
	}

	return &Response{Response: resp}, nil
}
