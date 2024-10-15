package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Opendb() {

	var err error

	DB, err = sql.Open("sqlite3", "./auth.db")

	if err != nil {
		fmt.Printf("erreur lors de l'ouverrture de la base de donnée: %v", err)
		return
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(30) NOT NULL UNIQUE,
		name VARCHAR(30) NOT NULL,
		password VARCHAR(15) NOT NULL
	);`
	_, err = DB.Exec(createTableSQL) // Exécuter la requête de création de table
	if err != nil {
		log.Fatal(err) // Gérer l'erreur si la création échoue
		return
	}

}

func CloseDB() {
	err := DB.Close()
	if err != nil {
		log.Fatal("Erreur lors de la fermeture de la base de données:", err)
	}
	fmt.Println("Base de données fermée avec succès.")
}
