package handler

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"simplified-tik-tok/biz/common"
	"simplified-tik-tok/biz/dal/mysql"
	"simplified-tik-tok/biz/model"
	"simplified-tik-tok/biz/mw"
	utils2 "simplified-tik-tok/biz/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// User login handler
func PublishAction(_ context.Context, c *app.RequestContext) {
	var publishActionStruct struct {
		ActionType string                `form:"title" json:"title" query:"title"`
		Token      string                `form:"token" json:"token" query:"token"`
		Data       *multipart.FileHeader `form:"data" vd:"?"`
	}

	if err := c.BindAndValidate(&publishActionStruct); err != nil {
		c.JSON(http.StatusOK, common.NewParameterError("feed param form wrong"))
		return
	}

	userID, err := mw.GetIDFromTokenString(publishActionStruct.Token)
	if err != nil {
		c.JSON(http.StatusOK, common.NewPrivilegeError("no privilege"))
		return
	}

	file := publishActionStruct.Data
	timeStamp := time.Now().Unix()
	relativePath := strconv.Itoa(int(userID)) + "-" + strconv.FormatInt(timeStamp, 10)
	videoFilePath := common.VideoPath + relativePath + path.Ext(file.Filename) // 重命名视频名称
	coverFilePath := common.CoverPath + relativePath + ".png"                  // 重命名视频名称

	err = c.SaveUploadedFile(file, fmt.Sprint(videoFilePath)) // 存储文件
	if err != nil {
		c.JSON(http.StatusOK, common.NewFileError("video save error"))
		return
	}

	utils2.GetCompressedVideo(videoFilePath)

	err = utils2.GetSnapshot(videoFilePath, coverFilePath, 1)
	if err != nil {
		os.Remove(videoFilePath)
		c.JSON(http.StatusOK, common.NewFileError("cover convert error"))
		return
	}

	title := c.PostForm("title")

	// 视频信息存储到数据库中
	videoRecord := model.Video{
		AuthorID:      userID,
		Createtime:    timeStamp,
		Title:         title,
		Filepath:      videoFilePath,
		Playurl:       common.PlayerURLPrefix + relativePath + path.Ext(file.Filename),
		Coverurl:      common.CoverURLRrefix + relativePath + ".png",
		Favoritecount: 0,
		Commentcount:  0,
	}
	// 同时删除文件
	if err = mysql.CreateVideo(&videoRecord); err != nil {
		os.Remove(videoFilePath)
		os.Remove(coverFilePath)
		c.JSON(http.StatusOK, common.NewDatabaseError("insert video info wrong"))
		return
	}
	c.JSON(http.StatusOK, utils.H{
		"status_code": 0,
		"status_msg":  "success",
	})
}
