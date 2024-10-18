package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

var homework = make(map[string]bool)

var User struct {
	name     string
	passw    string
	homework map[string]bool
}

type InfoHomework struct {
	Info string
}

var users = map[string]string{
	"petya":  "1111",
	"masha":  "0101",
	"karina": "2222",
}

func getHomework(c echo.Context) error {
	hw := c.FormValue("hw")
	homework[hw] = false

	return c.String(http.StatusOK, "subject accepted, "+hw+" is waiting for you")
}

func showHomework(c echo.Context) error {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		return err
	}

	var uncompletedHw, completedHw string

	for key, compl := range homework {
		if !compl {
			uncompletedHw += fmt.Sprintf("subject: %s, isn't completed\n", key)
		} else {
			completedHw += fmt.Sprintf("subject: %s, is completed\n", key)
		}
	}

	info := InfoHomework{
		Info: fmt.Sprintf("Uncompleted homework:\n%s\nCompleted homework:\n%s", uncompletedHw, completedHw),
	}
	return tmpl.Execute(c.Response(), info)
}

func updateHomework(c echo.Context) error {
	subject := c.QueryParam("hw")

	if _, exists := homework[subject]; !exists {
		return c.String(http.StatusBadRequest, "subject doesn't exist")
	} else {
		homework[subject] = true
		return c.String(http.StatusOK, "congrats! "+subject+" is completed!")
	}
}

func logIn(c echo.Context) error {
	name := c.QueryParam("name")

	passw := c.QueryParam("passw")

	if _, exists := users[name]; !exists || users[name] != passw {
		return c.String(http.StatusBadRequest, "invalid passw or name")
	}

	return c.String(http.StatusOK, fmt.Sprintf("you logged in as %s", name))
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.File("index.html")
	})

	e.POST("/submit", getHomework)

	e.GET("/hw-show", showHomework)

	e.GET("/hw-compl", updateHomework)

	e.GET("/log-in", logIn)

	e.Logger.Fatal(e.Start(":1323"))
}
