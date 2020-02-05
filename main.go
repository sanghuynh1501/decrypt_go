package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"syscall/js"

	"github.com/mergermarket/go-pkcs7"
)

type Time struct {
	WeekNumber    int
	UtcOffset     string
	UtcDatetime   string
	Unixtime      int
	Timezone      string
	RawOffset     int
	DstUntil      string
	DstOffset     int
	DstFrom       string
	Dst           bool
	DayOfYear     int
	DayOfWeek     int
	DateTime      string
	ClientIp      string
	Asbbreviation string
}

// Cipher key must be 32 chars long because block size is 16 bytes
const CIPHER_KEY = "abcdefghijklmnopqrstuvwxyz012345"

func CheckTime(unix_time int) bool {
	resp, err := http.Get("http://worldtimeapi.org/api/timezone/America/Argentina/Salta")
	if err != nil {
		log.Println("error to connect")
		return false
	}
	defer resp.Body.Close()
	var time Time
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	// this is where the magic happens, I pass a pointer of type Person and Go'll do the rest
	err = json.Unmarshal(body, &time)

	if err != nil {
		panic(err)
	}

	log.Println("time_stamp ", unix_time)
	log.Println("unix_time ", time.Unixtime)

	if unix_time > time.Unixtime {
		return false
	}

	return true
}

// Decrypt decrypts cipher text string into plain text string
func Decrypt(encrypted string) (string, error) {
	key := []byte(CIPHER_KEY)
	cipherText, _ := hex.DecodeString(encrypted)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(cipherText) < aes.BlockSize {
		panic("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	if len(cipherText)%aes.BlockSize != 0 {
		panic("cipherText is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText, _ = pkcs7.Unpad(cipherText, aes.BlockSize)
	return fmt.Sprintf("%s", cipherText), nil
}

func main() {
	time_stamp := js.Global().Get("time_stamp").Int()
	// hash_string := js.Global().Get("hash_string").String()
	checktime := CheckTime(time_stamp)
	log.Println("checktime ", checktime)
	if !checktime {
		log.Println("set plant_string")
		js.Global().Set("plant_string", "false")
	}
}
