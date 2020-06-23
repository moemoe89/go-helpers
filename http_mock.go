package helpers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
)

type ctrl struct {
	statusCode int
	response   interface{}
}

func (c *ctrl) mockHandler(w http.ResponseWriter, r *http.Request) {
	resp := []byte{}

	rt := reflect.TypeOf(c.response)
	if rt.Kind() == reflect.String {
		resp = []byte(c.response.(string))
	} else if rt.Kind() == reflect.Struct || rt.Kind() == reflect.Ptr {
		resp, _ = json.Marshal(c.response)
	} else {
		resp = []byte("{}")
	}

	w.WriteHeader(c.statusCode)
	w.Write(resp)
}

func HttpMock(pattern string, statusCode int, response interface{}) *httptest.Server {
	c := &ctrl{statusCode, response}

	handler := http.NewServeMux()
	handler.HandleFunc(pattern, c.mockHandler)

	return httptest.NewServer(handler)
}
