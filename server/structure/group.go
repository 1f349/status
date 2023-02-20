package structure

import "time"

type Group struct {
	Id        int64     `xorm:"pk autoincr" json:"id"`
	Name      string    `json:"name"`
	Public    bool      `json:"public,omitempty"`
	Order     int64     `json:"order,omitempty"`
	CreatedAt time.Time `xorm:"created" json:"-"`
	UpdatedAt time.Time `xorm:"updated" json:"-"`
	Services  []Service `xorm:"-" json:"services"`
}
