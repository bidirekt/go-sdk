package petstore

import (
	"example.com/shared"
	"github.com/bidirekt/go-sdk"
)

type Mood string

const (
	MoodCalm    Mood = "calm"
	MoodPlayful Mood = "playful"
)

var _ = bidirekt.Enum(MoodCalm, MoodPlayful)

var _ = bidirekt.Provides("GET", "/pets/*",
	bidirekt.Response(200, shared.Pet{}),
)
