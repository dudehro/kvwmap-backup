package delivery

import (
	"github.com/kvwmap-backup/usecases"
	"github.com/kvwmap-backup/models"
	"github.com/kvwmap-backup/configuration"
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"errors"
)

func URIHandler() {
	r := mux.NewRouter()
	r.HandleFunc("/", list_backup_configs)
        r.HandleFunc("/backup-config", list_backup_configs)
        r.HandleFunc("/backup-config/{File}", show_backup_config)
	//tar
        r.HandleFunc("/backup-config/{File}/tar/new", tar_new)
	r.HandleFunc("/backup-config/{File}/tar/save", tar_save)
	r.HandleFunc("/backup-config/{File}/tar/edit/{Id}", tar_edit)
	r.HandleFunc("/backup-config/{File}/tar/delete/{Id}", tar_delete)
	//mysql
	r.HandleFunc("/backup-config/{File}/mysql/new", mysql_new)
	r.HandleFunc("/backup-config/{File}/mysql/save", mysql_save)
	r.HandleFunc("/backup-config/{File}/mysql/edit/{Id}", mysql_edit)
	r.HandleFunc("/backup-config/{File}/mysql/delete/{Id}", mysql_delete)

	//statische Inhalte, CSS+JS
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./interface/templates/static"))))

	r.NotFoundHandler = http.HandlerFunc(httpNotFound)

	listenport := config.GetConfigValFor(config.KeyHTTPPort)
	log.Printf("starting server on port %s", listenport)
        err := http.ListenAndServe(":" + listenport, r)
	if err != nil {
		log.Fatal(err)
	}
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

	request := api_usecases.SaveItem(taritem, vars["File"])

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

// /backup-config/{File}/tar/delete/{Id}
func tar_delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	item := structs.TarItem{}
	request := api_usecases.DeleteItem(vars["File"], vars["Id"], item)
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
	log.Printf("httpNotFound requested URI %s", r.RequestURI)
	http.ServeFile(w, r, "./interface/templates/status/notfound.html")
}

func httpErrorPanic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./interface/templates/status/panic.html")
}

func httpErrorStr(w http.ResponseWriter, r *http.Request, e string) {
	log.Printf("httpErrorStr error:%s", e)
	err := errors.New(e)
	httpError(w,r,err)
}

func httpError(w http.ResponseWriter, r *http.Request, e error) {
	log.Printf("httpError error:%s", e.Error())
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

// /backup-config/{name}/mysql/new
func mysql_new(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var htmldata structs.HTMLTemplateData
	htmldata.Vars = vars
	templ, err := template.ParseFiles("./interface/templates/mysqlitem.html")
	if err != nil {
		httpError(w, r, err)
	}
	templ.Execute(w, htmldata)
}

// /backup-config/{File}/mysql/save
func mysql_save(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var mysqlitem structs.MysqlDumpItem
	mysqlitem.ContainerId = r.PostFormValue("containerid")
	log.Printf("ContainerId:%s",mysqlitem.ContainerId)
	mysqlitem.DbName = r.PostFormValue("dbname")
	mysqlitem.DockerNetwork = r.PostFormValue("dockernetwork")
	mysqlitem.TargetName = r.PostFormValue("targetname")

	request := api_usecases.SaveItem(mysqlitem, vars["File"])
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

// /backup-config/{File}/mysql/edit
func mysql_edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templ, err := template.ParseFiles("./interface/templates/mysql.html")
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


// /backup-config/{File}/mysql/delete/{Id}
func mysql_delete(w http.ResponseWriter, r *http.Request) {
	log.Println("mysql_delete")
	vars := mux.Vars(r)
	item := structs.MysqlDumpItem{}
	request := api_usecases.DeleteItem(vars["File"], vars["Id"], item)
	if request.Success {
		httpSuccess(w, r)
	} else {
		httpErrorStr(w, r, request.Errors[0])
	}
}
