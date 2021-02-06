package main

import (
	"log"
	"net/http"

	"lawf"
)

func main() {
	r := lawf.New()
	r.GET("/", func(c *lawf.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})

	r.GET("/hello", func(c *lawf.Context) {
		// expect /hello?name=LA
		c.String(200, "Hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *lawf.Context) {
		// expect /hello/Lavch
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *lawf.Context) {
		c.JSON(http.StatusOK, lawf.Head{
			"filepath": c.Param("filepath"),
		})
	})

	r.POST("/login", func(c *lawf.Context) {
		c.JSON(http.StatusOK, lawf.Head{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *lawf.Context) {
			c.HTML(http.StatusOK, "<h1>Hello World!</h1>")
		})

		v1.GET("/hello", func(c *lawf.Context) {
			c.String(http.StatusOK, "hello %s, you're at v1: %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	{
		v2.GET("hello/:name", func(c *lawf.Context) {
			c.String(http.StatusOK, "hello %s, you're at v2: %s\n", c.Param("name"), c.Path)
		})

		v2.POST("/login", func(c *lawf.Context) {
			c.JSON(http.StatusOK, lawf.Head{
				"username": c.PostForm("v2:username"),
				"password": c.PostForm("v2:password"),
			})
		})
	}

	if err := r.Run(":9999"); err != nil {
		log.Fatal(err)
	}
}
