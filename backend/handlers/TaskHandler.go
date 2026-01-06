package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type (
	GetTaskResponse struct {
		ID       int64  `json:"id"`
		Desc     string `json:"description"`
		Deadline int64  `json:"deadline"`
	}

	CreateTaskRequest struct {
		Desc     string `json:"description"`
		Deadline int64  `json:"deadline"`
	}

	CreateTaskResponse struct {
		ID int64 `json:"id"`
	}

	UpdateTaskRequest struct {
		Desc     string `json:"description"`
		Deadline int64  `json:"deadline"`
	}

	Task struct {
		ID       int64
		Desc     string
		Deadline int64
	}
)

type TaskStorageInterface interface {
	CreateTask(task Task) int64
	GetTask(id int64) (Task, error)
	UpdateTask(task Task) (Task, error)
	DeleteTask(id int64) error
}

type TaskHandler struct {
	Storage TaskStorageInterface
}

type TaskStorage struct {
	mu            sync.RWMutex
	TaskIDCounter int64
	Tasks         map[int64]Task
}

func (t *TaskHandler) CreateTask(c *fiber.Ctx) error {

	var request CreateTaskRequest
	err := c.BodyParser(&request)
	if err != nil {
		return err
	}

	task := Task{
		ID:       t.Storage.TaskIDCounter,
		Desc:     request.Desc,
		Deadline: request.Deadline,
	}

	id := t.storage.CreateTask(task)
	t.storage.taskIDCounter++

	return c.Status(201).JSON(CreateTaskResponse{
		ID: id,
	})

}

func (t *TaskHandler) GetTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	task, err := t.storage.GetTask(idInt64)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.Status(200).JSON(GetTaskResponse{
		ID:       task.ID,
		Desc:     task.Desc,
		Deadline: task.Deadline,
	})

}

func (t *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).SendString("Id should exists")
	}
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(400).SendString("Error Json")
	}

	var request UpdateTaskRequest
	err = c.BodyParser(&request)
	if err != nil {
		return c.Status(400).SendString("Error Json")
	}

	task := Task{
		ID:       idInt64,
		Desc:     request.Desc,
		Deadline: request.Deadline,
	}

	_, err = t.storage.UpdateTask(task)
	if err != nil {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)

}

func (t *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).SendString("Id should exists")
	}
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(400).SendString("Error ID")
	}

	err = t.storage.DeleteTask(idInt64)
	if err != nil {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)

}

func (t *TaskStorage) CreateTask(task Task) int64 {
	t.tasks[task.ID] = task
	return task.ID

}

func (t *TaskStorage) GetTask(id int64) (Task, error) {

	result, ok := t.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	return result, nil

}

func (t *TaskStorage) UpdateTask(task Task) (Task, error) {
	t.tasks[task.ID] = task
	updated, ok := t.tasks[task.ID]
	if !ok {
		return Task{}, fmt.Errorf("Item not found")
	}

	return updated, nil

}

func (t *TaskStorage) DeleteTask(id int64) error {

	_, ok := t.tasks[id]
	if !ok {
		return fmt.Errorf("Item not found")
	}
	delete(t.tasks, id)

	return nil

}
