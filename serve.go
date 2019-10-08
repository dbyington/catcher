package main

import (
	"net/http"

	"github.com/dbyington/pitcher"
	"github.com/sirupsen/logrus"
)

var W1 = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := r.Context().Value("log").(*logrus.Logger)
		log.Debug("My Middleware Calling next")
		next.ServeHTTP(w, r)
		log.Debug("My Middleware Done")
	})
}

var home = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value("log").(*logrus.Logger)
	log.Debug("home handler")
	w.Header().Set("content-type", "text/html")
	reply := `<h2>Welcome Home</h2>
<a href="/about"">About</a><br/><a href="/product/1234/detail">Product 1234</a>`
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(reply))
})

var about = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value("log").(*logrus.Logger)
	log.Debug("about handler")
	w.Header().Set("content-type", "text/html")

	reply := `<h2>About</h2>
<a href="/home"">Home</a><br/><a href="/product/1234/detail">Product 1234</a>`
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(reply))
})

var product = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value("log").(*logrus.Logger)
	log.Debug("product handler")
	params := r.Context().Value("parameters").(map[string]string)
	w.Header().Set("content-type", "text/html")

	reply := `<h2>Product `+params["id"]+` Detail</h2>
`+params["part"]+`<br/>`
	reply += r.URL.Query().Encode() + "<br/>"
	reply += `<a href="/about"">About</a><br/><a href="/home">home</a><a href="/product/new/fancy">New Products</a>`
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(reply))
})

func main() {
	app := pitcher.NewApp(":8888")

	app.Use(W1)

	app.GET("/", home)
	app.GET("/home", home)
	app.GET("/about", about)
	app.POST("/product", product)
	app.GET("/product/:id", product)
	app.PUT("/product/:id", product)
	app.GET("/product/:id/detail/:part", product)
	app.GET("/product/new/fancy", product)

	if err := app.ListenAndServe(); err != nil {
		app.Log.Errorf("error serving: %s", err)
	}
}
