package petstore

import (
	"net/url"
)

type Owner struct {
	Name string `json:"name"`
}

type Pet struct {
	Owner   Owner   `json:"owner" bidirekt:"field"`
	Friends []Owner `json:"friends" bidirekt:"field"`
	Website url.URL `json:"website" bidirekt:"field"`
}
