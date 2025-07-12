package handlers

import (
	"net/http"
	"wirth_hotel/config"

	"github.com/gin-gonic/gin"
)

func ShowHomePage(c *gin.Context) {
	session, _ := config.GetSession(c)
	user, loggedIn := session.Values["user"]

	c.HTML(http.StatusOK, "home.html", gin.H{
		"LoggedIn": loggedIn,
		"User":     user,
	})
}

func ShowAboutPage(c *gin.Context) {
	session, _ := config.GetSession(c)
	user, loggedIn := session.Values["user"]
	c.HTML(http.StatusOK, "about.html", gin.H{
		"LoggedIn": loggedIn,
		"User":     user,
	})
}

func ShowReviewPage(c *gin.Context) {
	session, _ := config.GetSession(c)
	user, loggedIn := session.Values["user"]
	c.HTML(http.StatusOK, "review.html", gin.H{
		"LoggedIn": loggedIn,
		"User":     user,
	})
}

// Profile page menampilkan booking user
type Booking struct {
	BookingID  string
	Name       string
	Email      string
	CheckIn    string
	CheckOut   string
	GuestRoom  string
	GuestCount int
	TotalPrice int
	CreatedAt  string
}

func ShowProfilePage(c *gin.Context) {
	session, _ := config.Store.Get(c.Request, "session")
	fullname, _ := session.Values["fullname"].(string)
	email, ok := session.Values["user"].(string)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	db := config.GetDB()

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=?)", email).Scan(&exists)
	if err != nil || !exists {
		// User tidak ditemukan, hapus sesi dan redirect login
		delete(session.Values, "user")
		delete(session.Values, "fullname")
		session.Save(c.Request, c.Writer)
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	rows, err := db.Query(`
		SELECT booking_id, name, email, check_in, check_out, guest_room, guest_count, total_price, created_at
		FROM bookings
		WHERE email = ?
		ORDER BY created_at DESC
	`, email)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Gagal mengambil data booking", "UserName": fullname})
		return
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var b Booking
		err := rows.Scan(&b.BookingID, &b.Name, &b.Email, &b.CheckIn, &b.CheckOut, &b.GuestRoom, &b.GuestCount, &b.TotalPrice, &b.CreatedAt)
		if err == nil {
			bookings = append(bookings, b)
		}
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"Bookings": bookings,
		"UserName": fullname,
		"LoggedIn": true,
	})
}

// Handler di handlers.go
func DownloadReceipt(c *gin.Context) {
	bookingID := c.Param("bookingID")
	db := config.GetDB()

	var booking Booking
	err := db.QueryRow(`
        SELECT booking_id, name, email, check_in, check_out, guest_room, guest_count, total_price, created_at
        FROM bookings
        WHERE booking_id = ?`, bookingID).
		Scan(&booking.BookingID, &booking.Name, &booking.Email, &booking.CheckIn, &booking.CheckOut,
			&booking.GuestRoom, &booking.GuestCount, &booking.TotalPrice, &booking.CreatedAt)

	if err != nil {
		c.String(http.StatusNotFound, "Booking not found: "+err.Error())
		return
	}

	// Render halaman receipt
	c.HTML(http.StatusOK, "confirmation.html", gin.H{
		"BookingID":      booking.BookingID,
		"Name":           booking.Name,
		"Email":          booking.Email,
		"GuestRoom":      booking.GuestRoom,
		"GuestCount":     booking.GuestCount,
		"CheckIn":        booking.CheckIn,
		"CheckOut":       booking.CheckOut,
		"PriceFormatted": formatRupiah(booking.TotalPrice / booking.GuestCount), // asumsi
		"TotalNights":    1,
		"TotalFormatted": formatRupiah(booking.TotalPrice),
	})
}
