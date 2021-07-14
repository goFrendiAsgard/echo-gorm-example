package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	connectionString := os.Getenv("MYSQL_CONNECTION_STRING")
	if connectionString == "" {
		connectionString = "root:toor@tcp(localhost:3306)/alta"
	}
	db, err := initDb(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	e := echo.New()
	e.GET("/", hello)
	bookController := NewBookController(db)
	e.GET("/books", bookController.GetBooks)
	e.POST("/books", bookController.CreateBook)
	e.Start(":8080")
}

func initDb(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(connectionString))
	if err != nil {
		return db, err
	}
	err = db.AutoMigrate(&Book{})
	return db, err
}

func hello(c echo.Context) error {
	c.String(http.StatusOK, "hello world")
	return nil
}

type Book struct {
	gorm.Model
	Title  string `json:"title,omitempty" form:"title"`
	Author string `json:"author,omitempty" form:"author"`
}

type BookController struct {
	db *gorm.DB
}

func NewBookController(db *gorm.DB) *BookController {
	return &BookController{
		db: db,
	}
}

func (bc *BookController) GetBooks(c echo.Context) error {
	books := []Book{}
	if tx := bc.db.Find(&books); tx.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, tx.Error)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success getting books",
		"result":  books,
	})
}

func (bc *BookController) CreateBook(c echo.Context) error {
	book := Book{}
	c.Bind(&book)
	if tx := bc.db.Save(&book); tx.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, tx.Error)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success inserting boook",
		"result":  book,
	})
}
