package api

import (
	//"fmt"
	"log"

	//"github.com/go-resty/resty/v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/sora083/ab-provider/model"
)

// func FetchTicketsInfos(req *model.SearchReq) *resty.Response {
// 	// Create a Resty Client
// 	client := resty.New()

// 	log.Println("Before start to call API", req)

// 	resp, err := client.R().
// 		EnableTrace().
// 		SetQueryParams(map[string]string{
// 			"dept":   req.Departure,
// 			"city":   req.Arrival,
// 			"ymd":    req.DepartureDate,
// 			"key":    "",
// 			"format": "json",
// 		}).
// 		Get("http://webservice.recruit.co.jp/ab-road-air/ticket/v1/")

// 	log.Println("After call API")

// 	// Explore response object
// 	fmt.Println("Response Info:")
// 	fmt.Println("Error      :", err)
// 	fmt.Println("Status Code:", resp.StatusCode())
// 	fmt.Println("Status     :", resp.Status())
// 	fmt.Println("Time       :", resp.Time())
// 	fmt.Println("Received At:", resp.ReceivedAt())
// 	// fmt.Println("Body       :\n", resp)
// 	fmt.Println()

// 	return resp
// }

func FetchTicketsInfos(req *model.SearchReq) ([]byte, error) {

	values := url.Values{}
	values.Add("dept", req.Departure)
	values.Add("city", req.Arrival)
	values.Add("ymd", req.DepartureDate)
	values.Add("key", "")
	values.Add("format", "json")

	httpUrl := "http://webservice.recruit.co.jp/ab-road-air/ticket/v1/"

	timeout := time.Duration(30 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(httpUrl + "?" + values.Encode())
	if err != nil {
		return nil, err
	}

	// 関数を抜ける際に必ずresponseをcloseするようにdeferでcloseを呼ぶ
	defer resp.Body.Close()

	body, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		log.Fatal(error)
	}

	return body, nil
}
