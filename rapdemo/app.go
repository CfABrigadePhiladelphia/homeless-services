//The RAP site is composed of three parts:
//-static pages - for the public
//-API - for RESTful CRUD ops on Datastore
//-logged in pages - as a web based way access the CRUD ops - not in the initial release
//-an import page for the bulk updates

//google app engine handles auth fairly well on its own

package rapdemo

import (
	"net/http"
	"time"

	"appengine"
)

const basePath = "rapdemo"

func init() {
	//basePath = "rapdemo"
	fs := http.FileServer(http.Dir(basePath + "/static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.Handle("/auth", appHandler(authdemo))

	//bulk import from csv
	http.Handle("/csvimport", appHandler(csvimport))

	//api
	http.Handle("/resources", appHandler(resources))

	//handles the templated but otherwise mostly static html pages
	http.Handle("/", appHandler(serveTemplate))
}

//The resource type is what most of the application will focus on.
type resource struct {
	ID int64

	//display fields
	Category         string `json:"Category"`
	OrganizationName string `json:"Organization Name"`
	Address          string `json:"Address"`
	ZipCode          string `json:"Zip Code"`
	Days             string `json:"Days"`
	TimeOpen         string `json:"Time: Open"`
	TimeClose        string `json:"Time: Close"`
	PeopleServed     string `json:"People Served"`
	Description      string `json:"Description"`
	PhoneNumber      string `json:"Phone Number"`
	Location         appengine.GeoPoint

	//audit fields
	LastUpdatedTime time.Time `datastore:",noindex"`
	LastUpdatedBy   string    `datastore:",noindex"`
	IsActive        bool
}

//following the error pattern suggested in the Go Blog
//http://blog.golang.org/error-handling-and-go

type appError struct {
	Error   error
	Message string
	Code    int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		c := appengine.NewContext(r)
		c.Errorf("%v", e.Error)
		http.Error(w, e.Message, e.Code)
	}
}
