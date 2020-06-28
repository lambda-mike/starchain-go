package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/starchain/contracts"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type BlockchainMock struct{}

func (b BlockchainMock) RequestMessageOwnershipVerification(addr string) (string, error) {
	return addr + " OK", nil
}

func (b BlockchainMock) GetBlockByHeight(h int) (contracts.Block, error) {
	var block contracts.Block
	return block, errors.New("TODO")
}

func createApi() *httptest.Server {
	var blockchain contracts.Blockchain = BlockchainMock{}
	api := Create(&blockchain)
	return httptest.NewServer(api)
}

func TestHello(t *testing.T) {
	t.Log("Hello")
	{
		server := createApi()
		defer server.Close()
		t.Log("Server url: ", server.URL)
		t.Log("\tWhen called at /hello")
		{
			resp, err := http.Get(server.URL + "/hello")
			if err != nil {
				t.Fatalf("\t\tShould be able to create a request, err: %v", err)
			}
			t.Log("\t\tShould be able to create a request")
			if resp.StatusCode != 200 {
				t.Fatalf("\t\tShould get response 200 OK, got: %v", resp.StatusCode)
			}
			t.Log("\t\tShould get response 200 OK")
			body, err := ioutil.ReadAll(resp.Body)
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

func TestRequestValidation(t *testing.T) {
	t.Log("RequestValidation")
	{
		server := createApi()
		defer server.Close()
		t.Log("Server url: ", server.URL)
		t.Log("\tWhen called at /requestValidation")
		{
			addr := "1FzpnkhbAteDkU1wXDtd8kKizQhqWcsrWe"
			data, _ := json.Marshal(Address{addr})
			response, err := http.Post(server.URL+"/requestValidation", "application/json", bytes.NewReader(data))
			if err != nil {
				t.Fatalf("\t\tShould be able to submit a validation request, got err: %v", err)
			}
			t.Log("\t\tShould be able to submit a validation request")
			if response.StatusCode != 200 {
				t.Fatalf("\t\tShould get response 200 OK, got: %v", response.StatusCode)
			}
			t.Log("\t\tShould get response 200 OK")
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("\t\tShould not return err, got: %v", err)
			}
			t.Log("\t\tShould not return err")
			expected := addr + " OK"
			if string(body) != expected {
				t.Fatalf("\t\tShould return correct message: \"%s\", got: \"%s\"", expected, body)
			}
			t.Log("\t\tShould return correct message")
		}
	}
}

func TestGetBlockByHeight(t *testing.T) {
	t.Log("GetBlockByHeight")
	{
		server := createApi()
		defer server.Close()
		t.Log("Server url: ", server.URL)
		t.Log("\tGiven a need to test endpoint /block/:height")
		{
			t.Log("\tWhen called with 0 index")
			{
				response, err := http.Get(server.URL + "/block/0")
				if err != nil {
					t.Fatalf("\t\tShould be able to get a block, got err: %v", err)
				}
				t.Log("\t\tShould be able to get a block")
				if response.StatusCode != 200 {
					t.Fatalf("\t\tShould get response 200 OK, got: %v", response.StatusCode)
				}
				t.Log("\t\tShould get response 200 OK")
				// TODO add block to contracts
				// TODO decode body to block
				err = errors.New("TODO")
				if err != nil {
					t.Fatalf("\t\tShould not return err, got: %v", err)
				}
				t.Log("\t\tShould not return err")
			}
		}
	}
}
