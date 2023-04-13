# OpenAI Go SDK
This is an unofficial Go SDK for interacting with the OpenAI API, which allows you to leverage the power of OpenAI's models like GPT-4 within your Go applications. This SDK aims to simplify the process of integrating OpenAI's API with your projects by providing an easy-to-use interface.

**Note**: This SDK is not officially supported by OpenAI and comes with no warranty.

## Table of Contents
- [Features](#features)
- [Supported Services](#supported-services)
- [Installation](#installation)
- [Usage](#usage)
- License

## Features
- Easy integration with the OpenAI API
- Support for various OpenAI models
- Simple and intuitive function calls
- Extensible for future API updates

## Supported Services
The following services are currently supported by the SDK:

- [Models](https://platform.openai.com/docs/api-reference/models)
- [Completions](https://platform.openai.com/docs/api-reference/completions)
- [Chat Completions](https://platform.openai.com/docs/api-reference/chat) 
- [Edits](https://platform.openai.com/docs/api-reference/edits) 
- [Images](https://platform.openai.com/docs/api-reference/images)
- [Embeddings](https://platform.openai.com/docs/api-reference/embeddings) 
- [Audio](https://platform.openai.com/docs/api-reference/audio)
- [Files](https://platform.openai.com/docs/api-reference/files)
- [Fine-tunes](https://platform.openai.com/docs/api-reference/fine-tunes)
- [Moderations](https://platform.openai.com/docs/api-reference/moderations)

## Installation
You can install the OpenAI Go SDK using the following command:
```bash
go get -u github.com/uzziahlin/openai
```

## Usage
The following code snippet shows how to use the SDK to create a chat completion using the `gpt-3.5-turbo` model:
```go
package main

import (
	"context"
	"fmt"
	"github.com/uzziahlin/openai"
	"os"
)

func main() {
	// get api key, suggest to use environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")

	// create a new client
	client, err := openai.New(openai.App{
		ApiUrl: "https://api.openai.com",
		ApiKey: apiKey,
	})

	// if you want to use proxy, you can use the following code
	//client, err := openai.New(openai.App{
	//	ApiUrl: "https://api.openai.com",
	//	ApiKey: apiKey,
	//}, openai.WithProxy(&openai.Proxy{
	//	Url:      "your proxy url"
	//	Username: "your proxy username, if not need, just leave it",
	//	Password: "your proxy password, if not need, just leave it",
	//}))

	if err != nil {
		// handle error
	}

	// create a new chat session, not stream
	resp, err := client.Chat.Create(context.TODO(), &openai.ChatCreateRequest{
		Model: "gpt-3.5-turbo",
		Messages: []*openai.Message{
			// can add more messages here
			{
				Role:    "user",
				Content: "Hello, How are you?",
			},
		},
	})

	// To be compatible with the streaming API, the returned value is a channel
	// and since it is not a stream, there is only one element, which can be taken out directly
	res, ok := <-resp

	if !ok {
		// if the channel is closed, it means that the request has been completed
	}

	// handle response
	// for example, print the first choice
	fmt.Println(res.Choices[0].Message.Content)
}
```
if you want to use the streaming API, you can use the following code:
```go
// The flow of creating a client and a non-streaming flow is the same
// only when creating a session, the stream parameter is set to true
resp, err := client.Chat.Create(context.TODO(), &openai.ChatCreateRequest{
    Model: "gpt-3.5-turbo",
    Messages: []*openai.Message{
        // can add more messages here
        {
            Role:    "user",
            Content: "Hello, How are you?",
        },
    },
    Stream: true,
})

if err != nil {
    // handle error
}

for {
    res, ok := <-resp
    if !ok {
        // channel is closed
        break
    }
    
    // handle response
    // for example, print the response
    // the content of the response is in the res.Choices[0].Delta.Content
    fmt.Println(res.Choices[0].Delta.Content)
}
```
other services are similar to the above usage, so I won't repeat it here.

## License
This project is licensed under the Apache License 2.0. Please see the LICENSE file for more details.