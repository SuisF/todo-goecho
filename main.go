package main

import (
	"todo-goecho/dto"
	
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	tm := dto.NewTodoManager() //instance of todo manager

	e := echo.New()

	e.Use(middleware.Logger())//log all requests

	e.GET("/", func(c echo.Context) error { //sending json data back
		todos := tm.GetAll()
		return c.JSON(http.StatusOK, todos)
	})

	authenticatedGroup := e.Group("/todos", func(next echo.HandlerFunc) echo.HandlerFunc{
		return func(c echo.Context) error {
			authorization := c.Request().Header.Get("authorization")//check for auth-token
			if authorization != "auth-token" {
				c.Error(echo.ErrUnauthorized)
				return nil
			}

			next(c)
			return nil
		}
	})

	authenticatedGroup.POST("/create", func(c echo.Context) error {
		requestBody := dto.CreateTodoRequest{}
		
		//parses the body and binds it to the passed type
		err := c.Bind(&requestBody)
		if err != nil {
			return err
		}

		todo := tm.Create(requestBody)
		return c.JSON(http.StatusCreated, todo)
	})

	authenticatedGroup.PATCH("/:id/complete", func(c echo.Context) error {
		id := c.Param("id")

		err := tm.Complete(id)
		if err != nil {
			c.Error(err)
			return err
		}

		return nil	
	})

	authenticatedGroup.DELETE("/:id", func(c echo.Context) error {
		id := c.Param("id")

		err := tm.Remove(id)
		if err != nil {
			c.Error(err)
			return err
		}

		return nil
	})
	
	e.Start(":2500")


}