package shared

import "time"

type Timestamps struct {
	CreatedAt time.Time  `json:"createdAt" bidirekt:"field type=string description='RFC 3339 date and time'"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" bidirekt:"field type=string"`
}

func NewTimestamps() Timestamps {
	return Timestamps{CreatedAt: time.Now().UTC()}
}

func (this *Timestamps) Touch() {
	now := time.Now().UTC()
	this.UpdatedAt = &now
}
