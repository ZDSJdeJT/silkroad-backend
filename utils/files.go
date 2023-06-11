package utils

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func WriteToFile(filePath string, content string) error {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return errors.New("failed to create directory: " + err.Error())
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.New("failed to open file: " + err.Error())
	}

	stat, err := file.Stat()
	if err == nil && stat.IsDir() {
		if ferr := file.Close(); ferr != nil {
			return errors.New("failed to close file: " + ferr.Error())
		}
		return errors.New("file path is a directory")
	}

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(content)
	if err != nil {
		if werr := writer.Flush(); werr != nil {
			return errors.New("failed to flush buffer: " + werr.Error())
		}
		if ferr := file.Close(); ferr != nil {
			return errors.New("failed to close file: " + ferr.Error())
		}
		return errors.New("failed to write file: " + err.Error())
	}

	err = writer.Flush()
	if err != nil {
		if fcerr := file.Close(); fcerr != nil {
			return errors.New("failed to close file: " + fcerr.Error())
		}
		return errors.New("failed to flush buffer: " + err.Error())
	}

	if ferr := file.Close(); ferr != nil {
		return errors.New("failed to close file: " + ferr.Error())
	}

	return nil
}

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

func DeleteExpiredChunks() {
	dirPath := "./data/uploads/chunks"
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

func GetDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
		if !fileInfo.IsDir() {
			size += fileInfo.Size()
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return size, nil
}
