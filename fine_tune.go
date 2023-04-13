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
	"context"
	"encoding/json"
	"fmt"
)

const (
	FineTuneCreatePath   = "/fine-tunes"
	FineTuneListPath     = "/fine-tunes"
	FineTuneRetrievePath = "/fine-tunes/%s"
	FineTuneCancelPath   = "/fine-tunes/%s/cancel"
	EventsListPath       = "/fine-tunes/%s/events"
	ModelDeletePath      = "/models/%s"
)

type FineTuneService interface {
	Create(ctx context.Context, req *FineTuneCreateRequest) (*FineTune, error)
	List(ctx context.Context) (*FineTuneListResponse, error)
	Retrieve(ctx context.Context, id string) (*FineTune, error)
	Cancel(ctx context.Context, id string) (*FineTune, error)
	ListEvents(ctx context.Context, id string, stream ...bool) (chan *EventListResponse, error)
	DeleteModel(ctx context.Context, model string) (*ModelDeleteResponse, error)
}

type FineTuneCreateRequest struct {
	// TrainingFile is the ID to the training file
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
	Id              string           `json:"id"`
	Object          string           `json:"object"`
	Model           string           `json:"model"`
	CreateAt        int64            `json:"create_at"`
	Events          []*FineTuneEvent `json:"events"`
	FineTunedModel  string           `json:"fine_tuned_model"`
	Hyperparams     Hyperparams      `json:"hyperparams"`
	OrganizationId  string           `json:"organization_id"`
	ResultFiles     []*File          `json:"result_files"`
	Status          string           `json:"status"`
	ValidationFiles []*File          `json:"validation_files"`
	TrainingFiles   []*File          `json:"training_files"`
	UpdatedAt       int64            `json:"updated_at"`
}

type Hyperparams struct {
	BatchSize              int64   `json:"batch_size"`
	LearningRateMultiplier float64 `json:"learning_rate_multiplier"`
	NEpochs                int64   `json:"n_epochs"`
	PromptLossWeight       float64 `json:"prompt_loss_weight"`
}

type FineTuneEvent struct {
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
	Object string           `json:"object"`
	Data   []*FineTuneEvent `json:"data"`
}

type ModelDeleteResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

type FineTuneServiceOp struct {
	client *Client
}

// Create Creates a job that fine-tunes a specified model from a given dataset.
//Response includes details of the enqueued job including job status and the name of the fine-tuned models once complete.
func (f FineTuneServiceOp) Create(ctx context.Context, req *FineTuneCreateRequest) (*FineTune, error) {
	var resp FineTune
	err := f.client.Post(ctx, FineTuneCreatePath, req, &resp)
	return &resp, err
}

// List Returns a list of all fine-tuning jobs.
func (f FineTuneServiceOp) List(ctx context.Context) (*FineTuneListResponse, error) {
	var resp FineTuneListResponse
	err := f.client.Get(ctx, FineTuneListPath, nil, &resp)
	return &resp, err
}

// Retrieve Returns a fine-tuning job by ID.
func (f FineTuneServiceOp) Retrieve(ctx context.Context, id string) (*FineTune, error) {
	var resp FineTune
	err := f.client.Get(ctx, fmt.Sprintf(FineTuneRetrievePath, id), nil, &resp)
	return &resp, err
}

// Cancel Cancels a fine-tuning job.
func (f FineTuneServiceOp) Cancel(ctx context.Context, id string) (*FineTune, error) {
	var resp FineTune
	err := f.client.Post(ctx, fmt.Sprintf(FineTuneCancelPath, id), nil, &resp)
	return &resp, err
}

// ListEvents Returns a list of events for a fine-tuning job.
// If stream=true, the response will be a stream of events as they are generated.
// by default, the response will be a list of all events generated so far.
func (f FineTuneServiceOp) ListEvents(ctx context.Context, id string, stream ...bool) (chan *EventListResponse, error) {
	type Stream struct {
		Stream bool `url:"stream"`
	}

	var s Stream
	if len(stream) > 0 {
		s.Stream = stream[0]
	}

	if !s.Stream {
		var resp EventListResponse
		err := f.client.Get(ctx, fmt.Sprintf(EventsListPath, id), s, &resp)
		if err != nil {
			return nil, err
		}
		ch := make(chan *EventListResponse, 1)
		ch <- &resp
		close(ch)
		return ch, nil
	}

	es, err := f.client.GetByStream(ctx, fmt.Sprintf(EventsListPath, id), s)

	if err != nil {
		return nil, err
	}

	ch := make(chan *EventListResponse)

	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			case e, ok := <-es:
				if !ok {
					return
				}
				var resp EventListResponse
				err := json.Unmarshal([]byte(e.Data), &resp)
				if err != nil {
					f.client.logger.Error(err, "failed to unmarshal chat response")
					continue
				}
				select {
				case <-ctx.Done():
					return
				case ch <- &resp:
				}
			}
		}
	}()

	return ch, nil
}

// DeleteModel Deletes a fine-tuned model.
func (f FineTuneServiceOp) DeleteModel(ctx context.Context, model string) (*ModelDeleteResponse, error) {
	var resp ModelDeleteResponse
	err := f.client.Delete(ctx, fmt.Sprintf(ModelDeletePath, model), nil, &resp)
	return &resp, err
}
