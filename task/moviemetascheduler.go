package task

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/metatube-community/metatube-sdk-go/engine"
	"github.com/metatube-community/metatube-sdk-go/engine/providerid"
	"github.com/metatube-community/metatube-sdk-go/model"
	"gorm.io/gorm"
)

func StartScheduledTasks(db *gorm.DB, e *engine.Engine) {
	go func() {
		log.Println("Scheduled tasks started")
		runNumberStatusTask(db, e)
		for {
			now := time.Now()
			// 计算下一个凌晨1点的时间
			next := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, now.Location())
			if now.After(next) {
				// 如果当前时间已经过了今天的1点，则设置为明天的1点
				next = next.Add(24 * time.Hour)
			}

			// 等待到下一个执行时间
			duration := next.Sub(now)
			log.Printf("Next scheduled task will run at: %v (%v from now)", next.Format("2006-01-02 15:04:05"), duration)

			// 睡眠直到执行时间
			time.Sleep(duration)

			// 执行任务
			log.Printf("Running scheduled task at %v", time.Now().Format("2006-01-02 15:04:05"))
			runNumberStatusTask(db, e)
		}
	}()
}

func runNumberStatusTask(db *gorm.DB, e *engine.Engine) {
	// 查询number_status表中status为0的记录
	var numberStatuses []model.NumberStatus
	if err := db.Where("status = ?", 0).Find(&numberStatuses).Error; err != nil {
		log.Printf("Failed to query number_status table: %v", err)
		return
	}

	log.Printf("Found %d number prefixes with status 0", len(numberStatuses))

	for _, ns := range numberStatuses {
		processNumberPrefix(db, e, ns.NumberPrefix)
	}
}

func processNumberPrefix(db *gorm.DB, e *engine.Engine, numberPrefix string) {
	log.Printf("Processing number prefix: %s", numberPrefix)

	// 查询movie_metadata表中以numberPrefix开头的记录，并找出最大的后缀数字
	var movie model.MovieInfo
	if err := db.Where("number LIKE ?", fmt.Sprintf("%s-%%", numberPrefix)).Order("number DESC").First(&movie).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Failed to query movie_metadata table for prefix %s: %v", numberPrefix, err)
			return
		}
		// 如果没有找到记录，则从0开始
		log.Printf("No records found for prefix %s, starting from 0", numberPrefix)
	}

	// 找出最大的后缀数字
	maxSuffix := 0
	if movie.ID != "" {
		re := regexp.MustCompile(`(\d{3})$`)
		matches := re.FindStringSubmatch(movie.ID)
		if len(matches) == 2 {
			maxSuffix, _ = strconv.Atoi(matches[1])
		}
	}

	log.Printf("Max suffix for prefix %s is %d", numberPrefix, maxSuffix)

	// 如果最大后缀已经是999，则更新status为1
	if maxSuffix >= 999 {
		if err := db.Model(&model.NumberStatus{}).Where("number_prefix = ?", numberPrefix).Update("status", 1).Error; err != nil {
			log.Printf("Failed to update status for prefix %s: %v", numberPrefix, err)
		} else {
			log.Printf("Updated status to 1 for prefix %s", numberPrefix)
		}
		return
	}

	var fanzaId string
	// 从最大编号开始依次自增到999，获取电影信息
	for i := maxSuffix + 1; i <= 999; i++ {
		number := fmt.Sprintf("%s-%03d", numberPrefix, i)
		//优先通过fanza获取电影信息
		if fanzaId != "" {
			//将fanzaId的最后三位数字替换成i
			re := regexp.MustCompile(`\d{3}$`)
			fanzaId = re.ReplaceAllString(fanzaId, fmt.Sprintf("%03d", i))

			log.Printf("Fetching fanza movie info for id: %s", fanzaId)

			// 调用GetMovieInfoByProviderID方法，Provider是FANZA，id是number，lazy参数是true
			pid := providerid.ProviderID{
				Provider: "FANZA",
				ID:       fanzaId,
			}
			movieInfo, err := e.GetMovieInfoByProviderID(pid, true)
			if err == nil {
				// 成功获取到信息，直接打印
				log.Printf("%s: %s 简介：%s", number, movieInfo.Title, movieInfo.Summary)
				continue
			}
		}

		// 如果不能获取到，则调用SearchMovieAll方法
		searchResults, err := e.SearchMovieAll(number, false)
		if err != nil {
			log.Printf("Failed to search movie for %s: %v", number, err)
			continue
		}

		if len(searchResults) == 0 {
			log.Printf("No search results for %s", number)
			continue
		}

		// 取第一个值
		firstResult := searchResults[0]
		//log.Printf("Found search result for %s: %s from %s", number, firstResult.Title, firstResult.Provider)

		// 再调用GetMovieInfoByProviderID方法
		pid := providerid.ProviderID{
			Provider: firstResult.Provider,
			ID:       firstResult.ID,
		}

		movieInfo, err := e.GetMovieInfoByProviderID(pid, true)
		if err != nil {
			log.Printf("Failed to fetch movie info for %s after search: %v", number, err)
			continue
		}
		if movieInfo.Provider == "FANZA" {
			//fanza存在，后续通过fanza获取电影信息
			fanzaId = movieInfo.ID
		}

		// 直接打印
		log.Printf("Successfully fetched movie info for %s (via search): %s", number, movieInfo.Title)
	}
}
