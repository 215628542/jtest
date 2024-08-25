package mutilMapCache

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"hash/crc32"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func printAlloc(prefix string) {

	//runtime.KeepAlive(m) // Keeps a reference to m so that the map isn’t collected
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("%s %d KB\n", prefix, memStats.Alloc/1024)
}

func fillMap(bigMap *BigMap) {
	for i := 0; i < 1000000; i++ {

		msg := strings.Repeat("test", 300)
		//msg := "300"
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(100000000)
		bigMap.Set(cast.ToString(randNum), &msg)
	}

}
func clearMap(bigMap *BigMap) {
	//bigMap.flushExpireData()

	var i uint32 = 0
	for ; i < bigMap.bucketSize; i++ {

		for k, _ := range bigMap.cacheBuckets[i] {
			delete(bigMap.cacheBuckets[i], k)
		}
	}
}

var BigSyncMapKey1 = "test1"

//func TestMapCache3(t *testing.T) {
//
//	go AestMapCache2()
//
//	r := gin.Default()
//	r.GET("/ping", func(c *gin.Context) {
//		runtime.GC()
//		c.JSON(http.StatusOK, gin.H{
//			"message": "pong123",
//		})
//	})
//	r.Use()
//	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
//
//}

func TestMapCache3(t *testing.T) {

	InitMapCache(BigSyncMapKey1, 1, 1*time.Second, 2*time.Second)
	bigMap, err := GetBigMaps(BigSyncMapKey1)
	if err != nil {
		panic(err)
	}

	printAlloc("0 init")
	i := 0
	for {
		fillMap(bigMap)
		i++
		time.Sleep(3 * time.Second)
		printAlloc(fmt.Sprintf("%d after fillMap", i))

		clearMap(bigMap)
		i++
		time.Sleep(3 * time.Second)
		printAlloc(fmt.Sprintf("%d after clearMap", i))
		//runtime.GC()
		//i++
		//printAlloc(fmt.Sprintf("%d after clearMap and runtime.GC()", i))

		//DelBigMaps(BigSyncMapKey1)
		//time.Sleep(2 * time.Second)
		//printAlloc(fmt.Sprintf("%d delete bigMap", i))

		//InitMapCache(BigSyncMapKey1, 1, 1*time.Second, 2*time.Second)
		//bigMap, err = GetBigMaps(BigSyncMapKey1)
		//if err != nil {
		//	panic(err)
		//}
		//time.Sleep(2 * time.Second)
		//printAlloc(fmt.Sprintf("%d after clearMap", i))

		bigMap.flushExpireData()
		i++
		//runtime.GC()
		time.Sleep(2 * time.Second)
		printAlloc(fmt.Sprintf("%d bigMap flushExpireData", i))
	}
}

// =======================

var once sync.Once
var bigMaps map[string]*BigMap

type BigMap struct {
	// 缓存key
	mapKey string

	// 缓存数据桶数
	bucketSize uint32

	// 缓存数据过期清除时间
	expireTime time.Duration

	// 清除缓存数据定时器
	expireTicker time.Duration

	// 清除缓存数据时间
	flushTime time.Time

	// 分段来存储，提高数据并发访问效率
	cacheBuckets []map[string]*string
}

func InitMapCache(mapKey string, bucketSize uint32, expireTime, expireTicker time.Duration) error {

	_, err := NewBigMaps(mapKey, bucketSize, expireTime, expireTicker)
	if err != nil {
		return err
	}
	//go bigMap.flushExpireData()
	return nil
}

func NewBigMaps(mapKey string, bucketSize uint32, expireTime, expireTicker time.Duration) (*BigMap, error) {

	once.Do(func() {
		fmt.Println("=============== new ====================")
		bigMaps = make(map[string]*BigMap, 0)
	})
	bigMap, ok := bigMaps[mapKey]
	if ok {
		return bigMap, nil
	}

	if bucketSize < 1 {
		return nil, errors.New("缓存桶数量不能小于1")
	}
	if mapKey == "" {
		return nil, errors.New("mapKey不能为空")
	}
	if expireTime < 0 {
		return nil, errors.New("请设置缓存过期时间")
	}
	if expireTicker < 0 {
		return nil, errors.New("请设置缓存过期定时器执行时间")
	}

	bigMap = &BigMap{
		mapKey:       mapKey,
		bucketSize:   bucketSize,
		expireTime:   expireTime,
		expireTicker: expireTicker,
		flushTime:    time.Now(),
		cacheBuckets: make([]map[string]*string, bucketSize),
	}

	var i uint32 = 0
	for ; i < bucketSize; i++ {
		bigMap.cacheBuckets[i] = make(map[string]*string, 0)
	}
	bigMaps[mapKey] = bigMap

	return bigMap, nil
}

func GetBigMaps(mapKey string) (*BigMap, error) {
	bigMap, ok := bigMaps[mapKey]
	if ok {
		return bigMap, nil
	}
	return nil, errors.New("bigMap不存在")
}

func DelBigMaps(mapKey string) {
	bigMap, ok := bigMaps[mapKey]
	if ok {
		bigMap.cacheBuckets = make([]map[string]*string, 0)
	}
}

// 获取数据
func (b *BigMap) Get(key string) (interface{}, bool) {
	index := checksum(key, b.bucketSize)
	if index >= len(b.cacheBuckets) {
		return nil, false
	}

	val, ok := b.cacheBuckets[index][key]
	return val, ok
}

// 保存数据，并更新对应key的expire time
func (b *BigMap) Set(key string, val *string) error {
	index := checksum(key, b.bucketSize)
	if index >= len(b.cacheBuckets) {
		return errors.New("index out of range")
	}
	b.cacheBuckets[index][key] = val
	return nil
}

func (b *BigMap) delete(key string) error {
	index := checksum(key, b.bucketSize)
	if index >= len(b.cacheBuckets) {
		return errors.New("index out of range")
	}
	delete(b.cacheBuckets[index], key)
	return nil
}

// 根据crc32算法取key的index
func checksum(key string, bucketSize uint32) int {
	return int(crc32.ChecksumIEEE([]byte(key)) % bucketSize)
}

func (b *BigMap) flushExpireData() {

	b.cacheBuckets = make([]map[string]*string, b.bucketSize)
	var i uint32 = 0
	for ; i < b.bucketSize; i++ {
		b.cacheBuckets[i] = make(map[string]*string, 0)
	}
	//for k, _ := range b.cacheBuckets {
	//	b.cacheBuckets[k] = make(map[string]interface{}, 0)
	//}

	//ticker := time.NewTicker(b.expireTicker)
	//for {
	//	select {
	//	case <-ticker.C:
	//		if time.Since(b.flushTime) >= b.expireTime {
	//			b.cacheBuckets = make([]map[string]interface{}, 0)
	//			b.flushTime = time.Now()
	//		}
	//	}
	//}

}
