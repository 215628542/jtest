package mutilMapCache

import (
	"errors"
	"fmt"
	"hash/crc32"
	"sync"
	"time"
)

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
	cacheBuckets []map[string]interface{}
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
		cacheBuckets: make([]map[string]interface{}, bucketSize),
	}

	var i uint32 = 0
	for ; i < bucketSize; i++ {
		bigMap.cacheBuckets[i] = make(map[string]interface{}, 0)
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
		bigMap.cacheBuckets = make([]map[string]interface{}, 0)
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
func (b *BigMap) Set(key string, val interface{}) error {
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
	ticker := time.NewTicker(b.expireTicker)
	for {
		select {
		case <-ticker.C:
			if time.Since(b.flushTime) >= b.expireTime {
				b.cacheBuckets = make([]map[string]interface{}, 0)
				b.flushTime = time.Now()
			}
		}
	}
}
