package handler

import (
	"html/template"
	"log"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/signup", signupHandler)
}

// Handler pour la page de connexion
func loginHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login.html")
}

// Handler pour la page d'inscription
func signupHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "signup.html")
}

// Fonction pour rendre les templates
func renderTemplate(w http.ResponseWriter, tmpl string) {
	tmplPath := "templates/" + tmpl
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Erreur de rendu du template %s: %v", tmpl, err)
		http.Error(w, "Erreur du serveur interne", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("Erreur d'exécution du template %s: %v", tmpl, err)
		http.Error(w, "Erreur du serveur interne", http.StatusInternalServerError)
	}
}
