package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
		templates: template.Must(template.New("").Delims("[[", "]]").Funcs(funcs).ParseGlob("templates/*.tmpl")),
	}

	// 全てのリクエストで差し込みたいミドルウェア（ログとか）はここ
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ルーティング
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index.tmpl", echo.Map{})
	})

	e.GET("/search", func(c echo.Context) error {
		log.Println("request")
		// log.info("request: {}", request.toString())
		// if result.hasErrors() {
		// 	log.error("validation error")
		// }
		// List < SearchResult > searchResults = flightSearchService.flightSearch(request)
		// log.info("response: {}", searchResults.toString())
		log.Println("response")
		// model.addAttribute("results", searchResults)
		return c.Render(200, "searchResults.tmpl", echo.Map{})
	})

	// サーバー起動
	e.Start(":8080")
}
