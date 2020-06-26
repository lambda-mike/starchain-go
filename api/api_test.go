package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHello(t *testing.T) {
	t.Log("TestHello")
	{
		RegisterHandlers()
		t.Log("\tWhen called at /hello")
		{
			req, err := http.NewRequest("GET", "/hello", nil)
			if err != nil {
				t.Fatalf("\t\tShould be able to create a request, err: %v", err)
			}
			t.Log("\t\tShould be able to create a request")
			recorder := httptest.NewRecorder()
			http.DefaultServeMux.ServeHTTP(recorder, req)
			if recorder.Code != 200 {
				t.Fatalf("\t\tShould get response 200 OK, got: %v", recorder.Code)
			}
			t.Log("\t\tShould get response 200 OK")
			body, err := ioutil.ReadAll(recorder.Body)
			if err != nil {
				t.Fatalf("\t\tShould not return err, got: %v", err)
			}
			t.Log("\t\tShould not return err")
			if string(body) != "hello" {
				t.Fatalf("\t\tShould return string \"hello\", got: \"%s\"", body)
			}
			t.Log("\t\tShould return string \"hello\"")
		}
	}
}
