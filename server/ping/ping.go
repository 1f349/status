package ping

import (
	"bytes"
	"code.mrmelon54.com/melon/status/server/notifier"
	"code.mrmelon54.com/melon/status/server/structure"
	"fmt"
	"github.com/tcnksm/go-httpstat"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
	"xorm.io/xorm"
)

type Ping struct {
	engine         *xorm.Engine
	wait           *sync.WaitGroup
	mutex          *sync.RWMutex
	cachedServices map[int64]structure.Service
	cachedOnline   map[int64]bool
	cachedDowntime map[int64]time.Time
	serviceTime    map[int64]time.Time
	notifiers      []notifier.Notifier
}

func New(engine *xorm.Engine, notifiers []notifier.Notifier) *Ping {
	return &Ping{
		engine:         engine,
		wait:           &sync.WaitGroup{},
		mutex:          &sync.RWMutex{},
		cachedServices: make(map[int64]structure.Service),
		cachedOnline:   make(map[int64]bool),
		cachedDowntime: make(map[int64]time.Time),
		serviceTime:    make(map[int64]time.Time),
		notifiers:      notifiers,
	}
}

func (p *Ping) sendNotify(a string) {
	for _, i := range p.notifiers {
		i.SendMessage(a)
	}
}

func (p *Ping) Run() {
	p.wait.Add(1)
	go p.internalCall()
}

func (p *Ping) Wait() {
	p.wait.Wait()
}

func (p *Ping) internalCall() {
	defer p.wait.Done()

	log.Println("[Ping] Starting internal checking service")
	for {
		time.Sleep(50 * time.Millisecond)
		for _, service := range p.cachedServices {
			p.checkService(service)
		}
	}
}

func (p *Ping) Reload() {
	go p.dirtyReload()
}

func (p *Ping) dirtyReload() {
	p.mutex.Lock()
	var services []structure.Service
	err := p.engine.Find(&services)
	if err != nil {
		log.Printf("[Ping::dirtyReload] Database error: %s\n", err)
		return
	}
	for _, i := range services {
		p.cachedServices[i.Id] = i
	}
	p.mutex.Unlock()
}

func (p *Ping) checkService(service structure.Service) {
	if sTime, ok := p.serviceTime[service.Id]; ok {
		if time.Now().After(sTime) {
			p.internalCheckService(service)
		}
	} else {
		p.internalCheckService(service)
	}
}

func (p *Ping) internalCheckService(service structure.Service) {
	p.serviceTime[service.Id] = time.Now().Add(time.Second * time.Duration(service.CheckInternal))
	log.Println("[Ping] Checking service:", service.Name)
	resp, timings, err := p.makeRequest(service)
	if err != nil {
		p.makeError(service, err, timings)
		return
	}
	if int64(resp.StatusCode) != service.ExpectedStatus {
		p.makeFailure(service, resp, timings)
	} else {
		p.makeHit(service, resp, timings)
	}
}

func (p *Ping) makeRequest(service structure.Service) (*http.Response, *httpstat.Result, error) {
	buf := new(bytes.Buffer)
	var result httpstat.Result
	req, err := http.NewRequest(service.Method, service.Domain, buf)
	if err != nil {
		return nil, nil, err
	}
	req2 := req.WithContext(httpstat.WithHTTPStat(req.Context(), &result))
	resp, err := p.makeClient().Do(req2)
	if err != nil {
		return nil, nil, err
	}
	return resp, &result, nil
}

func (p *Ping) makeClient() *http.Client {
	return &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   time.Second * 15,
	}
}

func (p *Ping) makeHit(service structure.Service, resp *http.Response, timings *httpstat.Result) {
	_, _ = io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	t := time.Now()
	_, err := p.engine.Insert(structure.Hit{
		Id:       0,
		Service:  service.Id,
		Latency:  timings.Total(t),
		PingTime: timings.DNSLookup,
	})
	if err != nil {
		log.Printf("[Ping::makeHit(Service %d)] Database error: %s\n", service.Id, err)
	}

	// update internal cache state
	p.mutex.Lock()
	b, ok := p.cachedOnline[service.Id]
	if !ok {
		b = true
	}
	t2, ok := p.cachedDowntime[service.Id]
	if !ok {
		t2 = time.Now()
	}
	p.cachedOnline[service.Id] = true
	p.mutex.Unlock()
	if !b {
		go p.sendNotify(fmt.Sprintf("Service #%d '%s' is back online after %s downtime\n%s", service.Id, service.Name, t.Sub(t2).Truncate(time.Second), service.Domain))
	}
}

func (p *Ping) makeFailure(service structure.Service, resp *http.Response, timings *httpstat.Result) {
	_, _ = io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	t := time.Now()
	z := fmt.Sprintf("Expected status code %d but got %d", service.ExpectedStatus, resp.StatusCode)
	_, err := p.engine.Insert(structure.Failure{
		Issue:     z,
		ErrorCode: int64(resp.StatusCode),
		Service:   service.Id,
		PingTime:  uint64(timings.DNSLookup),
		Reason:    "status_code",
	})
	if err != nil {
		log.Printf("[Ping::makeFailure(Service %d)] Database error: %s\n", service.Id, err)
	}

	// update internal cache state
	p.mutex.Lock()
	b, ok := p.cachedOnline[service.Id]
	if !ok {
		b = false
	}
	if b {
		p.cachedOnline[service.Id] = false
		p.cachedDowntime[service.Id] = t
		go p.sendNotify(fmt.Sprintf("Service #%d '%s' has gone offline: %s\n%s", service.Id, service.Name, z, service.Domain))
	}
	p.mutex.Unlock()
}

func (p *Ping) makeError(service structure.Service, e error, timings *httpstat.Result) {
	var a uint64
	if timings != nil {
		a = uint64(timings.DNSLookup)
	}
	t := time.Now()
	_, err := p.engine.Insert(structure.Failure{
		Issue:     e.Error(),
		ErrorCode: 0,
		Service:   service.Id,
		PingTime:  a,
		Reason:    "error",
	})
	if err != nil {
		log.Printf("[Ping::makeError(Service %d)] Database error: %s\n", service.Id, err)
	}

	// update internal cache state
	p.mutex.Lock()
	b, ok := p.cachedOnline[service.Id]
	if !ok {
		b = false
	}
	if b {
		p.cachedOnline[service.Id] = false
		p.cachedDowntime[service.Id] = t
		go p.sendNotify(fmt.Sprintf("Service #%d '%s' has gone offline due to an internal error\n%s", service.Id, service.Name, service.Domain))
	}
	p.mutex.Unlock()
}
