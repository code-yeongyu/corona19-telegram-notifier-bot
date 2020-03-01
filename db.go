package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var dbInfo = os.Getenv("DB_INFO")

// GetRecentNumbers returns the recent number of records
func GetRecentNumbers() (numbers map[string]int) {
	db, err := sql.Open("mysql", dbInfo)
	if err != nil {
		return nil
	}
	defer db.Close()

	data := db.QueryRow("SELECT * FROM numbers  ORDER BY id DESC LIMIT 1")
	var id, confirmed, death, cured int
	data.Scan(&id, &confirmed, &death, &cured)
	numbers = make(map[string]int)
	numbers["confirmed"] = confirmed
	numbers["death"] = death
	numbers["cured"] = cured
	return
}

// AddNumbers adds the record of number
// the keys of the numbers map[string]int should be:
// 'confirmed': numbers of people who got confirmed
// 'death': numbers of people who died because of the corona19
// 'cured' : numbers of people who cured from the corona19
func AddNumbers(numbers map[string]int) {
	db, err := sql.Open("mysql", dbInfo)

	if err != nil {
		return
	}
	defer db.Close()
	query := fmt.Sprintf("INSERT INTO numbers(confirmed, death, cured) VALUES (%d, %d, %d)", numbers["confirmed"], numbers["death"], numbers["cured"])
	db.Exec(query)
}

// GetChatIDs returns the chat IDs of all the users
func GetChatIDs() (IDs []int64) {
	db, err := sql.Open("mysql", dbInfo)
	if err != nil {
		return nil
	}
	defer db.Close()

	rows, err := db.Query("SELECT chat_id FROM chatIDs")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var chatID int
		err = rows.Scan(&chatID)
		if err != nil {
			panic(err)
		}
		IDs = append(IDs, int64(chatID))
	}

	return
}

// AddChatID adds a chat ID to the database
func AddChatID(id int64) {
	db, err := sql.Open("mysql", dbInfo)
	if err != nil {
		return
	}
	defer db.Close()

	query := fmt.Sprintf("INSERT INTO chatIDs(chat_id) VALUES (%d)", id)
	db.Exec(query)
}
