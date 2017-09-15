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
	"gopkg.in/gomail.v2"
)

func main() {

	e := echo.New()

	runFixture()

	// the file server for rice. "app" is the folder where the files come from.
	publicAssetHandler := http.FileServer(rice.MustFindBox("public").HTTPBox())

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	e.Renderer = renderer
	rand.Seed(time.Now().UnixNano())

	// serves the index.html from rice
	e.GET("/", echo.WrapHandler(publicAssetHandler))
	e.GET("/settings", showSettingPage)

	staticAssetHandler := http.FileServer(rice.MustFindBox("static").HTTPBox())

	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", staticAssetHandler)))
	e.POST("/upload", createProduct)
	e.POST("/settings", updateSetting)
	e.POST("/booking", booking)
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

	bookingMsg struct {
		Message string `json:"message"`
	}
)

var (
	products    = map[int]*product{}
	seq         = 1
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	alertEmail  = "sajat.shrestha@gmail.com"
	fromEmail   = "chatmuse2018@gmail.com"
)

func updateSetting(c echo.Context) error {
	alertEmail = c.FormValue("email")
	return renderSettingPage(c, "Settings updated successfully")
}

func booking(c echo.Context) (err error) {
	b := new(bookingMsg)
	if err = c.Bind(b); err != nil {
		return
	}
	sendEmail("New booking done", b.Message)
	return c.JSON(http.StatusOK, b)
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

	p := &product{
		ID: seq,
	}

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
	products[p.ID] = p
	seq++

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
		"email":   alertEmail,
		"message": message,
	})
}

func sendEmail(subject string, message string) {
	m := gomail.NewMessage()

	m.SetHeader("From", fromEmail)
	m.SetHeader("To", alertEmail)

	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	d := gomail.NewDialer("smtp.gmail.com", 587, fromEmail, "sajat@123")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func runFixture() {
	p1 := &product{
		ID:          1,
		Title:       "Wildcraft Laptop Bags",
		Description: "Wildcraft Laptop Bags with Description",
		Price:       1419,
		WebUrl:      "https://smartdoko.com/shop/wildcraft-gravity-nylon-bag-8903338009566",
		ImageUrl:    "https://smartdoko.com/assets/upload/images/product/detail/IMG-0600cc412ddc0ff29dd2e437b89f9cc1.gif",
		Payload:     "flsymAvu",
	}

	p2 := &product{
		ID:          2,
		Title:       "Samsung Gear S2 Platinum",
		Description: "Step up your style with the Samsung Gear S2 classic. Genuine leather, precious metal and exceptional finishes come together to create a sophisticated wearable that goes with anything.",
		Price:       2040,
		WebUrl:      "https://smartdoko.com/shop/samsung-galaxy-s2-sports-gear-band",
		ImageUrl:    "https://smartdoko.com/assets/upload/images/product/detail/IMG-23d641578ee1dab84459058b8a02e1ee.jpg",
		Payload:     "krvixZGh",
	}

	p3 := &product{
		ID:          3,
		Title:       "Multipurpose Water Spray Gun ",
		Description: "Expandable & Flexible Plastic Water Garden Pipe with Spray Nozzle For Car Wash Pet Bath",
		Price:       1395,
		WebUrl:      "https://smartdoko.com/shop/multipurpose-water-spray-gun",
		ImageUrl:    "https://smartdoko.com/assets/upload/images/product/detail/IMG-332b397997496f4c2bcdf1d24ebb814a.gif",
		Payload:     "cQxAVrlD",
	}

	products[p1.ID] = p1
	products[p2.ID] = p2
	products[p3.ID] = p3

}
