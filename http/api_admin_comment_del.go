package http

import (
	"strconv"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ParamsCommentDel struct {
	ID string `mapstructure:"id" param:"required"`

	SiteName string `mapstructure:"site_name"`
	SiteID   uint
	SiteAll  bool
}

func ActionAdminCommentDel(c echo.Context) error {
	if isOK, resp := AdminOnly(c); !isOK {
		return resp
	}

	var p ParamsCommentDel
	if isOK, resp := ParamsDecode(c, ParamsCommentDel{}, &p); !isOK {
		return resp
	}

	id, err := strconv.Atoi(p.ID)
	if err != nil {
		return RespError(c, "invalid id")
	}

	// find site
	if isOK, resp := CheckSite(c, &p.SiteName, &p.SiteID, &p.SiteAll); !isOK {
		return resp
	}

	comment := model.FindComment(uint(id), p.SiteName)
	if comment.IsEmpty() {
		return RespError(c, "comment not found")
	}

	if err := DelComment(&comment); err != nil {
		return RespError(c, "comment delete error")
	}

	// 删除子评论
	hasErr := false
	children := comment.FetchChildren(func(db *gorm.DB) *gorm.DB { return db })
	for _, c := range children {
		err := DelComment(&c)
		if err != nil {
			hasErr = true
		}
	}
	if hasErr {
		return RespError(c, "children comment delete error")
	}

	return RespSuccess(c)
}

func DelComment(comment *model.Comment) error {
	// 清除 notify
	var notifications []model.Notify
	lib.DB.Where("comment_id = ?", comment.ID).Find(&notifications)

	for _, n := range notifications {
		if err := lib.DB.Delete(n).Error; err != nil {
			return err
		}
	}

	// 删除 comment
	err := lib.DB.Delete(comment).Error
	if err != nil {
		return err
	}

	return nil
}
