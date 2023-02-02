package structure

import "time"

type Failure struct {
	Id        int64     `xorm:"pk autoincr" json:"id"`
	Issue     string    `json:"-"`
	ErrorCode int64     `json:"-"`
	Service   int64     `json:"service"`
	PingTime  uint64    `json:"ping_time"`
	Reason    string    `json:"-"`
	CreatedAt time.Time `xorm:"created index" json:"-"`
}
