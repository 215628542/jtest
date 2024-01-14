package elasticSearch

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"test/model"
)

func Exec() {

	//var id uint64 = 1
	//ctx := context.Background()
	//createOrder(ctx, id)
	getData()
}

func createOrder(ctx context.Context, id uint64) {
	m := model.EsVmallOrder{Id: id}
	err := m.GetById(ctx)
	if err != nil {
		panic(err)
	}

	result, err := es.Index().Index("vmall_orders").Id(cast.ToString(id)).
		BodyJson(m).Refresh("true").Do(ctx)

	if err != nil {
		panic(err)
	} else if result == nil {
		panic("create fail")
	}
}

func getData() {
	ctx := context.Background()
	result, err := es.Search().Index("vmall_orders").
		Size(10).From(0).
		Sort("id", false).
		Pretty(true).Do(ctx)

	if err != nil {
		panic(err)
	}

	// 解析数据
	for _, hit := range result.Hits.Hits {
		hitByte, _ := hit.Source.MarshalJSON()
		fmt.Println(string(hitByte))
	}

}
