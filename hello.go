package commandone

import (
	"html/template"
	"net/http"
	"time"
	"appengine"
	"appengine/datastore"
)

type Reading struct {
	Timestamp time.Time
	Value     string
}

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/logit", logit)
}


func root(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	q := datastore.NewQuery("Reading").Ancestor(readingKey(c))
	readings := make([]Reading, 0, 10)
	if _, err := q.GetAll(c, &readings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := readingTemplate.Execute(w, readings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func logit(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	reading := Reading{
		Timestamp: time.Now(),
		Value: "98 degrees",
	}

	key := datastore.NewIncompleteKey(c, "Reading", readingKey(c))

	_, err := datastore.Put(c, key, &reading)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func readingKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Reading", "default_reading", 0, nil)
}

var readingTemplate = template.Must(template.New("book").Parse(`
<html>
  <head>
    <title>Go Guestbook</title>
  </head>
  <body>
  {{range .}}
      {{with .Timestamp}}
        <p><b>{{.}}</b> wrote:</p>
      {{else}}
        <p>An anonymous person wrote:</p>
      {{end}}
      <pre>{{.Value}}</pre>
    {{end}}
    <form action="/logit" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>
  </body>
</html>
`))
