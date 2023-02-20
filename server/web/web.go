package web

import (
	"code.mrmelon54.com/melon/status/server/ping"
	"code.mrmelon54.com/melon/status/server/structure"
	"code.mrmelon54.com/melon/status/server/utils"
	"encoding/gob"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"xorm.io/xorm"
)

type Web struct {
	oauthClient  *oauth2.Config
	stateManager *utils.StateManager
	pageTitle    string
	engine       *xorm.Engine
	ownerSub     string
	originUrl    string
	resourceUrl  string
	pingService  *ping.Ping
}

func New(engine *xorm.Engine, p *ping.Ping) *Web {
	gob.Register(structure.OpenIdMeta{})
	sessionStore := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	return &Web{
		stateManager: utils.NewStateManager(sessionStore),
		pageTitle:    os.Getenv("TITLE"),
		engine:       engine,
		ownerSub:     os.Getenv("OWNER"),
		originUrl:    os.Getenv("ORIGIN_URL"),
		resourceUrl:  os.Getenv("RESOURCE_URL"),
		pingService:  p,
	}
}

func (w *Web) SetupWeb() *http.Server {
	gob.Register(new(utils.WebServiceKeyType))

	router := mux.NewRouter()
	router.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		_, _ = rw.Write([]byte("Status Page API Endpoint\n"))
	})
	router.HandleFunc("/status", w.all)
	router.HandleFunc("/status/{service}", w.getServiceStatus)
	router.HandleFunc("/maintenance", w.all)

	return &http.Server{
		Addr:    os.Getenv("LISTEN"),
		Handler: router,
	}
}

type groupServiceBeanTriple struct {
	Group   structure.Group   `xorm:"extends"`
	Service structure.Service `xorm:"extends"`
}

func (w *Web) all(rw http.ResponseWriter, _ *http.Request) {
	var bean []groupServiceBeanTriple
	err := w.engine.Where("group.public = 1 and service.public = 1").Select("group.name, group.id, group.order, service.name, service.id, service.order").
		Table(&structure.Group{}).
		Join("INNER", &structure.Service{}, "group.id = service.group_id").
		Asc("group.order", "service.order").
		Find(&bean)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Failed to read from database", http.StatusInternalServerError)
		return
	}

	outputGroups := make([]*structure.Group, 0)

	var currentGroup *structure.Group
	for _, i := range bean {
		if currentGroup == nil || currentGroup.Id != i.Group.Id {
			a := i.Group
			a.Services = make([]structure.Service, 0)
			outputGroups = append(outputGroups, &a)
			currentGroup = &a
		}
		currentGroup.Services = append(currentGroup.Services, i.Service)
	}

	// Encode output
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(outputGroups)
	if err != nil {
		http.Error(rw, "Failed to encode data", http.StatusInternalServerError)
		return
	}
}

type RenderGroup struct {
	D        structure.Group
	Services []RenderService
}

type RenderService struct {
	D        structure.Service
	State    byte
	Beans    [92]byte
	CheckAgo int
}
