package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Item holds the data.
type Item struct {
	Data   string    `json:"data"`
	Expiry time.Time `json:"expiryDate"`
	TTL    int       `json:"ttl"`
}

// Vault holds the vault data and key.
type Vault struct {
	Vault string `json:"vault" binding:"required"`
	Key   string `json:"key" binding:"required"`
}

func main() {
	appURL := os.Getenv("VAULT_APP_URL")
	if appURL == "" {
		appURL = "*"
	}

	port := os.Getenv("VAULT_PORT")
	if port == "" {
		port = "8080"
	}

	router := echo.New()

	router.POST("/", createAction)
	router.POST("/decrypt", decryptAction)

	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{appURL},
	}))

	if err := router.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}

func createAction(c echo.Context) error {
	item := &Item{}
	if err := c.Bind(item); err != nil {
		return c.JSON(http.StatusBadRequest, `{"error": "invalid form data"}`)
	}

	if item.Data == "" {
		return c.JSON(http.StatusBadRequest, `{"error": "missing data"}`)
	}

	if item.TTL == 0 {
		return c.JSON(http.StatusBadRequest, `{"error": "missing ttl"}`)
	}

	currentTime := time.Now()
	item.Expiry = currentTime.Add(time.Duration(item.TTL) * time.Second)

	key := generateUniqueID(16)
	json, err := json.Marshal(&item)
	if err != nil {
		return c.JSON(http.StatusBadRequest, `{"error": "invalid form data"}`)
	}

	data := encrypt([]byte(key), string(json))

	var vault Vault
	vault.Key = key
	vault.Vault = data

	return c.JSON(http.StatusOK, vault)
}

func decryptAction(c echo.Context) error {
	var vault Vault
	err := c.Bind(&vault)
	if err != nil {
		return c.JSON(http.StatusBadRequest, `{"error": "invalid form data"}`)
	}

	if vault.Key == "" {
		return c.JSON(http.StatusBadRequest, `{"error": "missing key"}`)
	}

	if vault.Vault == "" {
		return c.JSON(http.StatusBadRequest, `{"error": "missing vault data"}`)
	}
	data := decrypt([]byte(vault.Key), vault.Vault)

	var item Item
	err = json.Unmarshal([]byte(data), &item)
	if err != nil {
		return c.JSON(http.StatusBadRequest, `{"error": "invalid form data"}`)
	}

	currentTime := time.Now()
	if currentTime.Unix() > item.Expiry.Unix() {
		return c.JSON(http.StatusBadRequest, `{"error": "invalid time"}`)
	}

	return c.JSON(http.StatusOK, item)
}

func generateUniqueID(length int) string {
	n := length
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)

	return strings.ToLower(s)
}

func encrypt(key []byte, message string) string {
	plainText := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return ""
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//returns to base64 encoded string
	encmess := base64.URLEncoding.EncodeToString(cipherText)
	return encmess
}

func decrypt(key []byte, securemess string) string {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return ""
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}

	if len(cipherText) < aes.BlockSize {
		return ""
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	decodedmess := string(cipherText)
	return decodedmess
}
