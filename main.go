package main

import "Authentification/database"

func main() {

	database.Opendb()

	defer database.CloseDB()

}
