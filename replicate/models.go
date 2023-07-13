package replicate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ModelsService service

type Model struct {
	URL           string       `json:"url"`
	Owner         string       `json:"owner"`
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	Visibility    string       `json:"visibility"`
	GithubURL     string       `json:"github_url"`
	PaperURL      string       `json:"paper_url"`
	LicenseURL    string       `json:"license_url"`
	RunCount      int          `json:"run_count"`
	CoverImageURL string       `json:"cover_image_url"`
	LatestVersion ModelVersion `json:"latest_version"`
}

type ModelVersion struct {
	ID         string    `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	CogVersion string    `json:"cog_version"`
}

func (s *ModelsService) Get(owner string, name string) (*Model, error) {
	url := fmt.Sprintf("models/%s/%s", owner, name)
	req, err := s.client.baseRequest("GET", url)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var model Model
	err = json.NewDecoder(resp.Body).Decode(&model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}
