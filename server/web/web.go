package web

import (
	"bytes"
	"code.mrmelon54.com/melon/status/server/ping"
	"code.mrmelon54.com/melon/status/server/res"
	"code.mrmelon54.com/melon/status/server/structure"
	"code.mrmelon54.com/melon/status/server/utils"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"xorm.io/xorm"
)

const (
	LoginFrameStart = "<!DOCTYPE html><html><head><script>window.opener.postMessage({user:"
	LoginFrameEnd   = "},\"%s\");window.close();</script></head></html>"
	CheckFrameStart = "<!DOCTYPE html><html><head><script>window.onload=function(){window.parent.postMessage({user:"
	CheckFrameEnd   = "},\"%s\");}</script></head></html>"
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
		oauthClient: &oauth2.Config{
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
			Scopes:       []string{"openid"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  os.Getenv("AUTHORIZE_URL"),
				TokenURL: os.Getenv("TOKEN_URL"),
			},
			RedirectURL: os.Getenv("REDIRECT_URL"),
		},
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
		open, err := res.Open("index.html")
		if err != nil {
			http.NotFound(rw, req)
			return
		}
		all, err := io.ReadAll(open)
		if err != nil {
			http.NotFound(rw, req)
			return
		}
		out := bytes.Replace(all, []byte("<!-- production inject -->"), []byte(fmt.Sprintf("<script>const process = {env:{SVELTE_APP_API_URL:\"%s\"}};</script>", w.originUrl)), 1)
		_, _ = rw.Write(out)
	})
	router.HandleFunc("/login", w.stateManager.SessionWrapper(w.loginPage))
	router.HandleFunc("/check", w.stateManager.SessionWrapper(w.checkPage))
	router.HandleFunc("/admin", w.stateManager.SessionWrapper(w.adminCheckWrapper(func(rw http.ResponseWriter, req *http.Request, state *utils.State, isAdmin bool) {
		if isAdmin {
			_, _ = rw.Write([]byte("is admin\n"))
		} else {
			_, _ = rw.Write([]byte("not admin\n"))
		}
		//vars := mux.Vars(req)
		//generateExamTimetable(rw, engine, vars["course"], vars["year"], 0.5)
	})))
	router.NotFoundHandler = http.HandlerFunc(res.Handler)

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			handler.ServeHTTP(rw, req)
		})
	})
	apiRouter.HandleFunc("/all", w.stateManager.SessionWrapper(w.adminCheckWrapper(w.all)))
	apiRouter.HandleFunc("/groups", w.stateManager.SessionWrapper(w.adminCheckWrapper(w.allGroups)))
	apiRouter.HandleFunc("/group", w.stateManager.SessionWrapper(w.adminCheckWrapper(w.getGroup)))
	apiRouter.HandleFunc("/services", w.stateManager.SessionWrapper(w.adminCheckWrapper(w.allServices)))
	apiRouter.HandleFunc("/service", w.stateManager.SessionWrapper(w.adminCheckWrapper(w.getService)))
	apiRouter.HandleFunc("/status", w.stateManager.SessionWrapper(w.adminCheckWrapper(w.getServiceStatus)))
	apiRouter.HandleFunc("/hits", w.stateManager.SessionWrapper(w.adminCheckWrapper(w.allHits)))
	apiRouter.HandleFunc("/failures", w.stateManager.SessionWrapper(w.adminCheckWrapper(w.allFailures)))

	return &http.Server{
		Addr:    os.Getenv("LISTEN"),
		Handler: router,
	}
}

type groupServiceBeanTriple struct {
	group   structure.Group   `xorm:"extends"`
	service structure.Service `xorm:"extends"`
	hit     structure.Hit     `xorm:"extends"`
	failure structure.Failure `xorm:"extends"`
}

func (w *Web) all(rw http.ResponseWriter, req *http.Request, state *utils.State, admin bool) {
	var bean []groupServiceBeanTriple
	r := w.engine.NewSession()
	if !admin {
		r = r.Where("group.public = 1 and service.public = 1")
	}
	err := r.Select("group.name, group.id, group.order, service.name, service.id, service.order, "+
		"failure.id, failure.created_at, hit.id, hit.created_at").
		Table(&structure.Group{}).
		Join("INNER", &structure.Service{}, "group.id = service.group_id").
		Join("INNER", &structure.Hit{}, "service.id = hit.service").
		Join("INNER", &structure.Failure{}, "service.id = failure.service").
		Asc("group.order", "service.order").
		Find(&bean)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Failed to read from database", http.StatusInternalServerError)
		return
	}

	// Encode output
	encoder := json.NewEncoder(rw)
	err = encoder.Encode(bean)
	if err != nil {
		http.Error(rw, "Failed to encode data", http.StatusInternalServerError)
		return
	}
}

func (w *Web) adminCheckWrapper(cb func(rw http.ResponseWriter, req *http.Request, state *utils.State, isAdmin bool)) func(rw http.ResponseWriter, req *http.Request, state *utils.State) {
	return func(rw http.ResponseWriter, req *http.Request, state *utils.State) {
		if myUser, ok := utils.GetStateValue[*structure.OpenIdMeta](state, utils.KeyUser); ok {
			if myUser == nil {
				cb(rw, req, state, false)
				return
			}
			if !myUser.Admin {
				cb(rw, req, state, false)
				return
			}
			cb(rw, req, state, true)
			return
		}
		cb(rw, req, state, false)
	}
}

func (w *Web) loginPage(rw http.ResponseWriter, req *http.Request, state *utils.State) {
	q := req.URL.Query()
	if q.Has("in_popup") {
		state.Put("login-in-popup", true)
	}
	if myUser, ok := utils.GetStateValue[*structure.OpenIdMeta](state, utils.KeyUser); ok {
		if myUser != nil {
			if doLoginPopup(rw, w.originUrl, state, myUser) {
				return
			}
			http.Redirect(rw, req, "/", http.StatusTemporaryRedirect)
			return
		}
	}
	if flowState, ok := utils.GetStateValue[uuid.UUID](state, utils.KeyState); ok {
		q := req.URL.Query()
		if q.Has("code") && q.Has("state") {
			if q.Get("state") == flowState.String() {
				exchange, err := w.oauthClient.Exchange(context.Background(), q.Get("code"))
				if err != nil {
					fmt.Println("Exchange token error:", err)
					return
				}
				state.Put(utils.KeyAccessToken, exchange.AccessToken)
				state.Put(utils.KeyRefreshToken, exchange.RefreshToken)

				buf := new(bytes.Buffer)
				req2, err := http.NewRequest(http.MethodGet, w.resourceUrl, buf)
				if err != nil {
					return
				}
				req2.Header.Set("Authorization", "Bearer "+exchange.AccessToken)
				do, err := http.DefaultClient.Do(req2)
				if err != nil {
					return
				}

				var meta structure.OpenIdMeta
				decoder := json.NewDecoder(do.Body)
				err = decoder.Decode(&meta)
				if err != nil {
					log.Println("Failed to decode external openid meta:", err)
					http.Error(rw, "500 Internal Server Error: Failed to fetch user info", http.StatusInternalServerError)
					return
				}
				meta.Admin = meta.Sub == w.ownerSub
				state.Put(utils.KeyUser, &meta)

				if doLoginPopup(rw, w.originUrl, state, &meta) {
					return
				}
				http.Redirect(rw, req, "/", http.StatusTemporaryRedirect)
				return
			}
			http.Error(rw, "OAuth flow state doesn't match\n", http.StatusBadRequest)
			return
		}
	}
	flowState := uuid.New()
	state.Put(utils.KeyState, flowState)
	http.Redirect(rw, req, w.oauthClient.AuthCodeURL(flowState.String(), oauth2.AccessTypeOffline), http.StatusTemporaryRedirect)
}

func (w *Web) checkPage(rw http.ResponseWriter, _ *http.Request, state *utils.State) {
	if myUser, ok := utils.GetStateValue[*structure.OpenIdMeta](state, utils.KeyUser); ok {
		if myUser != nil {
			exportUserDataAsJson(rw, w.originUrl, myUser, true)
			return
		}
	}
	rw.WriteHeader(http.StatusBadRequest)
}

func doLoginPopup(rw http.ResponseWriter, originUrl string, state *utils.State, meta *structure.OpenIdMeta) bool {
	if b, ok := utils.GetStateValue[bool](state, "login-in-popup"); ok {
		if b {
			exportUserDataAsJson(rw, originUrl, meta, false)
			return true
		}
	}
	return false
}

func exportUserDataAsJson(rw http.ResponseWriter, originUrl string, meta *structure.OpenIdMeta, checkMode bool) {
	start := LoginFrameStart
	end := LoginFrameEnd
	if checkMode {
		start = CheckFrameStart
		end = CheckFrameEnd
	}
	_, _ = rw.Write([]byte(start))
	encoder := json.NewEncoder(rw)
	_ = encoder.Encode(meta)
	_, _ = rw.Write([]byte(fmt.Sprintf(end, originUrl)))
}

func fillPage(w io.Writer, name, tempStr string, data any) {
	tmpl, err := template.New(name).Parse(tempStr)
	if err != nil {
		log.Printf("[Http::GeneratePage] Parse: %v\n", err)
		return
	}
	if data == nil {
		err = tmpl.Execute(w, struct{}{})
	} else {
		err = tmpl.Execute(w, data)
	}
	if err != nil {
		log.Printf("[Http::GeneratePage] Execute: %v\n", err)
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
