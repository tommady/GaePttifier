// https://github.com/GoogleCloudPlatform/appengine-angular-gotodos/blob/master/gotodos.go

package gaepttifer

import (
    "html/template"
    "net/http"
    "time"

    "appengine"
    "appengine/datastore"
    "appengine/user"
)

type Rule struct {
    Board   string
    Keyword string
    Date    time.Time
}

func init() {
    http.HandleFunc("/", rootPageHandler)
    http.HandleFunc("/admin/setup", setupPageHandler)
    http.HandleFunc("/admin/crawling", crawlingHandler)
}

func ruleListKey(c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "CrawlingRules", "default_rulelist", 0, nil)
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
    c := appengine.NewContext(r)
    // if current login user is not admin.
    if !user.IsAdmin(c) {
        errorHandler(w, r, http.StatusUnauthorized, "")
        return
    }

    if r.Method == "POST" {
        ru := Rule{
                    Board:    r.FormValue("board"),
                    Keyword:  r.FormValue("keyword"),
                    Date:     time.Now(),
        }

        key := datastore.NewIncompleteKey(c, "Rule", ruleListKey(c))
        if _, err := datastore.Put(c, key, &ru); err != nil {
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
    // doing some crawlings
}
