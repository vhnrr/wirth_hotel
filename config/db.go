package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

var db *sql.DB
var Store = sessions.NewCookieStore([]byte("secret-key"))

// InitDB untuk menginisialisasi koneksi database
func InitDB() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/wirthhotel" // ganti sesuai setting DB kamu
	db, err = sql.Open("mysql", dsn)              // pakai = supaya assign ke package-level db
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Database cannot be accessed: ", err)
	}

	fmt.Println("Database connected successfully!")
}

// GetDB untuk mengakses variabel db dari luar package
func GetDB() *sql.DB {
	return db
}

// GetSession untuk mendapatkan session dari context
func GetSession(c *gin.Context) (*sessions.Session, error) {
	return Store.Get(c.Request, "session")
}
