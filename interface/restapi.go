package restapi

import (
	"github.com/kvwmap-backup/usecases"
	"net/http"
//	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
)

func URIHandler() {
	r := mux.NewRouter()
        r.HandleFunc("/backup_config", list_backup_configs)
        r.HandleFunc("/backup_config/{file}", show_backup_config)

        http.ListenAndServe(":8082", r)
}

// /backup-config
func list_backup_configs(w http.ResponseWriter, r *http.Request) {
//	for _,file := range api_usecases.List_backup_configs() {
//		fmt.Fprintf(w, "%s", file)
//	}

	templ, err := template.ParseFiles("./interface/templates/configs.html")
	if err != nil {
		log.Fatal(err)
	}
	templ.Execute(w, api_usecases.List_backup_configs())
}

// /backup-config/{name}
func show_backup_config(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
//	fmt.Fprintf(w,"Requested URI is %s", r.URL.Path)
//	fmt.Fprintf(w,"config-id is %s", vars["file"])

	templ, err := template.ParseFiles("./interface/templates/config.html")
	if err != nil {
		log.Fatal(err)
	}
	templ.Execute(w, api_usecases.Load_backup_config(vars["file"]))

}
