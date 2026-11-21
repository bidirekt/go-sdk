package petstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bidirekt/go-sdk"
)

var ErrNotFound = errors.New("pet not found")

type Pet struct {
	ID     int64  `json:"id" bidirekt:"field"`
	Name   string `json:"name" bidirekt:"field"`
	Status string `json:"status" bidirekt:"field"`
}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

var _ = bidirekt.Consumes("petstore_api", http.MethodGet, "/pets",
	bidirekt.Response(http.StatusOK, []Pet{}),
)

func (this *Client) ListPets(ctx context.Context) ([]Pet, error) {
	var pets []Pet
	_, err := this.get(ctx, "/pets", &pets)
	return pets, err
}

var _ = bidirekt.Consumes("petstore_api", http.MethodGet, "/pets/*",
	bidirekt.Response(http.StatusOK, Pet{}),
	bidirekt.Response(http.StatusNotFound, nil),
)

func (this *Client) GetPet(ctx context.Context, id int64) (Pet, error) {
	var pet Pet
	status, err := this.get(ctx, fmt.Sprintf("/pets/%d", id), &pet)
	if status == http.StatusNotFound {
		return Pet{}, ErrNotFound
	}
	return pet, err
}

func (this *Client) get(ctx context.Context, path string, body any) (int, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, this.BaseURL+path, nil)
	if err != nil {
		return 0, err
	}
	response, err := this.HTTPClient.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return response.StatusCode, nil
	}
	return response.StatusCode, json.NewDecoder(response.Body).Decode(body)
}
