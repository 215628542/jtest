package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cast"
	"sort"
	"strings"
)

// 瑞银信签名
func GenerateSignForData() {
	data := map[string]interface{}{
		"context":       "瑞银信测试商品-0423 数量（1）",
		"money":         "0.01",
		"pay_time":      "2023-04-23 11:07:18",
		"shop_id":       "2dd969ad8164c3b651f803ae8cb8447f",
		"shop_order_id": "vm407854005264519168",
		//"sign":          "AD4220A0B7C27FA16FA319C697B34848",
		"source": "wx",
	}
	generateSign(data)
}

func generateSign(requestBody map[string]interface{}) {
	keys := getMapSortKeys(requestBody)
	stringA := combineMapByKeys(keys, requestBody)

	signKey := "BC134BE67BDD1B3ACB98228C79B34D35"

	stringSignTemp := fmt.Sprintf("%s&key=%s", stringA, signKey)
	fmt.Println("签名字符串:", stringSignTemp)
	sg := strings.ToUpper(md5Encrypt(stringSignTemp))
	fmt.Println("签名:", sg)
}

func getMapSortKeys(data map[string]interface{}) []string {
	keys := make([]string, 0)
	for key := range data {
		val := cast.ToString(data[key])

		if val != "" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	return keys
}

func combineMapByKeys(keys []string, data map[string]interface{}) (res string) {
	for _, key := range keys {
		res += fmt.Sprintf("%s=%v&", key, data[key])
	}
	return strings.Trim(res, "&")
}

func md5Encrypt(sec string) string {
	m := md5.New()
	m.Write([]byte(sec))
	return hex.EncodeToString(m.Sum(nil))
}
