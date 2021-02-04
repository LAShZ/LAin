package main

import (
	"log"
	"net/http"

	"lawf"
)

func main() {
	r := lawf.New()
	r.GET("/", func(c *lawf.Context) {
		c.HTML(http.StatusOK, "<h1>Hello LA</h1>")
	})

	r.GET("/hello", func(c *lawf.Context) {
		// expect /hello?name=LA
		c.String(200, "Hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *lawf.Context) {
		c.JSON(http.StatusOK, lawf.Head{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	if err := r.Run(":9999"); err != nil {
		log.Fatal(err)
	}
}
