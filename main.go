package main

// $ curl http://localhost:9999/
// URL.Path = "/"
// $ curl http://localhost:9999/hello
// Header["Accept"] = ["*/*"]
// Header["User-Agent"] = ["curl/7.54.0"]
// curl http://localhost:9999/world
// 404 NOT FOUND: /world

import (
	"net/http"

	"gee"
)

func main() {
	r := gee.New()
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

	r.Post("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostFrom("username"),
			"password": c.PostFrom("password"),
		})
	})

	r.Run(":9999")
}