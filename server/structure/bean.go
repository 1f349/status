package structure

import "time"

type BeanState int

const (
	BeanStateUnknown BeanState = iota
	BeanStateHit
	BeanStateFailure
)

type Bean struct {
	State     BeanState `json:"state,omitempty"`
	CreatedAt time.Time `json:"-"`
	Time      int64     `json:"time,omitempty"`
}
