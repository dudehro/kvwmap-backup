package delivery

import (
	"github.com/kvwmap-backup/usecases"
	"github.com/kvwmap-backup/models"
	"net/http"
//	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"errors"
//	"encoding/json"
//	"strconv"
)

func URIHandler() {
	r := mux.NewRouter()
	r.HandleFunc("/", list_backup_configs)
        r.HandleFunc("/backup-config", list_backup_configs)
        r.HandleFunc("/backup-config/{File}", show_backup_config)
        r.HandleFunc("/backup-config/{File}/tar/new", tar_new)
	r.HandleFunc("/backup-config/{File}/tar/save", tar_save)
	r.HandleFunc("/backup-config/{File}/tar/edit/{Id}", tar_edit)
	r.HandleFunc("/backup-config/{File}/tar/delete/{Id}", tar_delete)

	//statische Inhalte, CSS+JS
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./interface/templates/static"))))

	r.NotFoundHandler = http.HandlerFunc(httpNotFound)

        http.ListenAndServe(":8082", r)
}

// /backup-config
func list_backup_configs(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("./interface/templates/configs.html")
	if err != nil {
		httpError(w, r, err)
	} else {
		templ.Execute(w, api_usecases.List_backup_configs())
	}
}

// /backup-config/{name}
func show_backup_config(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templ, err := template.ParseFiles("./interface/templates/config.html")
	if err != nil {
		httpError(w, r, err)
	}

	var c structs.HTMLTemplateData
	c.Vars = vars
	c.Backup, err = api_usecases.Load_backup_config(vars["File"])
	if err != nil {
		httpError(w, r, err)
	} else {
		templ.Execute(w, c)
	}
}

// /backup-config/{name}/tar/new
func tar_new(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var htmldata structs.HTMLTemplateData
	htmldata.Vars = vars
	templ, err := template.ParseFiles("./interface/templates/taritem.html")
	if err != nil {
		httpError(w, r, err)
	}
	templ.Execute(w, htmldata)
}

// /backup-config/{File}/tar/save
func tar_save(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var taritem structs.TarItem
	taritem.Source = r.PostFormValue("source")
	taritem.TargetName = r.PostFormValue("target")
	taritem.Exclude = r.PostFormValue("exclude")

	request := api_usecases.Save_TarItem(taritem, vars["File"])

	if request.Success {
		httpSuccess(w,r)
	} else {
		if len(request.Errors) < 1 {
			httpErrorStr(w, r, "unbekannter Fehler")
		} else {
			httpErrorStr(w, r, request.Errors[0])
		}
	}

}

// /backup-config/{File}/tar/edit
func tar_edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
        templ, err := template.ParseFiles("./interface/templates/taritem.html")
        if err != nil {
                httpError(w, r, err)
        }
        taritem, request := api_usecases.Load_Taritem(vars["File"], vars["Id"])
	if !request.Success {
		httpErrorStr(w, r, request.Errors[0])
	}
	vars["source"] = taritem.Source
	vars["targetname"] = taritem.TargetName
	vars["exclude"] = taritem.Exclude

	var htmldata structs.HTMLTemplateData
	htmldata.Vars = vars
	

	templ.Execute(w, htmldata)

}

func tar_delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	request := api_usecases.Delete_TarItem(vars["File"], vars["Id"])
	if request.Success {
		httpSuccess(w, r)
	} else {
		httpErrorStr(w, r, request.Errors[0])
	}
}

func httpSuccess(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./interface/templates/status/success.html")
}

func httpNotFound(w http.ResponseWriter, r *http.Request) {
//	w.WriteHeader(http.StatusNotFound)
	http.ServeFile(w, r, "./interface/templates/status/notfound.html")
}

func httpErrorPanic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./interface/templates/status/panic.html")
}

func httpErrorStr(w http.ResponseWriter, r *http.Request, e string) {
	err := errors.New(e)
	httpError(w,r,err)
}

func httpError(w http.ResponseWriter, r *http.Request, e error) {
        templ, err := template.ParseFiles("./interface/templates/status/error.html")
	var error_map map[string]string
	error_map = make(map[string]string)
	error_map["Error"] = e.Error()
        if err != nil {
		log.Println("Fehler beim parsen: "+err.Error())
		httpErrorPanic(w, r)
        } else {
	        templ.Execute(w, error_map)
	}
}
