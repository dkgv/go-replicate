package replicate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrPredictionFailed   = errors.New("prediction failed")
	ErrPredictionCanceled = errors.New("prediction canceled")
)

type PredictionsService service

type PredictionEvent string

const (
	EventPredictionStarted   PredictionEvent = "start"
	EventPredictionOutputted PredictionEvent = "output"
	EventPredictionLogged    PredictionEvent = "logs"
	EventPredictionCompleted PredictionEvent = "completed"
)

type Prediction struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Urls    struct {
		Get    string `json:"get"`
		Cancel string `json:"cancel"`
	} `json:"urls"`
	StartedAt   time.Time `json:"started_at"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
	Source      string    `json:"source"`
	Status      string    `json:"status"`
	Input       any       `json:"input"`
	Output      any       `json:"output"`
	Error       any       `json:"error"`
	Logs        string    `json:"logs"`
	Metrics     struct {
		PredictTime float64 `json:"predict_time"`
	} `json:"metrics"`
	Webhook             string            `json:"webhook,omitempty"`
	WebhookEventsFilter []PredictionEvent `json:"webhook_events_filter,omitempty"`
}

type Webhook struct {
	CallbackURL string            `json:"webhook"`
	Events      []PredictionEvent `json:"webhook_events_filter"`
}

func (s *PredictionsService) Create(ctx context.Context, modelID string, input any) (*Prediction, error) {
	return s.CreateWithWebhook(ctx, modelID, input, Webhook{})
}

func (s *PredictionsService) CreateWithWebhook(ctx context.Context, modelID string, input any, webhook Webhook) (*Prediction, error) {
	req, err := s.client.baseRequest(ctx, "POST", "predictions")
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(
		struct {
			Version string            `json:"version"`
			Input   any               `json:"input"`
			Webhook string            `json:"webhook,omitempty"`
			Events  []PredictionEvent `json:"webhook_events_filter,omitempty"`
		}{
			Version: modelID,
			Input:   input,
			Webhook: webhook.CallbackURL,
			Events:  webhook.Events,
		},
	)
	if err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewReader(body))

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var prediction Prediction
	err = json.NewDecoder(resp.Body).Decode(&prediction)
	if err != nil {
		return nil, err
	}

	return &prediction, nil
}

func (s *PredictionsService) Get(ctx context.Context, id string) (*Prediction, error) {
	req, err := s.client.baseRequest(ctx, "GET", fmt.Sprintf("predictions/%s", id))
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var prediction Prediction
	err = json.NewDecoder(resp.Body).Decode(&prediction)
	if err != nil {
		return nil, err
	}

	return &prediction, nil
}

func (s *PredictionsService) Await(ctx context.Context, id string, destination any) error {
	var (
		prediction *Prediction
		err        error
	)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			prediction, err = s.Get(ctx, id)
			if err != nil {
				return err
			}

			switch prediction.Status {
			case "succeeded":
				b, err := json.Marshal(prediction.Output)
				if err != nil {
					return err
				}
				return json.Unmarshal(b, &destination)
			case "failed":
				return ErrPredictionFailed
			case "canceled":
				return ErrPredictionCanceled
			case "starting", "processing":
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (s *PredictionsService) Cancel(ctx context.Context, id string) error {
	req, err := s.client.baseRequest(ctx, "POST", fmt.Sprintf("predictions/%s/cancel", id))
	if err != nil {
		return err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
