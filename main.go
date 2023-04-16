package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"

	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {

	proxySever()

	app := fiber.New()
	app.All("/*", func(ctx *fiber.Ctx) error {
		proxyUrl := fmt.Sprintf("https://%s/%s", os.Getenv("PROXY_DOMAIN"), ctx.Params("*", ""))
		if err := proxy.Do(ctx, proxyUrl); err != nil {
			return err
		}
		ctx.Response().Header.Del(fiber.HeaderServer)
		return nil
	})
	if err := app.Listen(":5333"); err != nil {
		fmt.Printf("启动失败: %s", err.Error())
		os.Exit(1)
	}
}

func proxySever() {
	// Google 目标 URL
	googleURL, _ := url.Parse("https://www.baidu.com")

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(googleURL)

	// 自定义修改请求头和请求主体
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)

		req.URL.Scheme = googleURL.Scheme
		req.URL.Host = googleURL.Host
		req.Host = googleURL.Host
		req.Header.Set("Referer", googleURL.String())
		req.Header.Set("Origin", googleURL.String())
	}

	// 改写响应头
	modifyResponse := func(resp *http.Response) error {
		resp.Header.Set("Access-Control-Allow-Origin", "*")
		return nil
	}
	proxy.ModifyResponse = modifyResponse

	// 监听本地端口，接收并转发请求
	err := http.ListenAndServe(":8080", proxy)
	if err != nil {
		panic(err)
	}
}
