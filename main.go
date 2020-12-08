package main

import (
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/sanditya12/learngo/scrapper"
)

const fileName string = "jobs.csv"

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.File("./index.html")
	})

	e.POST("/scrapper", func(c echo.Context) error {
		defer os.Remove(fileName)
		key := strings.ToLower(scrapper.CleanString(c.FormValue("key")))
		scrapper.Scrap(key)
		return c.Attachment(fileName, fileName)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
