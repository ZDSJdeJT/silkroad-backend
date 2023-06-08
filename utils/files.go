package utils

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func MergeFiles(chunksPath string, mergedPath string) error {
	// 确保合并文件所在目录存在
	err := os.MkdirAll(filepath.Dir(mergedPath), 0755)
	if err != nil {
		return err
	}

	// 打开输出文件
	mergedFile, err := os.Create(mergedPath)
	if err != nil {
		return err
	}

	// 确保在函数返回之前关闭输出文件
	defer func() {
		cerr := mergedFile.Close()
		if cerr != nil {
			err = cerr
		}
	}()

	// 遍历所有切片文件，并逐个将它们的数据写入输出文件中
	err = filepath.Walk(chunksPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录和隐藏文件
		if info.IsDir() || filepath.Base(path)[0] == '.' {
			return nil
		}

		// 打开切片文件
		chunkFile, err := os.Open(path)
		if err != nil {
			return err
		}

		// 确保在函数返回之前关闭切片文件
		defer func() {
			cerr := chunkFile.Close()
			if cerr != nil && err == nil {
				err = cerr
			}
		}()

		// 将切片文件的内容复制到输出文件中
		_, err = io.Copy(mergedFile, chunkFile)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func DeleteOldChunks() {
	dirPath := "./data/chunks"
	dir, err := os.Open(dirPath)
	if err != nil {
		log.Printf("Failed to open directory %s: %v", dirPath, err)
		return
	}

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		log.Printf("Failed to read directory contents of %s: %v", dirPath, err)
		err := dir.Close()
		if err != nil {
			log.Printf("Failed to close directory: %v", err)
			return
		}
		return
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			subdirPath := filepath.Join(dirPath, fileInfo.Name())

			modTime := fileInfo.ModTime()
			age := time.Since(modTime)
			if age > 24*time.Hour {
				if err := os.RemoveAll(subdirPath); err != nil {
					log.Printf("Failed to delete directory %s: %v", subdirPath, err)
				}
			}
		}
	}

	if err := dir.Close(); err != nil {
		log.Printf("Failed to close directory %s: %v", dirPath, err)
	}
	log.Printf("Successfully deleted old chunks")
}
