package utils

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"path/filepath"

	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/imroc/req/v3"
	"go.uber.org/multierr"
)

func EnsureDirectoryExists(filePath string) error {
	directory := filepath.Dir(filePath)

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func FindExecutableNames(filter func(name string) bool, extraPaths ...string) ([]string, error) {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return nil, fmt.Errorf("PATH environment variable is empty")
	}

	var results []string
	seen := make(map[string]struct{})
	var errs error

	pathToCheck := extraPaths
	pathToCheck = append(pathToCheck, filepath.SplitList(pathEnv)...)

	for _, dir := range pathToCheck {
		entries, err := os.ReadDir(dir)
		if err != nil {
			// cannot read this directory at all: skip it silently
			continue
		}

		for _, entry := range entries {
			name := entry.Name()
			if _, dup := seen[name]; dup {
				continue
			}
			if !filter(name) {
				continue
			}

			// We've got a candidate name; try to stat it
			info, err := entry.Info()
			if err != nil {
				// record the error for this matching name
				errs = multierr.Append(errs, err)
				continue
			}

			// skip directories or non‐executable
			if info.IsDir() || info.Mode()&0111 == 0 {
				continue
			}

			results = append(results, path.Join(dir, name))
			seen[name] = struct{}{}
		}
	}

	if len(results) > 0 {
		return results, nil
	}
	if errs != nil {
		// return only the aggregated errors
		return nil, errs
	}
	// no matches and no file‐specific errors
	return nil, nil

}

var TmpFileDir = path.Join(os.TempDir(), "vaala-frp-panel-download")

// DownloadFile 下载文件到一个临时文件，返回临时文件路径
func DownloadFile(ctx context.Context, url string, proxyUrl string) (string, error) {
	os.MkdirAll(TmpFileDir, 0777)

	tmpPath, err := os.MkdirTemp(TmpFileDir, "downloads")
	if err != nil {
		return "", err
	}

	tmpFileName := generateRandomFileName("download", ".tmp")
	fileFullPath := path.Join(tmpPath, tmpFileName)

	cli := req.C()
	if len(proxyUrl) > 0 {
		cli = cli.SetProxyURL(proxyUrl)
	}

	err = cli.NewParallelDownload(url).
		SetConcurrency(5).
		SetSegmentSize(1024 * 1024 * 1).
		SetOutputFile(fileFullPath).
		SetFileMode(0777).
		SetTempRootDir(path.Join(TmpFileDir, "downloads_cache")).
		Do()
	if err != nil {
		logger.Logger(ctx).WithError(err).Error("download file from url error")
		return "", err
	}

	return fileFullPath, nil
}

// generateRandomFileName 生成一个随机文件名
func generateRandomFileName(prefix, extension string) string {
	randomStr := randomString(8)
	fileName := fmt.Sprintf("%s_%s%s", prefix, randomStr, extension)
	return fileName
}

// randomString 生成一个指定长度的随机字符串
func randomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = charset[rand.Intn(len(charset))]
	}

	return string(bytes)
}

// ExtractGZTo decompresses the srcGZ file into a temporary directory,
// renames the extracted file to newName, moves it to destDir, and sets executable permissions (0755).
// It returns the full path of the final file on success.
func ExtractGZTo(srcGZ, newName, destDir string) (string, error) {
	// 1. Open source .gz file
	f, err := os.Open(srcGZ)
	if err != nil {
		return "", fmt.Errorf("failed to open source gzip file %q: %w", srcGZ, err)
	}
	defer f.Close()

	// 2. Create gzip reader
	zr, err := gzip.NewReader(f)
	if err != nil {
		return "", fmt.Errorf("failed to create gzip reader for %q: %w", srcGZ, err)
	}
	defer zr.Close()

	// 3. Create temporary directory
	tmpDir, err := os.MkdirTemp("", "vaala-frp-panel-gz_extract_*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}
	// Note: tmpDir is not auto-deleted. Caller may clean up if desired.

	// 4. Create the output file in the temp directory with the new name
	tmpFilePath := filepath.Join(tmpDir, newName)
	outFile, err := os.Create(tmpFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file %q: %w", tmpFilePath, err)
	}
	defer outFile.Close()

	// 5. Decompress into temp file
	if _, err := io.Copy(outFile, zr); err != nil {
		return "", fmt.Errorf("failed to write decompressed data to %q: %w", tmpFilePath, err)
	}

	// 6. Ensure destination directory exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create destination directory %q: %w", destDir, err)
	}

	// 7. Move the file to the destination directory
	finalPath := filepath.Join(destDir, newName)
	if err := os.Rename(tmpFilePath, finalPath); err != nil {
		return "", fmt.Errorf("failed to move file to %q: %w", finalPath, err)
	}

	// 8. Set executable permission
	if err := os.Chmod(finalPath, 0755); err != nil {
		return "", fmt.Errorf("failed to set executable permission on %q: %w", finalPath, err)
	}

	return finalPath, nil
}
