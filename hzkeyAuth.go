package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/keyauth"
)

func main() {
	h := server.Default()
	h.Use(keyauth.New(
		keyauth.WithValidator(func(ctx context.Context, requestContext *app.RequestContext, s string) (bool, error) {
			// s 指的是验证字段 Authorization 默认的 scheme为Bearer 可以通过withkeylookup 统一修改
			// <Authorization>：<Bearer> <token>
			if s == "test_admin" {
				return true, nil
			}
			return false, nil
		}),
	))
	h.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		value, _ := c.Get("token")
		c.JSON(consts.StatusOK, utils.H{"ping": value})
	})
	h.Spin()
}
