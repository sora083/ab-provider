package api

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/sora083/ab-provider/model"
)

func FetchTicketsInfos(req *model.SearchReq) *resty.Response {
	// Create a Resty Client
	client := resty.New()

	log.Println("Before start to call API", req)

	resp, err := client.R().
		EnableTrace().
		SetQueryParams(map[string]string{
			"dept":   req.Departure,
			"city":   req.Arrival,
			"ymd":    req.DepartureDate,
			"key":    "",
			"format": "json",
		}).
		Get("http://webservice.recruit.co.jp/ab-road-air/ticket/v1/")

	log.Println("After call API")

	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("Error      :", err)
	fmt.Println("Status Code:", resp.StatusCode())
	fmt.Println("Status     :", resp.Status())
	fmt.Println("Time       :", resp.Time())
	fmt.Println("Received At:", resp.ReceivedAt())
	// fmt.Println("Body       :\n", resp)
	fmt.Println()

	return resp
}
