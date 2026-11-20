package petstore

import (
	"time"
)

type Owner struct {
	Name string `json:"name" bidirekt:"field"`
}

type Pet struct {
	BornAt    time.Time      `json:"bornAt" bidirekt:"field type=string"`
	VisitedAt *time.Time     `json:"visitedAt" bidirekt:"field type=string"`
	Photo     []byte         `json:"photo" bidirekt:"field type=string"`
	Extra     any            `json:"extra" bidirekt:"field type=Owner"`
	ByNumber  map[int]string `json:"byNumber" bidirekt:"field type=map values=string"`
	Visits    []time.Time    `json:"visits" bidirekt:"field type=array items=string"`
}
