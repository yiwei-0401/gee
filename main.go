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
	"log"
	"net/http"
	"time"

	"gee"
)

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
	//day1
	/*r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	})*/
	r.Get("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>hello Gee</h1>")
	})
	//day2
	/*r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		for k,v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})*/
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

	v1 := r.Group("/v1")
	fmt.Println(v1)
	{
		v1.Get("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee V1</h1>")
		})
		v1.Get("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "hello V1 %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

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

	r.Run(":9999")
}