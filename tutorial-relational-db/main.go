package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type Character struct {
	// OK, apparently for models we need to have the first character uppercased.
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Level int    `json:"level"`
}

type Account struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	CharactersRemaining int    `json:"characters_remaining"`
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

	// Get accounts.
	accounts, err := getAccounts(db)
	if err != nil {
		panic(err)
	}

	fmt.Println(accounts)

	// Context reference: https://golangbot.com/connect-create-db-mysql/.
	ctx := context.Background()

	// Add character.
	themis := Character{
		Name:  "Themis",
		Role:  "Elidibus",
		Level: 99,
	}
	id, err := addCharacter(db, ctx, themis)
	if err != nil {
		panic(err)
	}

	fmt.Println(themis.Name, "successfully added with ID:", id)

	// Get accounts.
	accounts, err = getAccounts(db)
	if err != nil {
		panic(err)
	}

	fmt.Println(accounts)

	// Update.
	themis = Character{
		Name:  "Themis",
		Role:  "Elidibus",
		Level: 90,
	}
	_, err = updateCharacter(db, ctx, id, themis)
	if err != nil {
		panic(err)
	}

	fmt.Println(themis.Name, "successfully updated", themis.Name)

	// Get list of characters to ensure successfully added.
	characters, err = getCharacters(db)
	if err != nil {
		panic(err)
	}

	fmt.Println("Characters found: ", characters)

	// Delete character.
	_, err = deleteCharacter(db, ctx, id)
	if err != nil {
		panic(err)
	}

	fmt.Println(themis.Name, "successfully deleted")

	// Ensure that the character has been deleted.
	characters, err = getCharacters(db)
	if err != nil {
		panic(err)
	}

	fmt.Println("Characters found: ", characters)

	// Ensure that the characters_remaining field has been updated.
	// Get accounts.
	accounts, err = getAccounts(db)
	if err != nil {
		panic(err)
	}

	fmt.Println(accounts)
}

// Queries.
// First iteration was using db.Query etc.
// Second iteration was using prepared statements.
// Third iteration was using transactions for write-based actions. Also removes unnecessary prepared statements.
func getAccounts(db *sql.DB) ([]Account, error) {
	var accounts []Account

	rows, err := db.Query("SELECT * from accounts")
	if err != nil {
		return nil, fmt.Errorf("getAccounts: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		// While row still exists, iterate.
		var acc Account
		if err := rows.Scan(&acc.ID, &acc.Name, &acc.CharactersRemaining); err != nil {
			return nil, fmt.Errorf("getAccounts: %v", err)
		}

		accounts = append(accounts, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getAccounts: %v", err)
	}

	return accounts, nil
}

func getCharacters(db *sql.DB) ([]Character, error) {
	var characters []Character

	rows, err := db.Query("SELECT * from characters")
	if err != nil {
		return nil, fmt.Errorf("getCharacters: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		// While row still exists, iterate.
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

	rows, err := db.Query("SELECT * from characters WHERE role=?", role)
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

func getCharacterById(db *sql.DB, id int) (Character, error) {
	var char Character

	row := db.QueryRow("SELECT * from characters WHERE id=?", id)
	if err := row.Scan(&char.ID, &char.Name, &char.Role, &char.Level); err != nil {
		if err == sql.ErrNoRows {
			return char, fmt.Errorf("getCharacterById %q: no matching character", id)
		}

		return char, fmt.Errorf("getCharacterById %q: %v", id, err)
	}

	return char, nil
}

// Error helper function.
func fail(err error, label string) (int64, error) {
	return 0, fmt.Errorf("%s: %v", label, err)
}

func addCharacter(db *sql.DB, ctx context.Context, char Character) (int64, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err, "addCharacter")
	}
	defer tx.Rollback()

	// Check for remaining characters allowed.
	var account Account
	err = tx.
		QueryRowContext(ctx, "SELECT * from accounts WHERE name=?", "admin").
		Scan(&account.ID, &account.Name, &account.CharactersRemaining)
	if err != nil {
		if err != sql.ErrNoRows {
			return fail(err, "addCharacter")
		}
	}

	if account.CharactersRemaining == 0 {
		return fail(fmt.Errorf("No characters creation allowed for admin"), "addCharacter")
	}

	// Check for character name existence.
	var existingChar Character
	err = tx.
		QueryRowContext(ctx, "SELECT * from characters WHERE name=?", char.Name).
		Scan(&existingChar.ID, &existingChar.Name, &existingChar.Role, &existingChar.Level)
	if err != nil {
		if err != sql.ErrNoRows {
			return fail(err, "addCharacter")
		}
	}

	if existingChar.ID != 0 {
		// Character exists.
		return fail(fmt.Errorf("character %s exists", char.Name), "addCharacter")
	}

	// Insert.
	result, err := tx.ExecContext(ctx, "INSERT INTO characters (name, role, level) VALUES (?, ?, ?)", char.Name, char.Role, char.Level)
	if err != nil {
		return fail(err, "addCharacter")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fail(err, "addCharacter")
	}

	// Update number of characters allowed.
	_, err = tx.ExecContext(ctx, "UPDATE accounts SET characters_remaining=? WHERE name=?", account.CharactersRemaining-1, "admin")
	if err != nil {
		return fail(err, "addCharacter")
	}

	if err = tx.Commit(); err != nil {
		return fail(err, "addCharacter")
	}

	return id, nil
}

func updateCharacter(db *sql.DB, ctx context.Context, id int64, char Character) (int64, error) {
	stmt, err := db.PrepareContext(ctx, "UPDATE characters SET name=?, role=?, level=? WHERE id=?")
	if err != nil {
		return 0, fmt.Errorf("updateCharacter: %v", err)
	}

	result, err := stmt.ExecContext(ctx, char.Name, char.Role, char.Level, id)
	if err != nil {
		return 0, fmt.Errorf("updateCharacter: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("updateCharacter: %v", err)
	}

	return rowsAffected, nil
}

func deleteCharacter(db *sql.DB, ctx context.Context, id int64) (int64, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err, "deleteCharacter")
	}
	defer tx.Rollback()

	// Check for remaining characters allowed.
	var account Account
	err = tx.
		QueryRowContext(ctx, "SELECT * from accounts WHERE name=?", "admin").
		Scan(&account.ID, &account.Name, &account.CharactersRemaining)
	if err != nil {
		if err != sql.ErrNoRows {
			return fail(err, "deleteCharacter")
		}
	}

	// Delete.
	result, err := db.ExecContext(ctx, "DELETE FROM characters WHERE id=?", id)
	if err != nil {
		return 0, fmt.Errorf("deleteCharacter: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("deleteCharacter: %v", err)
	}

	// Update number of characters allowed.
	_, err = tx.ExecContext(ctx, "UPDATE accounts SET characters_remaining=? WHERE name=?", account.CharactersRemaining+1, "admin")
	if err != nil {
		return fail(err, "deleteCharacter")
	}

	if err = tx.Commit(); err != nil {
		return fail(err, "deleteCharacter")
	}

	return rowsAffected, nil
}
