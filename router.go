package main

import (
	"github.com/labstack/echo/v4"
)

func router(e *echo.Echo, ctl *controller) {
	e.POST("/auth/sign_up", ctl.signUp)
	e.POST("/auth/sign_in", ctl.signIn)
	e.GET("/user", ctl.getUser, jwtMiddleware(), userIdMiddleware)
	e.GET("/corps/search", ctl.searchCorp, jwtMiddleware(), userIdMiddleware)
	e.GET("/corps/:corp_id", ctl.getCorp, jwtMiddleware(), userIdMiddleware)
	e.POST("/corps/:corp_id/comments", ctl.createCorpComment, jwtMiddleware(), userIdMiddleware)
	e.GET("/trending/corps", ctl.getTrendingCorps, jwtMiddleware(), userIdMiddleware)
}
