package openai

import "context"

type FineTuneService interface {
	Create(ctx context.Context, req *FineTuneCreateRequest) (*FineTune, error)
	List(ctx context.Context) (*FineTuneListResponse, error)
	Retrieve(ctx context.Context, id string) (*FineTune, error)
	Cancel(ctx context.Context, id string) (*FineTune, error)
	ListEvents(ctx context.Context, id string, stream ...bool) (*EventListResponse, error)
	DeleteModel(ctx context.Context, model string) (*ModelDeleteResponse, error)
}

type FineTuneCreateRequest struct {
	TrainingFile                  string    `json:"training_file"`
	ValidationFile                string    `json:"validation_file,omitempty"`
	Model                         string    `json:"model,omitempty"`
	NEpochs                       int       `json:"n_epochs,omitempty"`
	BatchSize                     int       `json:"batch_size,omitempty"`
	LearningRateMultiplier        float64   `json:"learning_rate_multiplier,omitempty"`
	PromptLossWeight              float64   `json:"prompt_loss_weight,omitempty"`
	ComputeClassificationMetrics  bool      `json:"compute_classification_metrics,omitempty"`
	ClassificationNClasses        int64     `json:"classification_n_classes,omitempty"`
	ClassificationPositiveClasses string    `json:"classification_positive_classes,omitempty"`
	ClassificationBetas           []float64 `json:"classification_betas,omitempty"`
	Suffix                        string    `json:"suffix,omitempty"`
}

type FineTune struct {
	Id              string      `json:"id"`
	Object          string      `json:"object"`
	Model           string      `json:"model"`
	CreateAt        int64       `json:"create_at"`
	Events          []*Event    `json:"events"`
	FineTunedModel  string      `json:"fine_tuned_model"`
	Hyperparams     Hyperparams `json:"hyperparams"`
	OrganizationId  string      `json:"organization_id"`
	ResultFiles     []*File     `json:"result_files"`
	Status          string      `json:"status"`
	ValidationFiles []*File     `json:"validation_files"`
	TrainingFiles   []*File     `json:"training_files"`
	UpdatedAt       int64       `json:"updated_at"`
}

type Hyperparams struct {
	BatchSize              int64   `json:"batch_size"`
	LearningRateMultiplier float64 `json:"learning_rate_multiplier"`
	NEpochs                int64   `json:"n_epochs"`
	PromptLossWeight       float64 `json:"prompt_loss_weight"`
}

type Event struct {
	Object   string `json:"object"`
	CreateAt int64  `json:"create_at"`
	Level    string `json:"level"`
	Message  string `json:"message"`
}

type FineTuneListResponse struct {
	Object string      `json:"object"`
	Data   []*FineTune `json:"data"`
}

type EventListResponse struct {
	Object string   `json:"object"`
	Data   []*Event `json:"data"`
}

type ModelDeleteResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

type FineTuneServiceOp struct {
	client *Client
}

func (f FineTuneServiceOp) Create(ctx context.Context, req *FineTuneCreateRequest) (*FineTune, error) {
	//TODO implement me
	panic("implement me")
}

func (f FineTuneServiceOp) List(ctx context.Context) (*FineTuneListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (f FineTuneServiceOp) Retrieve(ctx context.Context, id string) (*FineTune, error) {
	//TODO implement me
	panic("implement me")
}

func (f FineTuneServiceOp) Cancel(ctx context.Context, id string) (*FineTune, error) {
	//TODO implement me
	panic("implement me")
}

func (f FineTuneServiceOp) ListEvents(ctx context.Context, id string, stream ...bool) (*EventListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (f FineTuneServiceOp) DeleteModel(ctx context.Context, model string) (*ModelDeleteResponse, error) {
	//TODO implement me
	panic("implement me")
}
