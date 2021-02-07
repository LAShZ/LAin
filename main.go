package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"lawf"
)

type student struct {
	Name string
	Age  uint8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func v2Middlewares() lawf.HandlerFunc {
	return func(c *lawf.Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := lawf.New()
	r.Use(lawf.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "LA", Age: 20}
	stu2 := &student{Name: "LZH", Age: 21}

	r.GET("/", func(c *lawf.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *lawf.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", lawf.Head{
			"title":  "lawf",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	r.GET("/date", func(c *lawf.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", lawf.Head{
			"title": "lawf",
			"now":   time.Date(2021, 2, 7, 0, 0, 0, 0, time.UTC),
		})
	})

	r.GET("/hello", func(c *lawf.Context) {
		// expect /hello?name=LA
		c.String(200, "Hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *lawf.Context) {
		// expect /hello/Lavch
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.POST("/login", func(c *lawf.Context) {
		c.JSON(http.StatusOK, lawf.Head{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/hello", func(c *lawf.Context) {
			c.String(http.StatusOK, "hello %s, you're at v1: %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	v2.Use(v2Middlewares())
	{
		v2.GET("/hello/:name", func(c *lawf.Context) {
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
