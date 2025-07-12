package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wirth_hotel/config"

	"github.com/gin-gonic/gin"
)

type BookingDetails struct {
	Name       string
	Email      string
	CheckIn    string
	CheckOut   string
	GuestRoom  string
	GuestCount string
}

var roomPrices = map[string]int{
	"Old Money Room": 15000000,
	"Modern Room":    10000000,
	"Undersea Room":  25000000,
	"Viking Room":    25000000,
	"Mermaid Room":   30000000,
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateBookingID() string {
	return fmt.Sprintf("ORD-%d%d", time.Now().Unix(), rand.Intn(1000))
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func formatRupiah(n int) string {
	s := fmt.Sprintf("%d", n)
	var result []string
	reversed := reverse(s)

	for i, digit := range reversed {
		if i != 0 && i%3 == 0 {
			result = append(result, ".")
		}
		result = append(result, string(digit))
	}
	return reverse(strings.Join(result, ""))
}

// HandleBookingConfirmation simpan booking ke DB dan tampilkan konfirmasi
func HandleBookingConfirmation(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	checkIn := c.PostForm("check-in")
	checkOut := c.PostForm("check-out")
	guestRoom := c.PostForm("guest-room")
	guestCountStr := c.PostForm("guest-count")

	guestCount, err := strconv.Atoi(guestCountStr)
	if err != nil || guestCount < 1 {
		guestCount = 1 // default minimal 1 tamu
	}

	layout := "2006-01-02"
	tCheckIn, err1 := time.Parse(layout, checkIn)
	tCheckOut, err2 := time.Parse(layout, checkOut)
	totalNights := 1
	if err1 == nil && err2 == nil && tCheckOut.After(tCheckIn) {
		totalNights = int(tCheckOut.Sub(tCheckIn).Hours() / 24)
	}

	pricePerNight, ok := roomPrices[guestRoom]
	if !ok {
		pricePerNight = 0
	}

	totalPrice := pricePerNight * totalNights

	db := config.GetDB()
	bookingID := generateBookingID()

	// Gunakan time.Now() untuk created_at sesuai format MySQL
	createdAt := time.Now().Format("2006-01-02 15:04:05")

	_, err = db.Exec(`INSERT INTO bookings
		(booking_id, name, email, check_in, check_out, guest_room, guest_count, total_price, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		bookingID, name, email, checkIn, checkOut, guestRoom, guestCount, totalPrice, createdAt)

	if err != nil {
		c.String(http.StatusInternalServerError, "Gagal menyimpan data booking: "+err.Error())
		return
	}

	c.HTML(http.StatusOK, "confirmation.html", gin.H{
		"BookingID":      bookingID,
		"Name":           name,
		"Email":          email,
		"GuestRoom":      guestRoom,
		"GuestCount":     guestCount,
		"CheckIn":        checkIn,
		"CheckOut":       checkOut,
		"PriceFormatted": formatRupiah(pricePerNight),
		"TotalNights":    totalNights,
		"TotalFormatted": formatRupiah(totalPrice),
	})
}
