package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	RateAlert = errors.New("失败率过高")
)

type DefaultMetricCollector struct {
	successNum   uint32
	failureNum   uint32
	timeoutNum   uint32
	rejectionNum uint32
}

type Window struct {
	Buckets map[int64]*windowBucket
	Mutex   *sync.RWMutex
	Volume  uint32  // 时间窗口内熔断器的生效起始请求数
	Rate    float64 // 失败率 大于此失败率请求会熔断
}

type windowBucket struct {
	Value DefaultMetricCollector
}

func NewWindow(volume uint32, rate float64) *Window {
	r := &Window{
		Buckets: make(map[int64]*windowBucket),
		Mutex:   &sync.RWMutex{},
		Volume:  volume,
		Rate:    rate,
	}

	return r
}

func (r *Window) getCurrentBucket() *windowBucket {
	var bucket *windowBucket
	var ok bool

	time := time.Now().UnixNano() / 1e8 //每100毫秒作为一个桶

	//创建新的桶
	if bucket, ok = r.Buckets[time]; !ok {
		bucket = &windowBucket{}
		r.Buckets[time] = bucket

		//删除过期bucket
		r.removeOldBuckets()
	}

	return bucket

}

func (r *Window) removeOldBuckets() {

	time := time.Now().UnixNano()/1e8 - 10 //剔除1s前的桶

	for timestamp := range r.Buckets {
		if timestamp <= time {
			delete(r.Buckets, timestamp)
		}
	}

}

func (r *Window) Statistics(code int) error {

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	//获取当前bucket
	bucket := r.getCurrentBucket()

	switch code {
	case 0:
		bucket.Value.successNum++
	case 1:
		bucket.Value.failureNum++
	case 2:
		bucket.Value.timeoutNum++
	case 3:
		bucket.Value.rejectionNum++
	}

	total := bucket.Value.successNum + bucket.Value.failureNum + bucket.Value.timeoutNum + bucket.Value.rejectionNum
	if total < r.Volume {
		return nil
	}

	rate := float64(bucket.Value.failureNum+bucket.Value.timeoutNum+bucket.Value.rejectionNum) / float64(total)
	if rate < r.Rate {
		return nil
	}

	return errors.Wrap(RateAlert, fmt.Sprintf("total: %d,success: %d\r\nrate:%f\r\n达到预警值", total, bucket.Value.successNum, rate))
}

func main() {
	window := NewWindow(3, 0.67)

	for i := 0; i < 20; i++ {
		err := window.Statistics(rand.Intn(4))
		if errors.Is(err, RateAlert) {
			fmt.Printf("%v\n", err)
			time.Sleep(1 * time.Second)
		} else {
			fmt.Printf("success:%d\n", i)
			time.Sleep(10 * time.Microsecond)
		}

	}

}
