package http

import (
	"github.com/ArtalkJS/ArtalkGo/model"
	"github.com/labstack/echo/v4"
)

type ParamsMarkRead struct {
	NotifyKey string `mapstructure:"notify_key"`

	SiteName string `mapstructure:"site_name"`
	SiteID   uint
	SiteAll  bool
}

func ActionMarkRead(c echo.Context) error {
	var p ParamsMarkRead
	if isOK, resp := ParamsDecode(c, ParamsMarkRead{}, &p); !isOK {
		return resp
	}

	// find site
	if isOK, resp := CheckSite(c, &p.SiteName, &p.SiteID, &p.SiteAll); !isOK {
		return resp
	}

	// find notify
	notify := model.FindNotifyByKey(p.NotifyKey)
	if notify.IsEmpty() {
		return RespError(c, "notify key is wrong")
	}

	if notify.IsRead {
		return RespSuccess(c)
	}

	// update notify
	err := notify.SetRead()
	if err != nil {
		return RespError(c, "notify save error")
	}

	return RespSuccess(c)
}
