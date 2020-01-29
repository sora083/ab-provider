package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sora083/ab-provider/api"
	"github.com/sora083/ab-provider/model"
	validator "gopkg.in/go-playground/validator.v9"
)

// type SearchReq struct {
// 	Departure     string `form:"departure" query:"departure" validate:"required,len=3"`
// 	Arrival       string `form:"arrival" query:"arrival" validate:"required,len=3"`
// 	DepartureDate string `form:"departureDate" query:"departureDate" validate:"required,len=8"`
// }

// type Ticket struct {
// 	ID          string `json:"id,omitempty"`
// 	title       string `json:"title,omitempty"`
// 	last_update string `json:"last_update,omitempty"`
// 	//ist<Airline> airline string `json:"id,omitempty"`
// 	airline_type    string `json:"airline_type,omitempty"`
// 	airline_summary string `json:"airline_summary,omitempty"`
// 	//DeptDetail deptDetail string `json:"dept_detail,omitempty"`
// 	// private CityNumber city_number string `json:"id,omitempty"`
// 	// private List<City> city string `json:"id,omitempty"`
// 	term_min   int64  `json:"term_min,omitempty"`
// 	term_max   int64  `json:"term_max,omitempty"`
// 	seat_class string `json:"seat_class,omitempty"`
// 	dept_time  string `json:"dept_time,omitempty"`
// 	trip_type  string `json:"trip_type,omitempty"`
// 	price      int64  `json:"price,omitempty"`
// 	bland      string `json:"bland,omitempty"`
// 	//urls string `json:"urls,omitempty"`
// }

// type SearchResult struct {
// 	resultsStart     int64  `json:"results_start,omitempty"`
// 	resultsReturned  int64  `json:"results_returned,omitempty"`
// 	resultsAvailable int64  `json:"results_available,omitempty"`
// 	tickets          Ticket `json:"ticket,omitempty"`
// }

type SearchResults struct {
	//results SearchResult `json:"results,omitempty"`
	results string `json:"results,omitempty"`
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

// TODO: https://qiita.com/RunEagler/items/ad79fc860c3689797ccc
func DateValidation(fl validator.FieldLevel) bool {
	_, err := time.Parse("20060102", fl.Field().String())
	if err != nil {
		return false
	}
	return true
}

type Renderer struct {
	templates *template.Template
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	e.Validator = &Validator{validator: validator.New()}
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

	req := &model.SearchReq{}
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Request is failed: "+err.Error())
	}
	log.Println("Request:", req)

	if err := c.Validate(req); err != nil {
		return c.String(http.StatusBadRequest, "Validate is failed: "+err.Error())
	}

	response := api.FetchTicketsInfos(req)
	log.Println("response:", response)
	log.Println("response.Body:", response.Body)
	log.Println("response.RawBody:", response.RawBody)
	//log.Println("response.RawResponse:", response.RawResponse)
	// if err := api.Get(); err != nil {
	// 	return err
	// }

	// model.addAttribute("results", searchResults)
	// var r SearchResults
	// json.Unmarshal(response.Body, &r)

	var err error
	var body []byte
	if body, err = json.Marshal(response); err != nil {
		return c.String(http.StatusBadRequest, "Parse is failed: "+err.Error())
	}
	log.Println("RES:", body)
	return c.Render(200, "searchResults.html", echo.Map{})
}
