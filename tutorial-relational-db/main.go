package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type Character struct {
	// OK, apparently for models we need to have the first character uppercased.
	ID    string `json:"id"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Level int    `json:"level"`
}

func main() {
	// Alternatively, we can just `source .env.sh` before running,
	// but I guess this can work, too.
	file, err := ioutil.ReadFile("./.env.sh")
	if err != nil {
		panic(err)
	}

	content := string(file)
	contentSlice := strings.Split(content, "\n")
	m := map[string]string{}

	for _, row := range contentSlice {
		kv := strings.Split(row, "=")
		l := len(kv)

		if l == 2 {
			m[kv[0]] = kv[1]
		}
	}

	cfg := mysql.Config{
		User:                 m["MYSQL_USERNAME"],
		Passwd:               m["MYSQL_PASSWORD"],
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "characters",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")

	characters, err := getCharactersByRole(db, "Scion")
	if err != nil {
		panic(err)
	}

	fmt.Println("Characters found: ", characters)

	// Get Uberdanger.
	var uberdanger Character
	for _, char := range characters {
		if char.Name == "Urianger" {
			uberdanger = char
			break
		}
	}

	character, err := getCharacterById(db, uberdanger.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println("Uberdanger found: ", character)

	themis := Character{
		Name:  "Themis",
		Role:  "Elidibus",
		Level: 99,
	}
	id, err := addCharacter(db, themis)
	if err != nil {
		panic(err)
	}

	fmt.Println(themis.Name, "successfully added with ID:", id)

	themis = Character{
		Name:  "Themis",
		Role:  "Elidibus",
		Level: 90,
	}
	_, err = updateCharacter(db, id, themis)
	if err != nil {
		panic(err)
	}

	fmt.Println(themis.Name, "successfully updated", themis.Name)

	characters, err = getCharacters(db)
	if err != nil {
		panic(err)
	}

	fmt.Println("Characters found: ", characters)

	_, err = deleteCharacter(db, id)
	if err != nil {
		panic(err)
	}

	fmt.Println(themis.Name, "successfully deleted")

	characters, err = getCharacters(db)
	if err != nil {
		panic(err)
	}

	fmt.Println("Characters found: ", characters)
}

// Queries.
// First iteration was using db.Query etc.
// Second iteration was using prepared statements.
func getCharacters(db *sql.DB) ([]Character, error) {
	var characters []Character

	stmt, err := db.Prepare("SELECT * from characters")
	if err != nil {
		return nil, fmt.Errorf("getCharacters: %v", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("getCharacters: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		// While row sitll exists, iterate.
		var char Character
		if err := rows.Scan(&char.ID, &char.Name, &char.Role, &char.Level); err != nil {
			return nil, fmt.Errorf("getCharacters: %v", err)
		}

		characters = append(characters, char)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getCharacters: %v", err)
	}

	return characters, nil
}

func getCharactersByRole(db *sql.DB, role string) ([]Character, error) {
	var characters []Character

	stmt, err := db.Prepare("SELECT * from characters WHERE role=?")
	if err != nil {
		return nil, fmt.Errorf("getCharactersByRole %q: %v", role, err)
	}

	rows, err := stmt.Query(role)
	if err != nil {
		return nil, fmt.Errorf("getCharactersByRole %q: %v", role, err)
	}
	defer rows.Close()

	for rows.Next() {
		// While row sitll exists, iterate.
		var char Character
		if err := rows.Scan(&char.ID, &char.Name, &char.Role, &char.Level); err != nil {
			return nil, fmt.Errorf("getCharactersByRole %q: %v", role, err)
		}

		characters = append(characters, char)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getCharactersByRole %q: %v", role, err)
	}

	return characters, nil
}

func getCharacterById(db *sql.DB, id string) (Character, error) {
	var char Character

	stmt, err := db.Prepare("SELECT * from characters WHERE id=?")
	if err != nil {
		return char, fmt.Errorf("getCharacterById: %v", err)
	}

	row := stmt.QueryRow(id)
	if err := row.Scan(&char.ID, &char.Name, &char.Role, &char.Level); err != nil {
		if err == sql.ErrNoRows {
			return char, fmt.Errorf("getCharacterById %q: no matching character", id)
		}

		return char, fmt.Errorf("getCharacterById %q: %v", id, err)
	}

	return char, nil
}

func addCharacter(db *sql.DB, char Character) (int64, error) {
	stmt, err := db.Prepare("INSERT INTO characters (name, role, level) VALUES (?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("addCharacter: %v", err)
	}

	result, err := stmt.Exec(char.Name, char.Role, char.Level)
	if err != nil {
		return 0, fmt.Errorf("addCharacter: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addCharacter: %v", err)
	}

	return id, nil
}

func updateCharacter(db *sql.DB, id int64, char Character) (int64, error) {
	stmt, err := db.Prepare("UPDATE characters SET name=?, role=?, level=? WHERE id=?")
	if err != nil {
		return 0, fmt.Errorf("updateCharacter: %v", err)
	}

	result, err := stmt.Exec(char.Name, char.Role, char.Level, id)
	if err != nil {
		return 0, fmt.Errorf("updateCharacter: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("updateCharacter: %v", err)
	}

	return rowsAffected, nil
}

func deleteCharacter(db *sql.DB, id int64) (int64, error) {
	stmt, err := db.Prepare("DELETE FROM characters WHERE id=?")
	if err != nil {
		return 0, fmt.Errorf("deleteCharacter: %v", err)
	}

	result, err := stmt.Exec(id)
	if err != nil {
		return 0, fmt.Errorf("deleteCharacter: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("deleteCharacter: %v", err)
	}

	return rowsAffected, nil
}
