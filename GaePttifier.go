package gaepttifer

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

const (
	ruleDbName = "RuleList"
)

type Rule struct {
	Board    string
	TitleKey string
	Date     time.Time
	Email    string
}

func init() {
	http.HandleFunc("/", rootPageHandler)
	http.HandleFunc("/admin/setup", setupPageHandler)
	http.HandleFunc("/admin/crawling", crawlingHandler)
}

func ruleListKey(ctx appengine.Context) *datastore.Key {
	return datastore.NewKey(ctx, "CrawlingRules", "default_rulelist", 0, nil)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, err string) {
	w.WriteHeader(status)
	page := template.Must(template.ParseFiles(
		"static/_base.html",
		"static/baseError.html",
	))

	switch status {
	case http.StatusNotFound:
		page = template.Must(template.ParseFiles(
			"static/_base.html",
			"static/404.html",
		))
	case http.StatusInternalServerError:
		page = template.Must(template.ParseFiles(
			"static/_base.html",
			"static/500.html",
		))
	case http.StatusUnauthorized:
		page = template.Must(template.ParseFiles(
			"static/_base.html",
			"static/401.html",
		))
	}

	if err := page.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func rootPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound, "")
		return
	}

	page := template.Must(template.ParseFiles(
		"static/_base.html",
		"static/index.html",
	))

	if err := page.Execute(w, nil); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

// Setup crawling rules, only admin user.
func setupPageHandler(w http.ResponseWriter, r *http.Request) {
	// handle user authority by myself.
	ctx := appengine.NewContext(r)
	// if current login user is not admin.
	if !user.IsAdmin(ctx) {
		errorHandler(w, r, http.StatusUnauthorized, "")
		return
	}

	if r.Method == "POST" {
		rule := Rule{
			Board:    r.FormValue("board"),
			TitleKey: r.FormValue("title_key"),
			Date:     time.Now(),
		}

		// Get then set login user's Email into database for crawler to send
		if u := user.Current(ctx); u != nil {
			rule.Email = u.Email
		} else {
			errorHandler(w, r, http.StatusInternalServerError, "")
		}

		key := datastore.NewIncompleteKey(ctx, ruleDbName, ruleListKey(ctx))
		if _, err := datastore.Put(ctx, key, &rule); err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)

	} else if r.Method == "GET" {
		page := template.Must(template.ParseFiles(
			"static/_base.html",
			"static/admin/setup.html",
		))

		if err := page.Execute(w, nil); err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}

	} else {
		errorHandler(w, r, http.StatusNotFound, "")
	}
}

// For cron schedule to call to do the crawling jobs.
func crawlingHandler(w http.ResponseWriter, r *http.Request) {
	// query rules from database
	ctx := appengine.NewContext(r)
	q := datastore.NewQuery(ruleDbName).Ancestor(ruleListKey(ctx)).Order("-Date").Limit(10)

	rules := []Rule{}
	if _, err := q.GetAll(ctx, &rules); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// pass through to crawlers
	for i := 0; i < len(rules); i++ {
		go crawlers[i].Crawling(&rules[i])
	}
}
