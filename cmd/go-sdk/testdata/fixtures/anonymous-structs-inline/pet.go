package petstore

type Pet struct {
	Address struct {
		City   string `json:"city" bidirekt:"field"`
		Street string `json:"street"`
		Geo    struct {
			Lat float64 `json:"lat" bidirekt:"field"`
			Lng float64 `json:"lng" bidirekt:"field"`
		} `json:"geo" bidirekt:"field"`
	} `json:"address" bidirekt:"field"`
	Visits []struct {
		At string `json:"at" bidirekt:"field"`
	} `json:"visits" bidirekt:"field"`
	Notes map[string]struct {
		Text string `json:"text" bidirekt:"field"`
	} `json:"notes" bidirekt:"field"`
}
