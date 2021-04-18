package main

// $ curl http://localhost:9999/
// URL.Path = "/"
// $ curl http://localhost:9999/hello
// Header["Accept"] = ["*/*"]
// Header["User-Agent"] = ["curl/7.54.0"]
// curl http://localhost:9999/world
// 404 NOT FOUND: /world

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"gee"
)

type student struct {
	Name string
	Age int
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%2d", year, month, day)
}

func onlyForV2() gee.HandlerFunc{
	return func(c *gee.Context) {
		t := time.Now()
		//if a server error occurred
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")
	stu1 := &student{Name: "jiawei", Age: 10}
	stu2 := &student{Name: "jack", Age: 19}
	r.Get("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gee.H{
			"title" : "gee",
			"stuArr" : [2]*student{stu1, stu2},
		})
	})

	r.Get("/date", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now": time.Date(2012, 4, 18, 12,34,43,33, time.UTC),
		})
	})


	r.Get("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, youre at %s \n", c.Query("name"), c.Path)
	})
	r.Get("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})
	r.Get("/hello/*/ztest", func(c *gee.Context) {
		c.String(http.StatusOK, "ztest, you're at %s\n",  c.Path)
	})


	r.Get("/assets/*filepath", func(c *gee.Context){
		c.JSON(http.StatusOK, gee.H{"filepath":c.Param("filepath")})
	})

	r.Post("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostFrom("username"),
			"password": c.PostFrom("password"),
		})
	})


	v2 := r.Group("v2")
	v2.Use(onlyForV2())
	{
		v2.Post("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostFrom("username"),
				"password": c.PostFrom("password"),
			})
		})
		v2.Get("/home/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, youre at %s \n", c.Param("name"), c.Path)
		})
	}
	//todo 为啥访问不了涅
	r.Static("/assets", "/Users/jiawei/testFile")

	r.Run(":9999")
}