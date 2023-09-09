package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wakieu/drtbox/api"
	"github.com/wakieu/drtbox/database"
	"github.com/wakieu/drtbox/pages"
)

func main() {

	db, err := sql.Open("sqlite3", "./drtbox.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	SQL := "CREATE TABLE box (boxpath TEXT NOT NULL PRIMARY KEY, text TEXT);"
	_, err = db.Exec(SQL)
	if err != nil && err.Error() != "table box already exists" {
		log.Printf("%q\n", err)
		return
	}

	landingPageTemplate, err := pages.NewTemplate("pages/landing_page.html")
	if err != nil {
		panic(err)
	}
	boxPageTemplate, err := pages.NewTemplate("pages/box_page.html")
	if err != nil {
		panic(err)
	}

	boxRepo := database.NewBoxRepository(db)

	//Setup Api Server
	apiHandler := api.NewHandler(boxRepo)
	apiServer := http.NewServeMux()
	apiServer.HandleFunc("/", apiHandler.ServeHTTP)

	//Setup Page Server
	pageHandler := pages.NewHandler(boxRepo, landingPageTemplate, boxPageTemplate)
	pageServer := http.NewServeMux()
	pageServer.HandleFunc("/", pageHandler.ServeHTTP)

	//Running Api Server
	go func() {
		log.Println("ApiServer started on: http://localhost:3131")
		http.ListenAndServe("localhost:3131", apiServer)
	}()

	//Running Page Server
	log.Println("PageServer started on: http://localhost:3030")
	http.ListenAndServe("localhost:3030", pageServer)

}
