package handler

import (
	"context"
	"fmt"
	"net/http"

	"simplified-tik-tok/biz/common"
	"simplified-tik-tok/biz/dal/mysql"
	"simplified-tik-tok/biz/mw"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func RelationAction(_ context.Context, c *app.RequestContext) {
	var relationActionStruct struct {
		ToUserID   int64  `form:"to_user_id" json:"to_user_id" query:"to_user_id"`
		ActionType int32  `form:"action_type" json:"action_type" query:"action_type" vd:"($ == 1 || $ == 2); msg:'Illegal format'"`
		Token      string `form:"token" json:"token" query:"token"`
	}

	if err := c.BindAndValidate(&relationActionStruct); err != nil {
		c.JSON(http.StatusOK, common.NewParameterError("feed param form wrong"))
		return
	}

	fmt.Println(relationActionStruct.Token)

	userID, err := mw.GetIDFromTokenString(relationActionStruct.Token)
	if err != nil {
		c.JSON(http.StatusOK, common.NewPrivilegeError("no privilege"))
		return
	}
	// 关注
	if relationActionStruct.ActionType == 1 {
		err := mysql.CreateRelation(int64(userID), relationActionStruct.ToUserID)
		if err != nil {
			c.JSON(http.StatusOK, common.NewDatabaseError("create relation wrong"))
			return
		}
		err = mysql.AddUserFollowCountByID(userID)
		if err != nil {
			c.JSON(http.StatusOK, common.NewDatabaseError("add relation wrong"))
			return
		}
		err = mysql.AddUserFansCountByID(uint(relationActionStruct.ToUserID))
		if err != nil {
			c.JSON(http.StatusOK, common.NewDatabaseError("add relation wrong"))
			return
		}
		c.JSON(http.StatusOK, utils.H{
			"status_code": 0,
			"status_msg":  "success",
		})
		return
	}

	// 取关
	err = mysql.DeleteRelation(int64(userID), relationActionStruct.ToUserID)
	if err != nil {
		c.JSON(http.StatusOK, common.NewDatabaseError("delete relation wrong"))
		return
	}
	err = mysql.DecUserFollowCountByID(userID)
	if err != nil {
		c.JSON(http.StatusOK, common.NewDatabaseError("delete relation wrong"))
		return
	}
	err = mysql.DecUserFansCountByID(uint(relationActionStruct.ToUserID))
	if err != nil {
		c.JSON(http.StatusOK, common.NewDatabaseError("delete relation wrong"))
		return
	}
	c.JSON(http.StatusOK, utils.H{
		"status_code": 0,
		"status_msg":  "success",
	})
}
