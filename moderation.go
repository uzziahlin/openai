package openai

import "context"

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
	//TODO implement me
	panic("implement me")
}
