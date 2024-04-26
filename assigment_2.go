package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	username = "postgres"
	password = "123123"
	hostname = "localhost"
	port     = 5432
	db       = "postgres"
)

type Task struct {
	ID        int64
	Name      string
	Completed bool
}

func createTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS task (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE
  )
  `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTask(db *sql.DB, task Task) error {
	var result sql.Result
	var err error
	query := "INSERT INTO task (name, completed) VALUES ($1, $2)"

	result, err = db.Exec(query, task.Name, task.Completed)
	if err != nil {
		return fmt.Errorf("CreateTask: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("CreateTask: %w", err)
	}

	fmt.Printf("Rows affected: %d\n", affected)
	return nil
}

func GetTask(db *sql.DB) ([]*Task, error) {
	var (
		id        int64
		name      string
		completed bool
	)

	result := make([]*Task, 0, 10)

	rows, err := db.Query("SELECT id, name, completed FROM task")
	if err != nil {
		return nil, fmt.Errorf("GetTask: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, &name, &completed); err != nil {
			return nil, fmt.Errorf("GetTask: %w", err)
		}
		result = append(result, &Task{
			ID:        id,
			Name:      name,
			Completed: completed,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetTask: %w", err)
	}
	return result, nil
}

func DeleteTasks(db *sql.DB, id int) error {
	query := "DELETE FROM task WHERE id=$1"
	_, err := db.Exec(query, id)
	return err
}

func updateTask(db *sql.DB, id int, completed bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE task SET completed = $1 WHERE id = $2", completed, id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	DSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, hostname, port, db)
	db, err := sql.Open("postgres", DSN)

	if err != nil {
		fmt.Println(err, "here i am")
		return
	}

	if err = db.Ping(); err != nil {
		fmt.Println(err, "here")
		return
	}

	defer db.Close()
	fmt.Println("Successfully connected to postgres")

	if err != nil {
		fmt.Println(err)
		return
	}

	createTable(db)

	addTask := Task{
		Name:      "Jump",
		Completed: false,
	}

	err = CreateTask(db, addTask)

	if err != nil {
		fmt.Println(err)
		return
	}

	tasks, err := GetTask(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, task := range tasks {
		fmt.Printf("\nID: %d. Name: %s. Completed: %s.", task.ID, task.Name, task.Completed)
	}

	updateTaskId := 6
	err = updateTask(db, updateTaskId, true)
	if err != nil {
		panic(err)
	}

	deleteTaskId := 1
	if err := DeleteTasks(db, deleteTaskId); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\nAll commands were successfully completed")
}
