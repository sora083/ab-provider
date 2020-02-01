package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
	"fmt"
	//"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sora083/ab-provider/api"
	"github.com/sora083/ab-provider/model"
	validator "gopkg.in/go-playground/validator.v9"
)

type Url struct {
	Qr     string `json:"qr,omitempty"`
	Mobile string `json:"mobile,omitempty"`
	Pc     string `json:"pc,omitempty"`
}

type PriceType struct {
	Min int64  `json:"min,omitempty"`
	Max int64  `json:"max,omitempty"`
	Commission int64  `json:"commsission,omitempty"`
}

// type City struct {
// 	Code    string `json:"code,omitempty"`
// 	Name    string `json:"name,omitempty"`
// 	NonStop string `json:"nonstop,omitempty"`
// }

type Code struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

type Ticket struct {
	ID             string `json:"id,omitempty"`
	Title          string `json:"title,omitempty"`
	LastUpdate     string `json:"last_update,omitempty"`
	Airline        []Code `json:"airline,omitempty"`
	AirlineType    string `json:"airline_type,omitempty"`
	AirlineSummary string `json:"airline_summary,omitempty"`
	City           []Code `json:"city,omitempty"`
	TermMin        string `json:"term_min,omitempty"`
	TermMax        string `json:"term_max,omitempty"`
	SeatClass      Code   `json:"seat_class,omitempty"`
	DeptTIme       string `json:"dept_time,omitempty"`
	TripType       Code   `json:"trip_type,omitempty"`
	Price          PriceType  `json:"price,omitempty"`
	Brand          Code   `json:"brand,omitempty"`
	Urls           Url    `json:"urls,omitempty"`
}

type SearchResult struct {
	ResultsReturned  string   `json:"results_returned"`
	ResultsStart     int64    `json:"results_start"`
	ResultsAvailable string   `json:"results_available"`
	Ticket           []Ticket `json:"ticket,omitempty"`
}

type SearchResults struct {
	Results SearchResult `json:"results"`
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

func GetCity(cities []Code) string {
	//return strings.Join(cities, ",")
	// str string := nil
	// for i, s := range cities {
	// 	str = str + s.Name + "(" + s.Code + ") "
	// }
	// return str
	return "City"
}

func GetPrice(price PriceType) string {
	//return Price.Min + "〜" + Price.Max + "(" + Price.Commission + ")"
	return fmt.Sprintf("%d 〜 %d(%d)", price.Min, price.Max, price.Commission)
}

func main() {
	e := echo.New()
	e.Validator = &Validator{validator: validator.New()}
	// funcs := template.FuncMap{
	// 	"encode_json": func(v interface{}) string {
	// 		b, _ := json.Marshal(v)
	// 		return string(b)
	// 	},
	// }
	e.Renderer = &Renderer{
		//templates: template.Must(template.New("").Delims("[[", "]]").Funcs(funcs).ParseGlob("templates/*.html")),
		//templates: template.Must(template.New("").Delims("[[", "]]").ParseGlob("templates/*.html")),
		templates: template.Must(template.New("").ParseGlob("templates/*.html")),
		//templates: template.Must(template.ParseGlob("templates/*.html")),
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

	response, err := api.FetchTicketsInfos(req)
	if err != nil {
		log.Fatal(err)
	}

	var results SearchResults
	json.Unmarshal(response, &results)

	tickets := results.Results.Ticket
	//log.Print("RESPONSE: ", tickets)

	return c.Render(200, "search_results.html", tickets)
}
