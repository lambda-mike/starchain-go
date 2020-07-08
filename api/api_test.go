package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/starchain/contracts"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// We store raw JSON as data (it comes from http request)
var mockBlocks [4]contracts.Block = [...]contracts.Block{
	contracts.Block{`"Genesis Block"`, "123abc456", 0, "", "", 1592156792},
	contracts.Block{`"Regular Block"`, "789abc987", 0, "7a7b7c", "123abc456", 1592156794},
	contracts.Block{`"Other Block"`, "fff333", 0, "333fff", "789abc987", 1592156795},
	contracts.Block{`"Regular Block II"`, "789abc987", 0, "7a7b7c", "fff333", 1592156796},
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
		return block, errors.New("Error height is 1")
	default:
		return block, errors.New("Error default")
	}
}

func (b BlockchainMock) GetBlockByHash(h string) (contracts.Block, error) {
	var block contracts.Block
	switch h {
	case "123abc456":
		return mockBlocks[0], nil
	case "789abc987":
		return mockBlocks[1], nil
	default:
		return block, errors.New("Uknown hash error")
	}
}

func (b BlockchainMock) GetStarsByWalletAddress(addr string) []string {
	var stars []string = make([]string, 0)
	for _, b := range mockBlocks {
		if b.Owner == addr {
			stars = append(stars, b.Body) // In real implementation we will decode it
		}
	}
	return stars
}

func (b BlockchainMock) SubmitStar(star contracts.StarData) (contracts.Block, error) {
	var block contracts.Block
	if star.Message != "" {
		block := contracts.Block{string(star.Data), "1a32", 1, star.Address, mockBlocks[0].Hash, 1592156792}
		return block, nil
	} else {
		return block, errors.New("Empty message error!")
	}
}

func createApi() *httptest.Server {
	var blockchain contracts.BlockchainOperator = BlockchainMock{}
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

func TestGetBlocks(t *testing.T) {
	t.Log("GetBlocks")
	{
		server := createApi()
		defer server.Close()
		t.Log("Server url: ", server.URL)
		t.Log("\tGiven a need to test endpoint /blocks/:address")
		{
			t.Log("\tWhen called with proper address")
			{
				owner := mockBlocks[1].Owner
				response, err := http.Get(server.URL + "/blocks/" + owner)
				if err != nil {
					t.Fatalf("\t\tShould be able to get a block, got err: %v", err)
				}
				t.Log("\t\tShould be able to get blocks for given owner")
				if response.StatusCode != 200 {
					t.Fatalf("\t\tShould get response 200 OK, got: %v", response.StatusCode)
				}
				t.Log("\t\tShould get response 200 OK")
				body, _ := ioutil.ReadAll(response.Body)
				rawBlocks := [...]json.RawMessage{
					json.RawMessage(mockBlocks[1].Body),
					json.RawMessage(mockBlocks[3].Body),
				}
				blocksJson, err := json.Marshal(rawBlocks)
				if err != nil {
					t.Fatal("\t\tShould not get an error when marshalling rawBlocks: ", err)
				}
				if string(blocksJson) != string(body) {
					t.Fatalf("\t\tShould return correct body: %v, got: %v", string(blocksJson), string(body))
				}
				t.Log("\t\tShould return correct body")
			}
			t.Log("\tWhen called with wrong address")
			{
				response, err := http.Get(server.URL + "/blocks/666")
				if err != nil {
					t.Fatal("\t\tShould not get an error for non existing blocks: ", err)
				}
				t.Log("\t\tShould not get an error for non existing blocks")
				if response.StatusCode != 200 {
					t.Fatal("\t\tShould return OK status code, got: ", response.StatusCode)
				}
				t.Log("\t\tShould return OK status code")
				var stars []string
				if err := json.NewDecoder(response.Body).Decode(&stars); err != nil {
					t.Fatalf("\t\tShould decode response body, got err: %v; json: %v", err, response.Body)
				}
				if len(stars) != 0 {
					t.Fatalf("\t\tShould return empty slice, got: %v", stars)
				}
				if stars == nil {
					t.Fatalf("\t\tShould return not nil slice, got: nil")
				}
				t.Log("\t\tShould return empty slice")
			}
		}
	}
}

func TestSubmitStar(t *testing.T) {
	t.Log("SubmitStar")
	{
		server := createApi()
		defer server.Close()
		t.Log("Server url: ", server.URL)
		t.Log("\tGiven a need to test endpoint /submitStar")
		{
			t.Log("\tWhen called with JSON string data")
			{
				addr := "a1b2c3"
				msg := fmt.Sprintf("%s:%d:starRegistry", addr, 1592156792)
				starData := `"This is brand new Star"`
				star := StarDto{
					Address:   addr,
					Message:   msg,
					Data:      json.RawMessage(starData),
					Signature: "doesnotmatter",
				}
				data, _ := json.Marshal(star)
				response, err := http.Post(server.URL+"/submitStar", "application/json", bytes.NewReader(data))
				if err != nil {
					t.Fatalf("\t\tShould be able to post a star, got err: %v", err)
				}
				t.Log("\t\tShould be able to post a star")
				if response.StatusCode != http.StatusCreated {
					body, _ := ioutil.ReadAll(response.Body)
					t.Fatalf("\t\tShould get response 201 Created, got: %v, err: %v", response.StatusCode, string(body))
				}
				t.Log("\t\tShould get response 201 Created")
				var block contracts.Block
				if err := json.NewDecoder(response.Body).Decode(&block); err != nil {
					t.Fatalf("\t\tShould decode response body, got err: %v; json: %v", err, response.Body)
				}
				t.Logf("\t\tShould decode response body %s", response.Body)
				if block.Owner != addr {
					t.Fatalf("\t\tShould return block with owner: %v, got: %v", addr, block.Owner)
				}
				t.Logf("\t\tShould return block with correct owner")
				if block.Height != 1 {
					t.Fatalf("\t\tShould return block with height: %v, got: %v", 1, block.Height)
				}
				t.Logf("\t\tShould return block with correct height")
				if block.Body != starData {
					t.Fatalf("\t\tShould return block with data: %v, got: %v", starData, block.Body)
				}
				t.Logf("\t\tShould return block with correct data")
			}
			t.Log("\tWhen called with JSON object data")
			{
				addr := "a4b5c6"
				msg := fmt.Sprintf("%s:%d:starRegistry", addr, 1592156792)
				starData := `{"text":"This is brand new Star"}`
				star := StarDto{
					Address:   addr,
					Message:   msg,
					Data:      json.RawMessage(starData),
					Signature: "doesnotmatter",
				}
				data, _ := json.Marshal(star)
				response, err := http.Post(server.URL+"/submitStar", "application/json", bytes.NewReader(data))
				if err != nil {
					t.Fatalf("\t\tShould be able to post a star, got err: %v", err)
				}
				t.Log("\t\tShould be able to post a star")
				if response.StatusCode != http.StatusCreated {
					body, _ := ioutil.ReadAll(response.Body)
					t.Fatalf("\t\tShould get response 201 Created, got: %v, err: %v", response.StatusCode, string(body))
				}
				t.Log("\t\tShould get response 201 Created")
				var block contracts.Block
				if err := json.NewDecoder(response.Body).Decode(&block); err != nil {
					t.Fatalf("\t\tShould decode response body, got err: %v; json: %v", err, response.Body)
				}
				t.Logf("\t\tShould decode response body %s", response.Body)
				if block.Owner != addr {
					t.Fatalf("\t\tShould return block with owner: %v, got: %v", addr, block.Owner)
				}
				t.Logf("\t\tShould return block with correct owner")
				if block.Height != 1 {
					t.Fatalf("\t\tShould return block with height: %v, got: %v", 1, block.Height)
				}
				t.Logf("\t\tShould return block with correct height")
				if block.Body != starData {
					t.Fatalf("\t\tShould return block with data: %v, got: %v", starData, block.Body)
				}
				t.Logf("\t\tShould return block with correct data")
			}
			t.Log("\tWhen called with wrong data")
			{
				var star StarDto
				data, _ := json.Marshal(star)
				response, err := http.Post(server.URL+"/submitStar", "application/json", bytes.NewReader(data))
				if err != nil {
					t.Fatal("\t\tShould not get an error for malformed star data, got: ", err)
				}
				t.Log("\t\tShould not get an error for malformed star data")
				if response.StatusCode != http.StatusInternalServerError {
					t.Fatal("\t\tShould return InternalServerError status code and error, got: ", response.StatusCode)
				}
				body, _ := ioutil.ReadAll(response.Body)
				t.Log("\t\tShould return InternalServerError status code and error:", string(body))
			}
		}
	}
}

func TestValidate(t *testing.T) {
	t.Log("Validate")
	{
		server := createApi()
		defer server.Close()
		t.Log("Server url: ", server.URL)
		t.Log("\tGiven a need to test endpoint /validate")
		{
			t.Log("\tWhen called on fresh blockchain")
			{
				response, err := http.Get(server.URL + "/validate")
				if err != nil {
					t.Fatalf("\t\tShould be able to get validation result, got err: %v", err)
				}
				t.Log("\t\tShould be able to get validation result")
				if response.StatusCode != 200 {
					t.Fatalf("\t\tShould get response 200 OK, got: %v", response.StatusCode)
				}
				t.Log("\t\tShould get response 200 OK")
				// TODO decode response to Dto
			}
			t.Log("\tWhen called on blockchain with 3 blocks")
			{
				// TODO add blocks
				response, err := http.Get(server.URL + "/validate")
				if err != nil {
					t.Fatal("\t\tShould not get an error when validating chain, got: ", err)
				}
				t.Log("\t\tShould not get an error when validating chain")
				if response.StatusCode != 200 {
					t.Fatal("\t\tShould return OK status code, got: ", response.StatusCode)
				}
				t.Log("\t\tShould return OK status code")
				// TODO decode response to valid DTO
				//if err := json.NewDecoder(response.Body).Decode(&validationResult); err != nil {
				//	t.Fatalf("\t\tShould decode response body, got err: %v; json: %v", err, response.Body)
				//}
			}
			t.Log("\tWhen called on blockchain with tempered blocks")
			{
				// TODO modify blocks
				response, err := http.Get(server.URL + "/validate")
				if err != nil {
					t.Fatal("\t\tShould not get an error when validating chain, got: ", err)
				}
				t.Log("\t\tShould not get an error when validating chain")
				if response.StatusCode != 200 {
					t.Fatal("\t\tShould return OK status code, got: ", response.StatusCode)
				}
				t.Log("\t\tShould return OK status code")
				// TODO decode response to valid DTO
				//if err := json.NewDecoder(response.Body).Decode(&validationResult); err != nil {
				//	t.Fatalf("\t\tShould decode response body, got err: %v; json: %v", err, response.Body)
				//}
			}
		}
	}
}
