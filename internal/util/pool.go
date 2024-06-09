package util

import (
	"os"
	"strconv"
	"sync"
)

var (
	bufPool *sync.Pool
	once    sync.Once
	size    int
)

func init() {
	size = getSizeFromEnv()
}

func getSizeFromEnv() int {
	envSize := os.Getenv("BUF_POOL_SIZE")
	if envSize == "" {
		return 1024 // 默认缓冲区大小为1024
	}
	size, err := strconv.Atoi(envSize)
	if err != nil {
		return 1024 // 如果转换失败，使用默认值
	}
	return size
}

// GetBufPool returns the singleton buffer pool, initializing it if necessary
func GetBuf() []byte {
	once.Do(func() {
		bufPool = &sync.Pool{
			New: func() interface{} {
				buf := make([]byte, size)
				return &buf
			},
		}
	})
	// return bufPool.Get().([]byte)
	return *(bufPool.Get().(*[]byte))
}

func PutBuf(data []byte) {
	once.Do(func() {
		bufPool = &sync.Pool{
			New: func() any {
				buf := make([]byte, size)
				return buf
			},
		}
	})
	bufPool.Put(&data)
}
