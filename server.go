package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/GeertJohan/go.rice"
	"strconv"
	"os"
	"io"
	"fmt"
)

func main() {
	e := echo.New()

	// the file server for rice. "app" is the folder where the files come from.
	publicAssetHandler := http.FileServer(rice.MustFindBox("public").HTTPBox())

	// serves the index.html from rice
	e.GET("/", echo.WrapHandler(publicAssetHandler))

	staticAssetHandler := http.FileServer(rice.MustFindBox("static").HTTPBox())

	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", staticAssetHandler)))
	e.POST("/upload", createProduct)
	e.GET("/file", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.QueryParam("id"))
		product := products[id]
		if product == nil {
			return c.JSON(http.StatusNotFound, false)
		}
		return c.File(product.FilePath)
	})

	e.GET("/list", func(c echo.Context) error {
		keys := make([]product, len(products))

		i := 0
		for _, v := range products {
			keys[i] = *v
			i++
		}
		return c.JSON(http.StatusOK, keys)
	})

	e.Logger.Fatal(e.Start(":8081"))
}

type (
	product struct {
		ID          int     `json:"id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		FileName    string  `json:"fileName"`
		WebUrl      string  `json:"url"`
		//ignore this field from JSON parsing
		FilePath string `json:"-" `
	}
)

var (
	products = map[int]*product{}
	seq      = 1
)

func createProduct(c echo.Context) error {

	homePath := os.Getenv("HOME")

	title := c.FormValue("title")
	description := c.FormValue("description")

	price, parseErr := strconv.ParseFloat(c.FormValue("price"), 32)
	if parseErr != nil {
		// do something sensible
	}
	url := c.FormValue("url")

	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	p := &product{
		ID: seq,
	}

	newFileName := fmt.Sprintf("file-%d", p.ID)
	fullFilePath := homePath + "/" + newFileName
	// Destination
	dst, err := os.Create(fullFilePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	p.Title = title
	p.Description = description
	p.Price = price
	p.WebUrl = url
	p.FileName = newFileName
	p.FilePath = fullFilePath
	if err := c.Bind(p); err != nil {
		return err
	}
	products[p.ID] = p
	seq++

	return c.JSON(http.StatusOK, p)

}
