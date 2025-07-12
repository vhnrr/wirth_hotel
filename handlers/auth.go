package handlers

import (
	"database/sql"
	"net/http"
	"wirth_hotel/config"

	"github.com/gin-gonic/gin"
)

// Tampilkan halaman login
func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

// Tampilkan halaman register
func ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

// Handle register user baru (hash password)
func HandleRegister(c *gin.Context) {
	fullname := c.PostForm("fullname")
	email := c.PostForm("email")
	password := c.PostForm("password")
	dob := c.PostForm("dob")

	db := config.GetDB()

	var existing string
	err := db.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&existing)
	if err != sql.ErrNoRows {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Email sudah digunakan"})
		return
	}

	_, err = db.Exec("INSERT INTO users (fullname, email, password, dob) VALUES (?, ?, ?, ?)", fullname, email, password, dob)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Gagal mendaftar"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}

// Handle login dengan bcrypt cek password
func HandleLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	db := config.GetDB()

	var fullname, dbPassword string
	err := db.QueryRow("SELECT fullname, password FROM users WHERE email = ?", email).Scan(&fullname, &dbPassword)
	if err == sql.ErrNoRows {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Status": "failed",
			"error":  "Email tidak terdaftar",
		})
		return
	} else if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"Status": "failed",
			"error":  "Terjadi kesalahan server",
		})
		return
	}

	if dbPassword != password {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Status": "failed",
			"error":  "Password salah",
		})
		return
	}

	// Simpan ke session
	session, _ := config.Store.Get(c.Request, "session")
	session.Values["user"] = email
	session.Values["fullname"] = fullname
	session.Save(c.Request, c.Writer)

	c.Redirect(http.StatusSeeOther, "/")
}

// Logout user dan hapus session
func HandleLogout(c *gin.Context) {
	session, _ := config.GetSession(c)
	delete(session.Values, "user")
	session.Save(c.Request, c.Writer)
	c.Redirect(http.StatusSeeOther, "/login")
}

// Middleware untuk protect route agar hanya user login bisa akses
func RequireLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, _ := config.GetSession(c)
		_, ok := session.Values["user"]
		if !ok {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
