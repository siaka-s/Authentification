package main

import (
	"Authentification/database"
	"Authentification/handler"
	"log"
	"net/http"
)

func main() {

	database.Opendb()

	defer database.CloseDB()

	fs := http.FileServer(http.Dir("templates"))
	http.Handle("/templates/", http.StripPrefix("/templates/", fs))

	// Charger les routes
	handler.SetupRoutes()

	// Démarrer le serveur sur le port 8080
	log.Println("Serveur démarré sur http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
