package dto

import (
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

//Todo represents the todo
type Todo struct {
	ID			string `json:"id"`
	Title		string `json:"title"`
	IsComplete	bool `json:"isDone"`
}

type TodoManager struct {
	todos	[]Todo
	m		sync.Mutex //for atomic operations
}

func NewTodoManager() TodoManager {
	return TodoManager{
		todos: 	make([]Todo, 0),
		m: 		sync.Mutex{},
	}
}

func (tm *TodoManager) GetAll() []Todo {
	return tm.todos
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

func (tm *TodoManager) Create( createTodoRequest CreateTodoRequest ) Todo {
	tm.m.Lock()
	defer tm.m.Unlock()

	newTodo := Todo {
		ID:				strconv.FormatInt(time.Now().UnixMilli(), 10),
		Title:			createTodoRequest.Title,
		IsComplete:		false,	
	}

	tm.todos = append(tm.todos, newTodo)
	return newTodo
}

func (tm *TodoManager) Complete(ID string) error {
	tm.m.Lock()
	defer tm.m.Unlock()

	//find todo through id
	var todo *Todo
	var index int = -1

	for i,t := range tm.todos {
		if t.ID == ID {
			todo = &t
			index = i
		}
	}

	if todo == nil {
		return echo.ErrNotFound
	}

	if todo.IsComplete {
		err := echo.ErrBadRequest
		err.Message = "todo is already complete"
		return err
	}

	//update todo
	tm.todos[index].IsComplete = true
	return nil
}

func (tm *TodoManager) Remove(ID string) error {
	tm.m.Lock()
	defer tm.m.Unlock()

	index := -1

	for i,t := range tm.todos {
		if t.ID == ID {
			index = i
			break
		}
	}

	if index == -1 {
		return echo.ErrNotFound
	}

	tm.todos = append(tm.todos[:index], tm.todos[index+1:]...)
	return nil
}