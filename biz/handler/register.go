package handler

import (
	"context"
	"net/http"

	"simplified-tik-tok/biz/dal/mysql"
	"simplified-tik-tok/biz/model"
	"simplified-tik-tok/biz/mw"
	utils2 "simplified-tik-tok/biz/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// Register user register handler
func Register(ctx context.Context, c *app.RequestContext) {
	var registerStruct struct {
		Username string `form:"username" json:"username" query:"username" vd:"(len($) > 0 && len($) < 128); msg:'Illegal format'"`
		Password string `form:"password" json:"password" query:"password" vd:"(len($) > 0 && len($) < 128); msg:'Illegal format'"`
	}

	if err := c.BindAndValidate(&registerStruct); err != nil {
		c.JSON(http.StatusOK, utils.H{
			"status_code": http.StatusBadRequest,
			"status_msg":  err.Error(),
			"user_id":     0,
			"token":       "",
		})
		return
	}
	users, err := mysql.FindUserByName(registerStruct.Username)
	if err != nil {
		c.JSON(http.StatusOK, utils.H{
			"status_code": http.StatusBadRequest,
			"status_msg":  err.Error(),
			"user_id":     0,
			"token":       "",
		})
		return
	}

	if len(users) != 0 {
		c.JSON(http.StatusOK, utils.H{
			"status_code": http.StatusBadRequest,
			"status_msg":  "user already exists",
			"user_id":     0,
			"token":       "",
		})
		return
	}

	if err = mysql.CreateUsers([]*model.User{
		{
			Username: registerStruct.Username,
			Password: utils2.MD5(registerStruct.Password),
		},
	}); err != nil {
		c.JSON(http.StatusOK, utils.H{
			"status_code": http.StatusBadRequest,
			"status_msg":  err.Error(),
			"user_id":     0,
			"token":       "",
		})
		return
	}

	// TODO: register获取Token方式，再Login一遍
	mw.JwtMiddleware.LoginHandler(ctx, c)
}
