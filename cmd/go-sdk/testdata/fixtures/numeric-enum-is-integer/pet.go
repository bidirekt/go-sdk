package petstore

type Priority int

const (
	PriorityLow  Priority = 0
	PriorityHigh Priority = 10
)

type Pet struct {
	Priority Priority `json:"priority" bidirekt:"field"`
}
