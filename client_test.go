package openai

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
)

// newMockClient returns a mock client for testing purposes.
func newMockClient(baseUrl string) *Client {
	app := App{
		ApiKey: os.Getenv("OPENAI_API_KEY"),
		ApiUrl: baseUrl,
	}

	client, err := New(app)
	if err != nil {
		panic(fmt.Sprintf("Cannot create client: %v", err))
	}

	return client
}

// newMockServer returns a mock server for testing purposes.
func newMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// loadTestdata loads testdata from a file.
func loadTestdata(filename string) []byte {
	// 读取文件
	file, err := os.Open("testdata/" + filename)
	if err != nil {
		panic(fmt.Sprintf("Cannot load testdata: %v", filename))
	}
	defer file.Close()

	// 读取文件中的数据
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	return data
}

func loadMockData(filename string, out any) {
	testdata := loadTestdata(filename)
	if err := json.Unmarshal(testdata, &out); err != nil {
		panic(fmt.Sprintf("decode mock data error: %s", err))
	}
}
