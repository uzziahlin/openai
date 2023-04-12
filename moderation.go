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

import "context"

const (
	ModerationCreatePath = "/moderations"
)

type ModerationService interface {
	Create(ctx context.Context, req *ModerationCreateRequest) (*ModerationCreateResponse, error)
}

type ModerationCreateRequest struct {
	Input string `json:"input"`
	Model string `json:"model,omitempty"`
}

type ModerationCreateResponse struct {
	Id      string        `json:"id"`
	Model   string        `json:"model"`
	Results []*Moderation `json:"results"`
}

type Moderation struct {
	Categories     ModerationCategory      `json:"categories"`
	CategoryScores ModerationCategoryScore `json:"category_scores"`
	Flagged        bool                    `json:"flagged"`
}

type ModerationCategory struct {
	Hate            bool `json:"hate"`
	HateThreatening bool `json:"hate/threatening"`
	SelfHarm        bool `json:"self-harm"`
	Sexual          bool `json:"sexual"`
	SexualMinors    bool `json:"sexual/minors"`
	Violence        bool `json:"violence"`
	ViolenceGraphic bool `json:"violence/graphic"`
}

type ModerationCategoryScore struct {
	Hate            float64 `json:"hate"`
	HateThreatening float64 `json:"hate/threatening"`
	SelfHarm        float64 `json:"self-harm"`
	Sexual          float64 `json:"sexual"`
	SexualMinors    float64 `json:"sexual/minors"`
	Violence        float64 `json:"violence"`
	ViolenceGraphic float64 `json:"violence/graphic"`
}

type ModerationServiceOp struct {
	client *Client
}

func (m ModerationServiceOp) Create(ctx context.Context, req *ModerationCreateRequest) (*ModerationCreateResponse, error) {
	var resp ModerationCreateResponse
	err := m.client.Post(ctx, ModerationCreatePath, req, &resp)
	return &resp, err
}
