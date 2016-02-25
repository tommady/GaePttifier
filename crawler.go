package gaepttifer

import (
	"appengine"
	"appengine/datastore"
)

const (
	ruleDbName         = "RuleList"
	pttBaseUrl         = "https://www.ptt.cc/bbs/"
	defaultParsingPage = "/index"
	resultDbName       = "ResultList"
)

type Rule struct {
	Board    string
	TitleKey string
	Email    string
}

type Result struct {
	Email string
	Url   string
	Title string
}

func defaultRuleList(ctx *appengine.Context) *datastore.Key {
	return datastore.NewKey(*ctx, "CrawlingRules", "default_rulelist", 0, nil)
}

func GetAllRules(ctx *appengine.Context) ([]Rule, error) {
	q := datastore.NewQuery(ruleDbName).Ancestor(defaultRuleList(ctx))

	rules := []Rule{}
	if _, err := q.GetAll(*ctx, &rules); err != nil {
		return nil, ReportError("Error: on getting rule into datastore", err)
	}

	return rules, nil
}

func (rule *Rule) Set(ctx *appengine.Context) error {
	key := datastore.NewIncompleteKey(*ctx, ruleDbName, defaultRuleList(ctx))
	if _, err := datastore.Put(*ctx, key, rule); err != nil {
		return ReportError("Error: on putting rule into datastore", err)
	}

	return nil
}

func (rule *Crawler) Crawling() {
	// doing crawling
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
