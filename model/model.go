package model

type SearchReq struct {
	Departure     string `form:"departure" query:"departure" validate:"required,len=3"`
	Arrival       string `form:"arrival" query:"arrival" validate:"required,len=3"`
	DepartureDate string `form:"departureDate" query:"departureDate" validate:"required,len=8"`
}
