package lawf

import (
	"fmt"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/Lavch")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.path != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if ps["name"] != "Lavch" {
		t.Fatal("name should be equal to \"Lavch\"")
	}

	fmt.Printf("matched path: %s, params[\"name\"]: %s\n", n.path, ps["name"])
}
