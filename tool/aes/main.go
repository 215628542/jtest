package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

func main() {

	data := "YlGAJ/j52i0Y6vnQ4XquuGGsDbRcpRi0mc77Ru53bz3P/f4EIrgLrE53cRWfNwC7cdjHHxR5955/WxLU/gS5x/7hELmYkqvNbOJ74bDx9cU="
	apiAeskey := "MTWJuf7l7R0wcqQjAPq8TitslIbaK71t"
	apiIV := "1234567890abcdef"
	decode, err := Decrypt(data, apiAeskey, apiIV)
	fmt.Println(err)
	fmt.Println(decode)

}

func Decrypt(rawData string, key string, iv string) (raw string, err error) {
	defer func() {
		if errStr := recover(); errStr != nil {
			err = errors.New("请求异常")
		}
		return
	}()
	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		return
	}
	dByte := aesCbcDecrypt(data, []byte(key), iv)
	raw = string(dByte)
	return
}

// AEC解密（CBC模式）
func aesCbcDecrypt(cipherText []byte, key []byte, iv string) []byte {
	//指定解密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//指定初始化向量IV,和加密的一致
	ivByte := []byte(iv)
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCDecrypter(block, ivByte)
	//解密
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)
	//删除填充
	plainText = unPadding(plainText)
	return plainText
}

// 对密文删除填充 PKCS7UnPadding
func unPadding(cipherText []byte) []byte {
	//取出密文最后一个字节end
	end := cipherText[len(cipherText)-1]
	//删除填充
	cipherText = cipherText[:len(cipherText)-int(end)]
	return cipherText
}
