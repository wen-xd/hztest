package main

//func main() {
//	h := server.Default(server.WithHostPorts("localhost:8083"))
//
//	h.Use(basic_auth.BasicAuthForRealm(map[string]string{
//		"username1": "password1",
//		"username2": "password2",
//	}, "xiao", "dong"))
//
//	h.GET("/basicAuth", func(ctx context.Context, c *app.RequestContext) {
//		v, e := c.Get("dong") // can set key by user
//		if e != true {
//			c.JSON(200, "hello")
//		}
//		c.JSON(consts.StatusOK, "hello hertz"+v.(string))
//	})
//
//	h.Spin()
//}

//func main() {
//	r := gin.Default()
//
//	r.Use(gin.BasicAuthForRealm(map[string]string{
//		"wen": "wen1",
//	}, "xiao"))
//	r.GET("/basicAuth", func(c *gin.Context) {
//		v, ok := c.Get("user") // 默认key was user
//		if !ok {
//			c.JSON(200, gin.H{"h": "world"})
//			return
//		}
//		c.JSON(200, gin.H{"hello": "world " + v.(string)})
//	})
//	r.Run(":8080")
//}
