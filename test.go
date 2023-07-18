package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	url := "http://www.1123123baidu.com"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", "BAIDUID=FBFBD8ED0351289C2E57C40CDB3E51E7:FG=1; BIDUPSID=FBFBD8ED0351289CF59534D6B871C5CC; H_PS_PSSID=36551_38642_38831_39027_39023_38942_39014_38820_38824_26350_39041_39095_39100; PSTM=1689693729; BDSVRTM=44; BD_HOME=1")

	res, err := client.Do(req)
	fmt.Println(err)
	fmt.Println(123)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
