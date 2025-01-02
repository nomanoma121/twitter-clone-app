package handler

import (
	"database/sql"
	"errors"
	"log"
	"server/model"
	"sort"

	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type TodoHandler struct {
	db        *sqlx.DB
	validator *validator.Validate
}

func NewTodoHandler(db *sqlx.DB) *TodoHandler {
	return &TodoHandler{db: db, validator: validator.New()}
}

func (h *TodoHandler) Register(g *echo.Group) {
	g.GET("/todos", h.GetTodos)
	g.POST("/todos", h.CreateTodo)
	g.PUT("/todos/:id", h.UpdateTodo)
	g.DELETE("/todos/:id", h.DeleteTodo)
}

type GetTodosResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func (h *TodoHandler) GetTodos(c echo.Context) error {
	userID := c.Get("user_id").(int)

	var todos []model.Todo
	err := h.db.Select(&todos, "SELECT * FROM todos WHERE user_id = ?", userID)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(200, []GetTodosResponse{})
		}
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	sort.Slice(todos, func(i, j int) bool {
		return todos[i].CreatedAt.After(todos[j].CreatedAt)
	})

	res := make([]GetTodosResponse, len(todos))
	for i, todo := range todos {
		res[i] = GetTodosResponse{
			ID:        todo.ID,
			Title:     todo.Title,
			Completed: todo.Completed,
		}
	}

	return c.JSON(200, res)
}

type CreateTodoRequest struct {
	Title string `json:"title" validate:"required"`
}

type CreateTodoResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func (h *TodoHandler) CreateTodo(c echo.Context) error {
	userID := c.Get("user_id").(int)

	req := new(CreateTodoRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	res, err := h.db.Exec("INSERT INTO todos (user_id, title) VALUES (?, ?)", userID, req.Title)
	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	id, _ := res.LastInsertId()
	return c.JSON(201, CreateTodoResponse{
		ID:        int(id),
		Title:     req.Title,
		Completed: false,
	})
}

type UpdateTodoRequest struct {
	Id        int    `param:"id" validate:"required"`
	Title     string `json:"title" validate:"required"`
	Completed bool   `json:"completed"` // bool型の場合はrequiredを指定しない(ゼロ値==falseが入る)
}

func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	userID := c.Get("user_id").(int)

	req := new(UpdateTodoRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	if err := h.validator.Struct(req); err != nil {
		log.Println("validator error: ", err)
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	_, err := h.db.Exec("UPDATE todos SET title = ?, completed = ? WHERE id = ? AND user_id = ?", req.Title, req.Completed, req.Id, userID)
	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	return c.JSON(200, nil)
}

type DeleteTodoRequest struct {
	ID int `param:"id" validate:"required"`
}

func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	userID := c.Get("user_id").(int)

	req := new(DeleteTodoRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(400, map[string]string{"message": "Bad Request"})
	}

	_, err := h.db.Exec("DELETE FROM todos WHERE id = ? AND user_id = ?", req.ID, userID)
	if err != nil {
		return c.JSON(500, map[string]string{"message": "Internal Server Error"})
	}

	return c.JSON(200, nil)
}
