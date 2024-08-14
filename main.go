package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

type Message struct {
	MessageID int    `json:"messageID"`
	Message   string `json:"message"`
	UserID    string `json:"userID"`
}

func main() {
	router := httprouter.New()
	router.POST("/message/:id", sendMessage)
	router.GET("/messages/:id", getMessages)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func initDb() (*sql.DB, error) {
	pgURL := os.Getenv("POSTGRES_URL")

	db, err := sql.Open("postgres", pgURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func sendMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	db, err := initDb()
	if err != nil {
		writeInternalErr(w, err)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeInternalErr(w, err)
		return
	}
	defer r.Body.Close()

	query := fmt.Sprintf("INSERT INTO messages (message, user_id) VALUES ('%s', %s)", string(bodyBytes), params.ByName("id"))

	if _, err := db.Exec(query); err != nil {
		writeInternalErr(w, err)
		return
	}
}

func getMessages(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	db, err := initDb()
	if err != nil {
		writeInternalErr(w, err)
		return
	}

	query := fmt.Sprintf("SELECT * FROM messages WHERE id = %s", params.ByName("id"))

	rows, err := db.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "no messages found for user "+params.ByName("id"), http.StatusNotFound)
			return
		} else {
			writeInternalErr(w, err)
			return
		}
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.MessageID, &message.UserID, &message.Message); err != nil {
			writeInternalErr(w, err)
			return
		}

		messages = append(messages, &message)
	}

	response, err := json.Marshal(messages)
	if err != nil {
		writeInternalErr(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func writeInternalErr(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
