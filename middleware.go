package main

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func jwtMiddleware() echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte("zhuzhu"),
	}
	return echojwt.WithConfig(config)
}

func userIdMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authUser := c.Get("user").(*jwt.Token)
		claims := authUser.Claims.(*jwtCustomClaims)
		c.Set("userId", claims.UserId)
		return next(c)
	}
}
