package handler

import (
	"Authentification/database"
	"Authentification/model"
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Fonction pour définir les routes
func SetupRoutes() {
	http.HandleFunc("/", acceuilRoute)
	http.HandleFunc("/signup", signupRoute)
	http.HandleFunc("/login", loginRoute)
	http.HandleFunc("/dashboard", dashboardRoute)
}

// Handler pour le formulaire d'inscription
func signupRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Récupérer les valeurs du formulaire
		nom := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Hasher le mot de passe
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Erreur lors du hachage du mot de passe", http.StatusInternalServerError)
			return
		}

		// Créer un utilisateur
		user := model.User{
			Nom:      nom,
			Email:    email,
			Password: string(hashedPassword), // Enregistrer le mot de passe haché
		}

		// Insérer l'utilisateur dans la base de données
		err = insertUser(user)
		if err != nil {
			http.Error(w, "Erreur lors de l'insertion dans la base de données", http.StatusInternalServerError)
			return
		}

		// Rediriger vers la page de connexion après l'inscription
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		renderTemplate(w, "signup.html")
	}
}

func loginRoute(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode est POST ou GET
	if r.Method == http.MethodGet {
		renderTemplate(w, "login.html")
		return
	}

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Extraire les données du formulaire
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Vérifier que les champs ne sont pas vides
	if email == "" || password == "" {
		http.Error(w, "Email et mot de passe requis", http.StatusBadRequest)
		return
	}

	// Récupérer l'utilisateur dans la base de données en utilisant l'email
	var user model.User
	query := "SELECT id, email, password FROM users WHERE email = ?"
	err := database.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Email ou mot de passe incorrect", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Erreur serveur, veuillez réessayer", http.StatusInternalServerError)
		return
	}

	// Comparer le mot de passe haché avec le mot de passe fourni
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Email ou mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	// Si tout est correct, rediriger l'utilisateur vers la page de succès ou dashboard
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Fonction pour insérer un utilisateur dans la base de données
func insertUser(user model.User) error {
	db, err := sql.Open("sqlite3", "auth.db") // Assure-toi que le nom est correct
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO users (nom, email, password) VALUES (?, ?, ?)`
	_, err = db.Exec(query, user.Nom, user.Email, user.Password)
	if err != nil {
		// Journaliser l'erreur pour un diagnostic plus facile
		log.Printf("Erreur lors de l'insertion de l'utilisateur dans la base de données : %v", err)
		return err
	}
	return nil
}

// Fonction pour afficher les templates
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

func acceuilRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		renderTemplate(w, "acceuil.html")
	}
}

func dashboardRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		renderTemplate(w, "dashboard.html")
	}
}
