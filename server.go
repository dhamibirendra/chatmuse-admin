package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/GeertJohan/go.rice"
	"strconv"
	"os"
	"io"
	"fmt"
	"time"
	"math/rand"
)

func main() {

	e := echo.New()
	rand.Seed(time.Now().UnixNano())

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
		ImageUrl    string  `json:"imageUrl"`
		Payload     string  `json:"payload"`
		//ignore this field from JSON parsing
		FilePath string `json:"-" `
	}
)

var (
	products    = map[int]*product{}
	seq         = 1
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func createProduct(c echo.Context) error {

	homePath := os.Getenv("HOME")

	title := c.FormValue("title")
	description := c.FormValue("description")
	imageUrl := c.FormValue("imageUrl")

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
	p.ImageUrl = imageUrl
	p.FileName = newFileName
	p.FilePath = fullFilePath
	p.Payload = RandString()
	if err := c.Bind(p); err != nil {
		return err
	}
	products[p.ID] = p
	seq++

	return c.JSON(http.StatusOK, p)

}

func RandString() string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
