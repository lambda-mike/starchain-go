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

var mockBlocks [1]contracts.Block = [...]contracts.Block{
	contracts.Block{"Genesis Block", "123abc456", 0, "abcdef123", "", 1592156792},
}

type BlockchainMock struct{}

func (b BlockchainMock) RequestMessageOwnershipVerification(addr string) (string, error) {
	return addr + " OK", nil
}

func (b BlockchainMock) GetBlockByHeight(h int) (contracts.Block, error) {
	var block contracts.Block
	switch h {
	case 0:
		return mockBlocks[0], nil
	case 1:
		return block, errors.New("TODO")
	default:
		return block, errors.New("TODO")
	}
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
			data, _ := json.Marshal(AddressDto{addr})
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
			t.Log("\tWhen called with proper index")
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
				var block BlockDto
				if err := json.NewDecoder(response.Body).Decode(&block); err != nil {
					t.Fatalf("\t\tShould decode response body, got err: %v; json: %v", err, response.Body)
				}
				t.Log("\t\tShould decode response body")
				if block.Body != mockBlocks[0].Body ||
					block.Time != mockBlocks[0].Time ||
					block.Hash != mockBlocks[0].Hash ||
					block.Owner != mockBlocks[0].Owner ||
					block.PreviousBlockHash != mockBlocks[0].PreviousBlockHash ||
					block.Height != mockBlocks[0].Height {
					t.Fatalf("\t\tShould return genesis block, got: %v", block)
				}
				t.Log("\t\tShould return genesis block")
			}
			t.Log("\tWhen called with wrong index")
			{
				response, err := http.Get(server.URL + "/block/666")
				if err != nil {
					t.Fatal("\t\tShould not get an error for non existing block: ", err)
				}
				t.Log("\t\tShould not get an error for non existing block")
				if response.StatusCode != 404 {
					t.Fatal("\t\tShould return not found status code, got: ", response.StatusCode)
				}
				t.Log("\t\tShould return not found status code")
			}
		}
	}
}

func TestGetBlockByHash(t *testing.T) {
	t.Log("GetBlockByHash")
	{
		server := createApi()
		defer server.Close()
		t.Log("Server url: ", server.URL)
		t.Log("\tGiven a need to test endpoint /block/hash/:hash")
		{
			t.Log("\tWhen called with proper hash")
			{
				response, err := http.Get(server.URL + "/block/hash/" + mockBlocks[0].Hash)
				if err != nil {
					t.Fatalf("\t\tShould be able to get a block, got err: %v", err)
				}
				t.Log("\t\tShould be able to get a block")
				if response.StatusCode != 200 {
					t.Fatalf("\t\tShould get response 200 OK, got: %v", response.StatusCode)
				}
				t.Log("\t\tShould get response 200 OK")
				var block BlockDto
				if err := json.NewDecoder(response.Body).Decode(&block); err != nil {
					t.Fatalf("\t\tShould decode response body, got err: %v; json: %v", err, response.Body)
				}
				t.Log("\t\tShould decode response body")
				if block.Body != mockBlocks[0].Body ||
					block.Time != mockBlocks[0].Time ||
					block.Hash != mockBlocks[0].Hash ||
					block.Owner != mockBlocks[0].Owner ||
					block.PreviousBlockHash != mockBlocks[0].PreviousBlockHash ||
					block.Height != mockBlocks[0].Height {
					t.Fatalf("\t\tShould return genesis block, got: %v", block)
				}
				t.Log("\t\tShould return genesis block")
			}
			t.Log("\tWhen called with wrong hash")
			{
				response, err := http.Get(server.URL + "/block/hash/666")
				if err != nil {
					t.Fatal("\t\tShould not get an error for non existing block: ", err)
				}
				t.Log("\t\tShould not get an error for non existing block")
				if response.StatusCode != 404 {
					t.Fatal("\t\tShould return not found status code, got: ", response.StatusCode)
				}
				t.Log("\t\tShould return not found status code")
			}
		}
	}
}
