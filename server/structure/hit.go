package structure

import "time"

type Hit struct {
	Id        int64         `xorm:"pk autoincr" json:"id"`
	Service   int64         `json:"service"`
	Latency   time.Duration `json:"latency"`
	PingTime  time.Duration `json:"ping_time"`
	CreatedAt time.Time     `xorm:"created index" json:"-"`
}
