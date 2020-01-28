package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sora083/ab-provider/api"
)

type Renderer struct {
	templates *template.Template
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	funcs := template.FuncMap{
		"encode_json": func(v interface{}) string {
			b, _ := json.Marshal(v)
			return string(b)
		},
	}
	e.Renderer = &Renderer{
		templates: template.Must(template.New("").Delims("[[", "]]").Funcs(funcs).ParseGlob("templates/*.html")),
	}

	// 全てのリクエストで差し込みたいミドルウェア（ログとか）はここ
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ルーティング
	e.GET("/", Index)
	e.GET("/search", Search)

	// サーバー起動
	e.Start(":8080")
}

func Index(c echo.Context) error {
	return c.Render(200, "index.html", echo.Map{})
}

func Search(c echo.Context) error {
	log.Println("request")

	var departure = c.QueryParam("departure")
	var arrival = c.QueryParam("arrival")
	var departureDate = c.QueryParam("departureDate")
	// log.info("request: {}", request.toString())

	log.Println("departure: ", departure)
	log.Println("arrival: ", arrival)
	log.Println("departureDate: ", departureDate)

	// if result.hasErrors() {
	// 	log.error("validation error")
	// }

	api := api.FetchTicketsInfos()
	log.Println("response:", api)
	// if err := api.Get(); err != nil {
	// 	return err
	// }

	// model.addAttribute("results", searchResults)
	return c.Render(200, "searchResults.html", echo.Map{})
}

func featchSearchResults() {

}
