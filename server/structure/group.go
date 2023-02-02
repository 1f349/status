package structure

import "time"

type Group struct {
	Id        int64     `xorm:"pk autoincr" json:"id"`
	Name      string    `json:"name"`
	Public    bool      `json:"public"`
	Order     int64     `json:"order"`
	CreatedAt time.Time `xorm:"created" json:"-"`
	UpdatedAt time.Time `xorm:"updated" json:"-"`
}
