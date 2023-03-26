package mw

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/jwt"

	"simplified-tik-tok/biz/dal/mysql"
	"simplified-tik-tok/biz/model"
	utils2 "simplified-tik-tok/biz/utils"
)

var (
	JwtMiddleware *jwt.HertzJWTMiddleware
	IdentityKey   = "identity"
)

func InitJwt() {
	var err error
	JwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:         "test zone",
		Key:           []byte("secret key"),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		LoginResponse: func(_ context.Context, c *app.RequestContext, code int, token string, _ time.Time) {
			// TODO
			// 获取UserID方案：1、查库；2、从Token中拿
			// 1、
			var loginStruct struct {
				Username string `form:"username" json:"username" query:"username" vd:"(len($) > 0 && len($) < 32); msg:'Illegal format'"`
				Password string `form:"password" json:"password" query:"password" vd:"(len($) > 0 && len($) < 32); msg:'Illegal format'"`
			}
			if err := c.BindAndValidate(&loginStruct); err != nil {
				c.JSON(http.StatusBadRequest, utils.H{
					"status_code": code,
					"status_msg":  "wrong",
					"user_id":     0,
					"token":       token,
				})
				return
			}
			users, err := mysql.FindUserByName(loginStruct.Username)

			if err != nil || len(users) == 0 {
				c.JSON(http.StatusBadRequest, utils.H{
					"status_code": code,
					"status_msg":  "wrong",
					"user_id":     0,
					"token":       token,
				})
			}
			c.JSON(http.StatusOK, utils.H{
				"status_code": 0,
				"status_msg":  "success",
				"user_id":     users[0].ID,
				"token":       token,
			})
		},
		Authenticator: func(_ context.Context, c *app.RequestContext) (interface{}, error) {
			var loginStruct struct {
				Username string `form:"username" json:"username" query:"username" vd:"(len($) > 0 && len($) < 32); msg:'Illegal format'"`
				Password string `form:"password" json:"password" query:"password" vd:"(len($) > 0 && len($) < 32); msg:'Illegal format'"`
			}
			if err := c.BindAndValidate(&loginStruct); err != nil {
				return nil, err
			}
			users, err := mysql.CheckUser(loginStruct.Username, utils2.MD5(loginStruct.Password))
			if err != nil {
				return nil, err
			}
			if len(users) == 0 {
				return nil, errors.New("user already exists or wrong password")
			}
			return users[0], nil
		},
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					IdentityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		HTTPStatusMessageFunc: func(e error, ctx context.Context, _ *app.RequestContext) string {
			hlog.CtxErrorf(ctx, "jwt biz err = %+v", e.Error())
			return e.Error()
		},
		Unauthorized: func(_ context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(http.StatusOK, utils.H{
				"status_code": code,
				"status_msg":  message,
				"user_id":     0,
				"token":       "",
			})
		},
	})
	if err != nil {
		panic(err)
	}
}

func GetIDFromToken(ctx context.Context, c *app.RequestContext) (uint, error) {
	if JwtMiddleware == nil {
		return 0, errors.New("jwt not init")
	}
	token, err := JwtMiddleware.ParseToken(ctx, c)
	if err != nil {
		return 0, err
	}
	claims := jwt.ExtractClaimsFromToken(token)
	res := uint(claims[IdentityKey].(float64))
	return res, nil
}

func GetIDFromTokenString(tokenString string) (uint, error) {
	if JwtMiddleware == nil {
		return 0, errors.New("jwt not init")
	}
	token, err := JwtMiddleware.ParseTokenString(tokenString)
	if err != nil {
		return 0, err
	}
	claims := jwt.ExtractClaimsFromToken(token)
	res := uint(claims[IdentityKey].(float64))
	return res, nil
}
