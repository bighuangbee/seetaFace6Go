package common

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var dateTime time.Time

// 获取图片的时间戳
func getImageTimestamp(imageName string) (time.Time, error) {
	re := regexp.MustCompile(`\d{14}`) // 匹配时间戳 20241018143015
	match := re.FindString(imageName)
	if match == "" {
		return time.Time{}, fmt.Errorf("图片名称中未找到时间戳")
	}

	// 按照时间戳格式解析时间
	return time.Parse("20060102150405", match)
}

// 获取视频文件的时间范围
func getVideoTimeRange(videoName string) (time.Time, time.Time, error) {
	re := regexp.MustCompile(`(\d{2}\.\d{2}\.\d{2})-(\d{2}\.\d{2}\.\d{2})`) // 匹配14.45.55-14.46.52格式
	match := re.FindStringSubmatch(videoName)
	if len(match) < 3 {
		return time.Time{}, time.Time{}, fmt.Errorf("视频文件名格式不正确" + videoName)
	}

	startTime, err := parseVideoTime(match[1])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endTime, err := parseVideoTime(match[2])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return startTime, endTime, nil
}

// 解析视频文件的时间格式为 Time 对象
func parseVideoTime(videoTime string) (time.Time, error) {
	parts := strings.Split(videoTime, ".")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("无效的时间格式: %s", videoTime)
	}

	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])
	second, _ := strconv.Atoi(parts[2])

	// 当前日期加上时间
	videoTimestamp := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), hour, minute, second, 0, dateTime.Location())

	return videoTimestamp, nil
}

// 查找匹配的文件
func FindMatchingVideo(imageName string, videoDir string) (string, error) {
	imageTime, err := getImageTimestamp(imageName)
	if err != nil {
		return "", err
	}

	dateTime = imageTime

	filesname, err := GetFilesName(videoDir)
	if err != nil {
		return "", err
	}

	for _, filename := range filesname {
		if filepath.Ext(filename) == ".mp4" || filepath.Ext(filename) == ".dav" {
			// 获取视频文件的时间范围
			startTime, endTime, err := getVideoTimeRange(filename)
			if err != nil {
				return "", err
			}

			// 检查图片时间是否在视频时间范围内
			if (imageTime.Equal(startTime) || imageTime.After(startTime)) && (imageTime.Equal(endTime) || imageTime.Before(endTime)) {
				return filename, nil
			}
		}
	}

	return "", errors.New("未找到")
}

func GetFilesName(dir string) (files []string, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}
