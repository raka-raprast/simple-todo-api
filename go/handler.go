package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Todo struct {
    ID            int          `json:"id"`
    Title         string       `json:"title"`
    Description   string       `json:"description"`
    AddedDate     time.Time    `json:"added_date"`
    CompletedDate sql.NullTime `json:"completed_date"`
}

type UpdateTodoRequest struct {
	Title         string       `json:"title"`
	Description   string       `json:"description"`
	CompletedDate sql.NullTime `json:"completed_date"`
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM todos")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.AddedDate, &todo.CompletedDate)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}


func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	json.NewDecoder(r.Body).Decode(&todo)

	err := db.QueryRow("INSERT INTO todos (title, description, added_date) VALUES ($1, $2, $3) RETURNING id", todo.Title, todo.Description, time.Now()).Scan(&todo.ID)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}	

func GetTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var todo Todo
	err := db.QueryRow("SELECT * FROM todos WHERE id = $1", id).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.AddedDate, &todo.CompletedDate)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	rawJSON, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	log.Println("Raw JSON Payload:", string(rawJSON))

	var updatedTodo UpdateTodoRequest
	err = json.Unmarshal(rawJSON, &updatedTodo)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Updating TODO with ID %d. New data: %+v", id, updatedTodo)

    _, err = db.Exec("UPDATE todos SET title = $1, description = $2, completed_date = $3 WHERE id = $4",
        updatedTodo.Title, updatedTodo.Description, updatedTodo.CompletedDate.Time, id)
    if err != nil {
        log.Fatal(err)
    }

    var todo Todo
    err = db.QueryRow("SELECT * FROM todos WHERE id = $1", id).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.AddedDate, &todo.CompletedDate)
    if err != nil {
        log.Fatal(err)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusNoContent)
}
