package service

import (
	"net/http"
	model "simple-backend/internal/domain/todo"
	repo "simple-backend/internal/todo/repository/mysql"
	errorUtils "simple-backend/internal/utils/error"

	"gorm.io/gorm"
)

type todoService struct {
	Repository model.TodoRepoInterface
}

func Init(db *gorm.DB) model.TodoServiceInterface {
	return &todoService{
		Repository: repo.Init(db),
	}
}

func (t *todoService) GetAllTodo(field *model.TodoQueries) (*int64, []*model.TodoTable, *errorUtils.CustomError) {
	return t.Repository.GetAllTodo(field)
}

func (t *todoService) GetTodo(id uint, uid string) (*model.TodoTable, *errorUtils.CustomError) {
	todo := new(model.TodoTable)
	todo.Id = id
	todo.UserId = uid

	return t.Repository.GetTodo(todo)
}

func (t *todoService) CreateTodo(input *model.TodoCreate) *errorUtils.CustomError {
	table := new(model.TodoTable)
	table.UserId = input.UserId
	table.Title = input.Title
	table.Description = input.Description

	return t.Repository.CreateTodo(table)
}

func (t *todoService) UpdateTodo(input *model.TodoUpdate) (*model.TodoTable, *errorUtils.CustomError) {
	table := new(model.TodoTable)
	table.Id = input.Id
	table.UserId = input.UserId

	if input.Title != "" {
		table.Title = input.Title
	}

	if input.Description != "" {
		table.Description = input.Description
	}

	if input.IsCompleted != nil {
		table.IsCompleted = *input.IsCompleted
	}

	return t.Repository.UpdateTodo(table)
}

func (t *todoService) UpdateTodoCompleted(input *model.TodoUpdate) *errorUtils.CustomError {
	var err *errorUtils.CustomError
	updatedTodo := new(model.TodoTable)
	updatedTodo.Id = input.Id
	updatedTodo.UserId = input.UserId

	updatedTodo, err = t.Repository.GetTodo(updatedTodo)
	if err != nil {
		return err
	}

	if updatedTodo.UserId != input.UserId {
		return errorUtils.NewErrorResponse(
			http.StatusForbidden,
			"identify error",
			nil,
		)
	}

	if updatedTodo.IsCompleted == 1 {
		return errorUtils.NewErrorResponse(
			http.StatusBadRequest,
			"todo is completed",
			nil,
		)
	}

	return t.Repository.UpdateTodoCompleted(updatedTodo)
}

func (t *todoService) DeleteTodo(id uint, uid string) (*model.TodoTable, *errorUtils.CustomError) {
	var err *errorUtils.CustomError
	deletedTodo := new(model.TodoTable)
	deletedTodo.Id = id
	deletedTodo.UserId = uid

	deletedTodo, err = t.Repository.DeleteTodo(deletedTodo)

	return deletedTodo, err
}
