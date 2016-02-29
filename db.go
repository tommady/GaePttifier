package gaepttifer

import (
	"appengine"
	"appengine/datastore"

	."github.com/tommady/pttifierLib"
)

const (
	ruleDbName         = "RuleList"
	resultDbName       = "ResultList"
)



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

func SetRule(ctx *appengine.Context, rule *Rule) error {
	key := datastore.NewIncompleteKey(*ctx, ruleDbName, defaultRuleList(ctx))
	if _, err := datastore.Put(*ctx, key, rule); err != nil {
		return ReportError("Error: on putting rule into datastore", err)
	}

	return nil
}
