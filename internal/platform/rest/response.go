package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/HewlettPackard/hpecli/internal/platform/log"
)

type Response struct {
	*http.Response
	Error error
}

// Bytes returns the response's Body as []byte.
func (r *Response) Bytes() []byte {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	return body
}

// JSON returns the response's body as pretty printed []byte
func (r *Response) JSON() []byte {
	body := r.Bytes()
	dst := &bytes.Buffer{}

	if err := json.Indent(dst, body, "", "  "); err != nil {
		log.Logger.Debug("Unable to pretty-print output.")
		log.Logger.Debugf("body=%s", body)

		return body
	}

	return dst.Bytes()
}

// Unmarshall expected json response into passed interface type
func (r *Response) Unmarshall(t interface{}) error {
	if err := json.Unmarshal(r.Bytes(), &t); err != nil {
		log.Logger.Debugf("expcted login response, but received: %+v", t)
		return err
	}

	return nil
}
