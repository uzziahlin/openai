// Copyright 2023 Ken Lin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openai

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// newMockClient returns a mock client for testing purposes.
func newMockClient(baseUrl string, opts ...Option) *Client {
	app := App{
		ApiKey: os.Getenv("OPENAI_API_KEY"),
		ApiUrl: baseUrl,
	}

	client, err := New(app, opts...)
	if err != nil {
		panic(fmt.Sprintf("Cannot create client: %v", err))
	}

	return client
}

// newMockServer returns a mock server for testing purposes.
func newMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func newMockHandler(t *testing.T, method string, filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, method, r.Method)

		mockData := loadTestdata(filename)

		// 模拟网络延迟
		time.Sleep(3 * time.Second)
		_, _ = w.Write(mockData)
	}
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

// loadMockData loads mock data from a file.
func loadMockData(filename string, out any) {
	testdata := loadTestdata(filename)
	if err := json.Unmarshal(testdata, &out); err != nil {
		panic(fmt.Sprintf("decode mock data error: %s", err))
	}
}
