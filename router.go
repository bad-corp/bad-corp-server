package main

import (
	"github.com/labstack/echo/v4"
)

func router(e *echo.Echo, ctl *controller) {
	e.POST("/auth/sign_up", ctl.signUp)
	e.POST("/auth/sign_in", ctl.signIn)
	e.GET("/user", ctl.getUser, jwtMiddleware(), userIdMiddleware)
	e.POST("/subjects", ctl.createSubject, jwtMiddleware(), userIdMiddleware)
	e.GET("/subjects/random", ctl.getSubjectByRandom, jwtMiddleware(), userIdMiddleware)
	e.POST("/subjects/:subjectId/comments", ctl.createSubjectComment, jwtMiddleware(), userIdMiddleware)
	e.GET("/trending/subjects", ctl.getTrendingSubjects) // TODO: 暂未实现
	e.GET("/trending/kings", ctl.getTrendingKings)
	e.GET("/trending/queens", ctl.getTrendingQueens)
}
