package elasticSearch

import (
	"github.com/olivere/elastic"
	"net/http"
)

var es *elastic.Client

func init() {
	var err error
	es, err = elastic.NewClient(
		elastic.SetHttpClient(&http.Client{Transport: &http.Transport{Proxy: nil}}),
		elastic.SetSniff(false),
		elastic.SetURL("http://localhost:9200"), // URL地址
		elastic.SetBasicAuth("", ""),            // 账号密码
		//elastic.SetTraceLog(log.New(os.Stdout, "ELASTIC ", log.LstdFlags)),
	)
	if err != nil {
		panic(err.Error())
	}
}
