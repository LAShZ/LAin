package lawf

import (
	"log"
	"net/http"
	"strings"
)

// HandlerFunc defines the request handler used by lawf
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // store all groups
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // middlewares of the group
	parent      *RouterGroup  //parent of current group
	engine      *Engine       // engine of the webframework
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group creates a new Group remember all groups with the same engine
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}

func (group *RouterGroup) addRoute(method string, path string, handler HandlerFunc) {
	pattern := group.prefix + path
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add "GET" request
func (group *RouterGroup) GET(path string, handler HandlerFunc) {
	group.addRoute("GET", path, handler)
}

// POST defines the method to add "POST" request
func (group *RouterGroup) POST(path string, handler HandlerFunc) {
	group.addRoute("POST", path, handler)
}

//  Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// Use adds middlewares to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}
