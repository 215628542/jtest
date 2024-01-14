package main

import (
	"fmt"
	"github.com/pkg/errors"
	"hash/crc32"
	"sync"
	"testing"
	"time"
)

// 多级缓存
func TestMapCache(t *testing.T) {

	InitSyncMapCache(BigSyncMapKey1, 3, 6*time.Second, 3*time.Second)
	syncMap1, err := GetBigSyncMaps("test1")
	if err != nil {
		panic(err)
	}
	err = syncMap1.Set("aa", "aa_value")
	if err != nil {
		panic(err)
	}
	v1, _ := syncMap1.Get("aa")
	fmt.Println(v1)

	InitSyncMapCache(BigSyncMapKey2, 3, 8*time.Second, 4*time.Second)
	syncMap2, err := GetBigSyncMaps("test2")
	if err != nil {
		panic(err)
	}
	err = syncMap2.Set("bb", "bb_value")
	if err != nil {
		panic(err)
	}
	v2, _ := syncMap2.Get("bb")
	fmt.Println(v2)

	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-ticker.C:
			v2, bo := syncMap2.Get("bb")
			fmt.Println("======")
			fmt.Println(v2)
			fmt.Println(bo)
		}
	}
}

// ====== syncMap类

const (
	BigSyncMapKey1 = "test1"
	BigSyncMapKey2 = "test2"
)

var once sync.Once
var bigSyncMaps map[string]*BigSyncMap

type BigSyncMap struct {
	// 缓存key
	mapKey string

	// 缓存数据桶数
	bucketSize uint32

	// 缓存数据过期清除时间
	expireTime time.Duration

	// 清除缓存数据定时器
	expireTicker time.Duration

	// 分段来存储，提高数据并发访问效率
	cacheBuckets []*sync.Map

	// 存储每个key的过期时间
	expires *sync.Map
}

func InitSyncMapCache(mapKey string, bucketSize uint32, expireTime, expireTicker time.Duration) error {

	bigSyncMap, err := NewBigSyncMaps(mapKey, bucketSize, expireTime, expireTicker)
	if err != nil {
		return err
	}
	go bigSyncMap.flushExpireData()
	return nil
}

func NewBigSyncMaps(mapKey string, bucketSize uint32, expireTime, expireTicker time.Duration) (*BigSyncMap, error) {

	once.Do(func() {
		fmt.Println("=============== new ====================")
		bigSyncMaps = make(map[string]*BigSyncMap, 0)
	})
	bigSyncMap, ok := bigSyncMaps[mapKey]
	if ok {
		return bigSyncMap, nil
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

	bigSyncMap = &BigSyncMap{
		mapKey:       mapKey,
		bucketSize:   bucketSize,
		expireTime:   expireTime,
		expireTicker: expireTicker,
		cacheBuckets: make([]*sync.Map, bucketSize),
		expires:      &sync.Map{},
	}

	var i uint32 = 0
	for ; i < bucketSize; i++ {
		bigSyncMap.cacheBuckets[i] = &sync.Map{}
	}
	bigSyncMaps[mapKey] = bigSyncMap

	return bigSyncMap, nil
}

func GetBigSyncMaps(mapKey string) (*BigSyncMap, error) {
	bigSyncMap, ok := bigSyncMaps[mapKey]
	if ok {
		return bigSyncMap, nil
	}
	return nil, errors.New("bigSyncMap不存在")
}

// 获取数据
func (b *BigSyncMap) Get(key string) (interface{}, bool) {
	index := checksum(key, b.bucketSize)
	if index >= len(b.cacheBuckets) {
		return nil, false
	}
	return b.cacheBuckets[index].Load(key)
}

// 保存数据，并更新对应key的expire time
func (b *BigSyncMap) Set(key string, val interface{}) error {
	index := checksum(key, b.bucketSize)
	if index >= len(b.cacheBuckets) {
		return errors.New("index out of range")
	}
	b.cacheBuckets[index].Store(key, val)
	b.expires.Store(key, time.Now())
	return nil
}

func (b *BigSyncMap) delete(key string) error {
	index := checksum(key, b.bucketSize)
	if index >= len(b.cacheBuckets) {
		return errors.New("index out of range")
	}
	b.cacheBuckets[index].Delete(key)
	return nil
}

// 根据crc32算法取key的index
func checksum(key string, bucketSize uint32) int {
	return int(crc32.ChecksumIEEE([]byte(key)) % bucketSize)
}

func (b *BigSyncMap) flushExpireData() {
	ticker := time.NewTicker(b.expireTicker)
	for {
		select {
		case <-ticker.C:
			fmt.Println(b.mapKey + "-开始清除过期缓存数据")
			b.expires.Range(func(key, expireTime interface{}) bool {
				if time.Since(expireTime.(time.Time)) > b.expireTime {
					b.delete(key.(string))
					b.expires.Delete(key.(string))
					fmt.Println(b.mapKey + "-已删除缓存key-" + key.(string))
				}
				return true
			})
		}
	}
}
