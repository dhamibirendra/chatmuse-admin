package main

import (
	"net/http"
	"html/template"
	"github.com/labstack/echo"
	"github.com/GeertJohan/go.rice"
	"strconv"
	"os"
	"io"
	"fmt"
	"time"
	"math/rand"
	"./src"
)

type (
	bookingMsg struct {
		Message string `json:"message"`
	}
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func main() {

	e := echo.New()
	rand.Seed(time.Now().UnixNano())
	src.RunFixture()

	// the file server for rice. "app" is the folder where the files come from.
	publicAssetHandler := http.FileServer(rice.MustFindBox("public").HTTPBox())

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	e.Renderer = renderer

	// serves the index.html from rice
	e.GET("/", echo.WrapHandler(publicAssetHandler))
	staticAssetHandler := http.FileServer(rice.MustFindBox("static").HTTPBox())
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", staticAssetHandler)))
	e.POST("/upload", createProduct)
	e.POST("/settings", updateSetting)
	e.POST("/booking", doBooking)
	e.GET("/file", fetchFile)
	e.GET("/list", showProductsList)
	e.GET("/settings", showSettingPage)

	e.Logger.Fatal(e.Start(":8081"))
}

func updateSetting(c echo.Context) error {
	src.AlertEmail = c.FormValue("email")
	return renderSettingPage(c, "Settings updated successfully")
}

func doBooking(c echo.Context) (err error) {
	b := new(bookingMsg)
	if err = c.Bind(b); err != nil {
		return
	}
	src.SendEmail("New booking done", b.Message)
	return c.JSON(http.StatusOK, b)
}

func fetchFile(c echo.Context) error {
	id, _ := strconv.Atoi(c.QueryParam("id"))
	product := src.Products[id]
	if product == nil {
		return c.JSON(http.StatusNotFound, false)
	}
	return c.File(product.FilePath)
}

func showProductsList(c echo.Context) error {
	return c.JSON(http.StatusOK, src.GetProductsList())
}

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

	p := src.NewProduct()

	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("file")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

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

		p.FileName = newFileName
		p.FilePath = fullFilePath
	}

	p.Title = title
	p.Description = description
	p.Price = price
	p.WebUrl = url
	p.ImageUrl = imageUrl
	p.Payload = RandString()
	if err := c.Bind(p); err != nil {
		return err
	}
	src.Products[p.ID] = p
	src.ProductSeq++

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"message": "Product added successfully",
	})

}

func RandString() string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func showSettingPage(c echo.Context) error {
	return renderSettingPage(c, "")
}

func renderSettingPage(c echo.Context, message string) error {
	return c.Render(http.StatusOK, "settings.html", map[string]interface{}{
		"email":   src.AlertEmail,
		"message": message,
	})
}
