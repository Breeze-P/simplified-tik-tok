package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/disintegration/imaging"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

// 生成封面
func GetSnapshot(videoPath string, snapshotPath string, frameNum int) (err error) {
	buf := bytes.NewBuffer(nil)
	err = ffmpeg_go.Input(videoPath).Filter("select", ffmpeg_go.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg_go.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()

	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err
	}

	err = imaging.Save(img, snapshotPath)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err
	}

	return nil
}

// 视觉无损压缩
func GetCompressedVideo(src string) {
	cmd := exec.Command("echo", "y", "|", "ffmpeg", "-i", src, "-c:v", "libx264", "-x264-params", "crf=18", src)
	cmd.Run()
}
