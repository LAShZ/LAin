package lawf

import (
	"net/http"
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
	c := newContext(w, req)
	engine.router.handle(c)
}

func (group *RouterGroup) addRoute(method string, path string, handler HandlerFunc) {
	group.engine.router.addRoute(method, path, handler)
}

// GET defines the method to add "GET" request
func (group *RouterGroup) GET(path string, handler HandlerFunc) {
	group.engine.addRoute("GET", path, handler)
}

// POST defines the method to add "POST" request
func (group *RouterGroup) POST(path string, handler HandlerFunc) {
	group.addRoute("POST", path, handler)
}

//  Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
