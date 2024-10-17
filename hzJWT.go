package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/jwt"
)

type login struct {
	Username string `form:"username,required" json:"username,required"`
	Password string `form:"password,required" json:"password,required"`
}

var identityKey = "id"

func PingHandler(ctx context.Context, c *app.RequestContext) {
	user, _ := c.Get(identityKey)
	c.JSON(200, utils.H{
		"message": fmt.Sprintf("username:%v", user.(*User).UserName),
	})
}

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

func main() {
	h := server.Default()

	// the jwt middleware
	// use LoginHandler 	: Authenticator   -> PayloadFunc -> LoginResponse
	// use MiddlewareFunc 	: IdentityHandler -> Authorizator
	authMiddleware, err := jwt.New(&jwt.HertzJWTMiddleware{
		TokenHeadName: "wen",
		Realm:         "test zone",
		Key:           []byte("secret key"),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		IdentityKey:   identityKey,
		// 2.用于设置登录时为 token 添加自定义负载信息的函数，如果不传入这个参数，则 token 的 payload 部分默认存储 token 的过期时间和创建时间，如下则额外存储了用户名信息。
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		// IdentityHandler 作用在登录成功后的每次请求中，用于设置从 token 提取用户信息的函数。这里提到的用户信息在用户成功登录时，触发 PayloadFunc 函数，已经存入 token 的负载部分。
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return &User{
				UserName: claims[identityKey].(string),
			}
		},
		// 网页验证 配合 HertzJWTMiddleware.LoginHandler 使用，登录时触发，用于认证用户的登录信息
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var loginVals login
			if err := c.BindAndValidate(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
				return &User{
					UserName:  userID,
					LastName:  "Hertz",
					FirstName: "CloudWeGo",
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		// 用于设置已认证的用户路由访问权限的函数，如下函数通过验证用户名是否为 admin，从而判断是否有访问路由的权限。
		//如果没有访问权限，则会触发 Unauthorized 参数中声明的 jwt 流程验证失败的响应函数。
		Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
			if v, ok := data.(*User); ok && v.UserName == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(code, map[string]interface{}{
				"code":    code,
				"message": message,
			})
		},
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	h.POST("/login", authMiddleware.LoginHandler)

	h.NoRoute(authMiddleware.MiddlewareFunc(), func(ctx context.Context, c *app.RequestContext) {
		claims := jwt.ExtractClaims(ctx, c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, map[string]string{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := h.Group("/auth")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/ping", PingHandler)
	}

	h.Spin()
}
