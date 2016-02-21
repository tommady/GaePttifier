package gaepttifer

import (
	"net/http"

    "appengine"
	"appengine/datastore"
	"appengine/user"

	"github.com/PuerkitoBio/goquery"
)

const (
	pttBaseUrl         = "https://www.ptt.cc/bbs/"
	defaultParsingPage = "/index"
    resultDbName = "ResultList"
)

type Crawler struct {
	Email string
	Url   string
	Title string
}

func resultListKey(ctx appengine.Context) *datastore.Key {
	return datastore.NewKey(ctx, "CrawlingResults", "default_resultlist", 0, nil)
}

// for some specific board need over 18 years old check
func getPttRespond(url string) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle error
	}

	req.Header.Set("Cookie", "over18=1")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle error
	}

	return res
}

// do the crawling job then write results into database for the email reaper to handle.
func (crw *Crawler) Crawling(rule *Rule) {
    ctx := appengine.NewContext(r)
	crawlUrl := pttBaseUrl + rule.Board + defaultParsingPage

	doc, err := goquery.NewDocumentFromResponse(getPttRespond(crawlUrl))
	if err != nil {
	    // handle error
	}

	crw.Email = rule.Email
	crw.Url = "https://test, test"
	crw.Title = "tsettest"

    key := datastore.NewIncompleteKey(ctx, resultDbName, resultListKey(ctx))
    if _, err := datastore.Put(ctx, key, crw); err != nil {
        // handle error
    }

}
