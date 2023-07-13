package replicate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PredictionsService service

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
}

func (s *PredictionsService) Create(model Model, input any) (*Prediction, error) {
	req, err := s.client.baseRequest("POST", "predictions")
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(
		map[string]interface{}{
			"version": model.LatestVersion.ID,
			"input":   input,
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

func (s *PredictionsService) Get(id string) (*Prediction, error) {
	req, err := s.client.baseRequest("GET", fmt.Sprintf("predictions/%s", id))
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

func (s *PredictionsService) Await(id string, destination any) error {
	var (
		prediction *Prediction
		err        error
	)

	ok := make(chan bool)
	go func() {
		for {
			prediction, err = s.Get(id)
			if err != nil {
				ok <- false
				return
			}

			switch prediction.Status {
			case "succeeded":
				ok <- true
				return

			case "failed", "canceled":
				ok <- false
				return

			case "starting", "processing":
				time.Sleep(1 * time.Second)
			}
		}
	}()

	success := <-ok
	if !success {
		return fmt.Errorf("Prediction failed or was canceled")
	}

	b, err := json.Marshal(prediction.Output)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &destination)
}

func (s *PredictionsService) Cancel(id string) error {
	req, err := s.client.baseRequest("POST", fmt.Sprintf("predictions/%s/cancel", id))
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
