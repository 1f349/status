package ping

import (
	"code.mrmelon54.com/melon/status/server/structure"
	"log"
	"time"
	"xorm.io/xorm"
)

type Maintenance struct {
	engine *xorm.Engine
}

func NewMaintenance(engine *xorm.Engine) *Maintenance {
	return &Maintenance{engine: engine}
}

func (m *Maintenance) Run() {
	go m.internalCall()
}

func (m *Maintenance) internalCall() {
	d := 92 * 24 * time.Hour
	ticker := time.Hour

	log.Println("[Maintenance::internalCall()] Clearing up database before starting")
	m.deleteAllExtraItems(d)

	for {
		select {
		case <-time.After(ticker):
			m.deleteAllExtraItems(d)
		}
	}
}

func (m *Maintenance) deleteAllExtraItems(d time.Duration) {
	m.deleteAllSince(&structure.Hit{}, d)
	m.deleteAllSince(&structure.Failure{}, d)
}

func (m *Maintenance) deleteAllSince(table any, d time.Duration) {
	_, err := m.engine.Where("created_at < from_unixtime(?)", time.Now().Add(-d).Unix()).Delete(table)
	if err != nil {
		log.Printf("[Maintenance::deleteAllSince(\"%#v\")] %s", table, err)
		return
	}
}
