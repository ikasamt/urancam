package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func convertToHLS(filePath string) (filenames []string, err error) {
	// HLS形式のファイル名（m3u8）
	baseFilePath := strings.TrimSuffix(filePath, filepath.Ext(filePath))
	hlsPlaylistPath := baseFilePath + ".m3u8"

	// HLSセグメントファイルのパターン（ts）
	hlsSegmentPath := baseFilePath + "_%03d.ts"

	// ffmpegコマンドを構築して実行
	cmd := exec.Command("ffmpeg", "-i", filePath, "-profile:v", "baseline", "-level", "3.0", "-s", "640x360", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-hls_segment_filename", hlsSegmentPath, "-f", "hls", hlsPlaylistPath)
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	// 生成されたSegmentファイルとhls Playlistファイルをfilenamesに格納して返す
	filenames = append(filenames, hlsPlaylistPath)

	for i := 0; i < 10; i++ {
		fn := fmt.Sprintf("%s_%03d.ts", baseFilePath, i)
		if !Exists(fn) {
			break
		}
		filenames = append(filenames, fn)
	}
	return
}
