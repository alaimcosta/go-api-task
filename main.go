package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Desciption string    `json:"description"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type TaskService struct {
	DB          *sql.DB
	TaskChannel chan Task
}

// método, adicionar uma Tesk (tarefa) ###################
func (t *TaskService) AddTask(ts *Task) error {
	query := "INSERT INTO tasks (title, description, status, created_at) VALUES (?, ?, ?, ?)"
	result, err := t.DB.Exec(query, ts.Title, ts.Desciption, ts.Status, ts.CreatedAt)

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	ts.ID = int(id)
	return err
}

func (t *TaskService) UpdateTaskStatus(ts Task) error {
	query := "UPDATE tasks SET status = ? WHERE id = ?"
	_, err := t.DB.Exec(query, ts.Status, ts.ID)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// DeleteTask
func (t *TaskService) DeleteTask(ts string) error {
	query := "DELETE FROM tasks WHERE id = ?"
	_, err := t.DB.Exec(query, ts)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// DeleteTasks
func (t *TaskService) DeleteTasks() error {
	query := "DELETE FROM tasks"
	_, err := t.DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// ####################

func (t *TaskService) ListTasks() ([]Task, error) {
	rows, err := t.DB.Query("Select * from tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close() //fecha a conexão

	var tasks []Task
	for rows.Next() {
		var task Task
		//altera o valor na variável
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Desciption,
			&task.Status,
			&task.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}

// ProcessTasks -> processador de tarefas ################
func (t *TaskService) ProcessTasks() {
	for task := range t.TaskChannel {

		log.Printf("Processing task: %s", task.Title)
		time.Sleep(5 * time.Second)
		task.Status = "completed"

		t.UpdateTaskStatus(task)
		log.Printf("Task %s processed", task.Title)
	}
}

// Handler -> web services
func (t *TaskService) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//altero o status da tarefa
	task.Status = "pending"
	task.CreatedAt = time.Now()

	err = t.AddTask(&task)

	if err != nil {
		http.Error(w, "Error add task", http.StatusInternalServerError)
		return
	}

	//agora manda a tarefa para um canal, após adicioná-la

	t.TaskChannel <- task //mando a task para processamento através do canal

	w.WriteHeader(http.StatusCreated)
}

func (t *TaskService) HandleListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := t.ListTasks()
	if err != nil {
		http.Error(w, "Error listing tasks", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// HandleDeleteTask -> delete one task
func (t *TaskService) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := t.DeleteTask(id)
	if err != nil {
		http.Error(w, "Erro ao deletar a tarefa", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Println("Tarefa deletada com sucesso!")
}

// HandleDeleteTask -> delete one task
func (t *TaskService) HandleDeleteTasks(w http.ResponseWriter, r *http.Request) {
	err := t.DeleteTasks()
	if err != nil {
		http.Error(w, "Erro ao deletar as tarefas", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Println("Todas as tarefas foram deletadas com sucesso!")
}

func main() {
	db, err := sql.Open("sqlite3", "./data/tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	taskChannel := make(chan Task)

	taskService := TaskService{
		DB:          db,
		TaskChannel: taskChannel,
	}

	//funciona em background
	go taskService.ProcessTasks()

	//servidor web
	http.HandleFunc("POST /tasks", taskService.HandleCreateTask)
	http.HandleFunc("GET /tasks", taskService.HandleListTasks)
	http.HandleFunc("DELETE /task/{id}", taskService.HandleDeleteTask)
	http.HandleFunc("DELETE /tasks", taskService.HandleDeleteTasks)

	log.Println("Server running on port 8081")
	http.ListenAndServe(":8081", nil)
}
