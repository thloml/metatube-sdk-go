package openaigen

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestFanzaMovieAPIWithGoroutines(t *testing.T) {
	baseURL := "http://100.74.242.89:26801/v1/movies/FANZA"
	lazyParam := "True"

	// 使用信号量控制并发数，最多同时运行3个
	concurrentLimit := 3
	semaphore := make(chan struct{}, concurrentLimit)

	// 使用 WaitGroup 等待所有协程完成
	var wg sync.WaitGroup

	// 记录开始时间
	startTime := time.Now()

	// 循环调用接口，从 pred00001 到 pred00999
	for i := 1; i <= 999; i++ {
		// 格式化ID，确保前导零
		id := fmt.Sprintf("1rctd00%03d", i)

		// 增加 WaitGroup 计数
		wg.Add(1)

		// 启动协程处理请求
		go func(movieID string) {
			// 函数结束时减少 WaitGroup 计数
			defer wg.Done()

			// 获取信号量，控制并发数（最多3个同时运行）
			semaphore <- struct{}{}
			// 函数结束时释放信号量
			defer func() { <-semaphore }()

			// 构建完整的URL
			url := fmt.Sprintf("%s/%s?lazy=%s", baseURL, movieID, lazyParam)

			// 发起HTTP GET请求
			resp, err := http.Get(url)
			if err != nil {
				t.Errorf("Failed to call API for ID %s: %v", movieID, err)
				return
			}

			// 确保响应体被关闭
			defer resp.Body.Close()

			// 检查响应状态码
			if resp.StatusCode != http.StatusOK {
				t.Errorf("API returned non-OK status for ID %s: %d", movieID, resp.StatusCode)
				return
			}

			// 读取响应体
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Failed to read response body for ID %s: %v", movieID, err)
				return
			}

			// 输出响应结果
			t.Logf("ID: %s, Response: %s", movieID, string(body))
		}(id)
	}

	// 等待所有协程完成
	wg.Wait()

	// 输出总耗时
	duration := time.Since(startTime)
	t.Logf("Total execution time: %v", duration)
}
