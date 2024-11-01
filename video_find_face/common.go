package video_find_face

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// copyFile 复制文件，如果目标文件已存在则覆盖
func CopyFile(src, dst string) error {
	var srcF io.Reader
	var err error

	if strings.HasPrefix(src, "http") {
		resp, err := http.Get(src)
		if err != nil {
			return fmt.Errorf("failed to download source file: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to download source file: status code %d", resp.StatusCode)
		}

		srcF = resp.Body
	} else {
		// 处理本地文件
		srcF, err = os.Open(src)
		if err != nil {
			return fmt.Errorf("failed to open source file: %w", err)
		}
		defer func() {
			if e := srcF.(*os.File).Close(); e != nil && err == nil {
				err = e
			}
		}()
	}

	if err != nil {
		return err
	}

	dstF, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		if e := dstF.Close(); e != nil && err == nil {
			err = e
		}
	}()

	written, err := io.Copy(dstF, srcF)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}
	if written == 0 {
		return fmt.Errorf("failed to written len: 0", err)
	}

	err = dstF.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}

func ExtractIP(url string) string {
	parts := strings.Split(url, "@")
	if len(parts) > 1 {
		hostPart := parts[1]
		ip := strings.Split(hostPart, ":")[0]
		return ip
	}
	return ""
}
