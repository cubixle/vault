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
	"os"
	"strings"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
)

// Item holds the data.
type Item struct {
	Data   string    `json:"data"`
	Expiry time.Time `json:"expiryData"`
	TTL    int32     `json:"ttl"`
}

// Vault holds the vault data and key.
type Vault struct {
	Vault string `json:"vault"`
	Key   string `json:"key"`
}

func main() {
	router := gin.Default()
	router.Use(CORS(os.Getenv("VAULT_APP_URL")))
	router.POST("/", createAction)
	router.POST("/decrypt", decryptAction)
	router.Run(":7014")
}

// CORS handles setting up the CORS security headers.
func CORS(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func show(c *gin.Context) {
	log.Println("asdasdsa")
}

func createAction(c *gin.Context) {
	var item Item
	c.BindJSON(&item)

	key := generateUniqueID(16)
	json, _ := json.Marshal(&item)
	data := encrypt([]byte(key), string(json))

	var vault Vault
	vault.Key = key
	vault.Vault = data

	c.JSON(200, vault)
}

func decryptAction(c *gin.Context) {
	var vault Vault
	err := c.BindJSON(&vault)
	if err != nil {
		panic(err)
	}

	data := decrypt([]byte(vault.Key), vault.Vault)

	var item Item
	err = json.Unmarshal([]byte(data), &item)
	if err != nil {
		panic(err)
	}

	c.JSON(200, item)
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
