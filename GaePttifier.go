package gaepttifer

import (
	"appengine"
	"appengine/user"
	"html/template"
	"net/http"
)

func init() {
	http.HandleFunc("/", rootPageHandler)
	http.HandleFunc("/admin/setup", setupPageHandler)
	http.HandleFunc("/admin/crawling", crawlingHandler)
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
		}

		// Get then set login user's Email into database for crawler to send
		if u := user.Current(ctx); u != nil {
			rule.Email = u.Email
		} else {
			errorHandler(w, r, http.StatusInternalServerError, "")
		}

		if err := rule.Set(&ctx); err != nil {
			fmt.Fprintln(w, err)
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
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
	ctx := appengine.NewContext(r)

	crawlers, err := GetAllRules(&ctx)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
	}

	for i := 0; i < len(crawlers); i++ {
		go crawlers[i].Crawling()
	}
}
