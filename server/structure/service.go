package structure

import "time"

type Service struct {
	Id                 int64     `xorm:"pk autoincr" json:"id"`
	Name               string    `json:"name"`
	Domain             string    `json:"-"`
	Expected           string    `json:"-"`
	ExpectedStatus     int64     `json:"-"`
	CheckInternal      int64     `json:"-"`
	Method             string    `json:"-"`
	PostData           []byte    `xorm:"blob" json:"-"`
	Order              int64     `json:"order"`
	Public             bool      `json:"public"`
	GroupId            int64     `json:"group_id"`
	Permalink          string    `xorm:"index" json:"permalink"`
	CreatedAt          time.Time `xorm:"created" json:"-"`
	UpdatedAt          time.Time `xorm:"updated" json:"-"`
	NotifyAfter        int64     `json:"-"`
	AllowNotifications bool      `json:"-"`
}
