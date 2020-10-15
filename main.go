package webperforapi

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/goware/urlx"
	"github.com/m7shapan/njson"
)

//GetWebPerformanceScore obtains general lighthouse result score
func GetWebPerformanceScore(domain, apiKey string) (int, error) {
	fName := "GetWebPerformanceScore"
	log.Printf("[%s]: Starting the function", fName)
	defer log.Printf("[%s]: Ending the function", fName)

	type googleWPResult struct {
		Score float32 `njson:"lighthouseResult.categories.performance.score"`
	}
	var r googleWPResult
	d, err := urlx.Parse(domain)
	if err != nil {
		log.Printf("[%s]: Error occured when parsing the domain - %s. Error: %v", fName, domain, err)
		return -1, err
	}
	u := &url.URL{
		Scheme:   "https",
		Host:     "www.googleapis.com",
		Path:     "pagespeedonline/v5/runPagespeed",
		RawQuery: fmt.Sprintf("key=%s&url=%s", apiKey, d.String()),
	}
	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("[%s]: Error occured when doing get request for the domain - %s. Error: %v", fName, domain, err)
		return -1, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[%s]: Error occured when reading the response. Status of response: %s. Error: %v", fName, resp.Status, err)
		return -1, err
	}
	if err := njson.Unmarshal(data, &r); err != nil {
		log.Printf("[%s]: Error occured when tried to unmarshal the response. The response: %s. Error: %v", fName, data, err)
		return -1, err
	}
	return int(r.Score * 100.0), nil
}
