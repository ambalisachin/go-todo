package controllers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//IV         = "1461618689689168"(x-iv)
//passphrase = "noenonrgkgneroiw"(x-key)

//EncryptDataHandler function is a handler for a web request that encrypts data
func EncryptDataHandler(ctx *gin.Context) {
	//declares a variable called "requestBody" with type "interface{}". An interface{} is a type that is used any type of data.
	// This is useful as it allows the variable to store any type of data without needing to specify a specific type.
	var requestBody interface{}
	//The ShouldBindJSON method takes the request body from the "context" and then decodes it into the "requestBody" variable.
	// The result of this operation is that the "requestBody" variable now contains the data from the JSON request body.
	ctx.ShouldBindJSON(&requestBody)
	//marshals the request body into a plain text format. The json.Marshal function takes an interface{} type and returns a byte slice that is the plain text representation of the request body.
	//If there is an error, it is stored in the err variable.
	plainText, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(requestBody)
	//encrypt a plaintext using the x-key & x-iv values retrieved from the request header.AESEncrypt func takes the plaintext, x-key (byte array) and x-iv ( string) as arguments & returns the encrypted data.
	encryptedData := AESEncrypt(string(plainText), []byte(ctx.Request.Header.Get("x-key")), ctx.Request.Header.Get("x-iv"))
	//encode a string of data that has been encrypted into a Base64 string.base64.StdEncoding.EncodeToString() func takes an encryptedData parameter(byte array)&returns a string that is the encoded version of the data.
	//The resulting encoded string is stored in the variable encryptedString.
	encryptedString := base64.StdEncoding.EncodeToString(encryptedData)
	//return an encrypted string as a JSON response.1st parameter is the HTTP status code(200 (OK))&2nd parameter is the encrypted string, which is returned in the response.
	ctx.JSON(http.StatusOK, encryptedString)
}
//DecryptDataHandler function is a handler for a web request that decrypts data
func DecryptDataHandler(ctx *gin.Context) {
	//creates a variable called requestBody that is a map[string]string.the requestBody variable is a map  that has string keys and string values.
	//This can be used to store request body data from an HTTP request.
	var requestBody map[string]string
	//Refer line no 23 & 24
	ctx.ShouldBindJSON(&requestBody)
	fmt.Println(requestBody, requestBody["Encrypted-data"])
	//	decodes a string that was previously encoded in Base64. The string is stored in the requestBody variable as the value for the key "Encrypted-data".
	// The decoded data is stored in the variable encryptedData.  (_) is used to ignore the second return value from the DecodeString call.
	encryptedData, _ := base64.StdEncoding.DecodeString(requestBody["Encrypted-data"])
	//decrypting an encrypted data with an AES Decrypt function.passing the encrypted data, a byte array ("x-key" ) & ("x-iv") header in the request.
	// The function then returns the decrypted text.
	decryptedText := AESDecrypt(encryptedData, []byte(ctx.Request.Header.Get("x-key")), ctx.Request.Header.Get("x-iv"))
	fmt.Println("\n decrypted data:", string(decryptedText))
	//send an HTTP response to with a status code of 200 (OK) & decrypted text as the response body.
	//The ctx is a context object which is used to send the response, the http.StatusOK is used to set the status code of the response and the string(decryptedText) is used to set the response body.
	ctx.JSON(http.StatusOK, string(decryptedText))
}
//AESEncrypt used to encrypt a given source of bytes using the AES algorithm. 
func AESEncrypt(src string, key []byte, IV string) []byte {
	// Create the AES cipher block key
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
	}
	if src == "" {
		fmt.Println("plain content empty")
	}
	// create a new CBC encryption from block and IV
	ecb := cipher.NewCBCEncrypter(block, []byte(IV))
	//Convert string variable called src into a slice of bytes
	//the slice if bytes is then stored in a variable called content
	content := []byte(src)
	//padding used to pad palintext before encryption
	//content is being padded using block size of block cipher before encryption
	//this ensure that the paintext has uniform length  and encryted pro
	content = PKCS5Padding(content, block.BlockSize())
	//create a byte slice called crypted with a size equal to length of content agrument
	//the byte slice is initializes to contain all zero values
	crypted := make([]byte, len(content))
	//ECB mode of operation to encrypt the given content using the provided key.
	// The "crypted" variable is used to store the encrypted output,
	//while content variable contains the original plaintext data to be encrypted.
	//The CryptBlocks() performs  actual encryption, using the ECB mode.
	ecb.CryptBlocks(crypted, content)
	return (crypted)

}

// PKCS5Padding process is to add extra bytes to the end of the data
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	//Calculate the required padding size for a given ciphertext.
	//The padding size is calculated by subtracting the length of the ciphertext from the blocksize,
	//then taking the modulus of the result.
	padding := blockSize - len(ciphertext)%blockSize
	//This code in will create a byte slice containing a repeated value of the byte given in the variable padding.
	//For example, if padding is set to 5, the byte slice will contain 5 bytes with each byte having a value of 5.
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	//appends the padtext to the ciphertext.
	//append() fun takes in 2 slices as arguments & appends elements of the 2nd slice to the elements of the 1st slice
	//& returns the combined slice. i.e ciphertext & padtext slices will be combined and the combined slice will be return.
	return append(ciphertext, padtext...)
}
//AESDecrypt is used to decrypt data using the AES algorithm. 
func AESDecrypt(crypt []byte, key []byte, IV string) []byte {
	// Create the AES cipher block from key
	block, err := aes.NewCipher(key)

	if err != nil {
		fmt.Println("key error1", err)
	}
	if len(crypt) == 0 {
		fmt.Println("plain content empty")
	}
	//Create a new CBC encryption from block and IV
	//The CBC decrypter will be used to decrypt data that has been encrypted with a CBC cipher.
	//The IV is used to help make the decryption process more secure, as it helps to ensure that
	//the same plaintext data does not produce the same ciphertext each time it is encrypted.
	ecb := cipher.NewCBCDecrypter(block, []byte(IV))
	//Creates a new slice of bytes called decrypted that is the same length as the existing slice of bytes called crypt.
	//decrypt the crypt slice of bytes and store the results in the decrypted slice of bytes.
	decrypted := make([]byte, len(crypt))
	//ECB mode of the AES to encrypt a block of data.
	//It takes two parameters: the decrypted data (decrypted) and the key used to encrypt the data (crypt).
	//It then encrypts the data using the AES algorithm and stores the encrypted data in the decrypted parameter.
	ecb.CryptBlocks(decrypted, crypt)
	//PKCS5Trimming is used to trim the decrypted data.
	//This will take the decrypted data,remove any extra padding that was added and return the trimmed data.
	return PKCS5Trimming(decrypted)
}

// PKCS5Trimming  func removes padding from the byte array of encrypted data
// as an argument and returns a new byte array with the padding removed.
func PKCS5Trimming(encrypt []byte) []byte {
	//Set variable as padding to last element of the slice as encrypt.
	// variable padding is set to the last element of the slice because last element of slice is padding used in encryption.
	padding := encrypt[len(encrypt)-1]
	//Remove padding from a string (encrypt).
	//int(padding)is convert the padding variable to an integer.
	//len(encrypt)is used to get the length of the string.
	//[:len(encrypt)-int(padding)] is used to remove the last int(padding) characters from the string.
	return encrypt[:len(encrypt)-int(padding)]
}
