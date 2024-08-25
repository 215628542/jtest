package shangyi

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// 加解密
func Decrypt(rawData string, source string) (raw string, err error) {
	var keyIvMap = make(map[string]map[string]string, 0)
	keyIvMap["company"] = make(map[string]string, 0)
	keyIvMap["company"]["key"] = "cw8jYudat!fHKnLkpsTm^VJ&5PylBMW7"
	keyIvMap["company"]["iv"] = "acOtV2Bb3QZIJYdR"
	keyIvMap["vmall"] = make(map[string]string, 0)
	keyIvMap["vmall"]["key"] = "1CkzgEcMPXtVix9NBJa3RLW2FlpBobw5"
	keyIvMap["vmall"]["iv"] = "KMUcC9WOPc8ezdGE"

	var key, iv string
	if source == "" {
		source = "company"
	}
	key = keyIvMap[source]["key"]
	iv = keyIvMap[source]["iv"]

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
