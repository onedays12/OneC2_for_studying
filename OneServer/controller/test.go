package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Test 返回一个固定的 hello oneday 页面，用于 1.5 接口连通性验证
func (c *Controller) Test(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.String(http.StatusOK,
		`<!DOCTYPE html>  
<html>  
<head>  
    <title>Test</title></head>  
<body>  
    <h1>hello oneday</h1></body>  
</html>`)
}
