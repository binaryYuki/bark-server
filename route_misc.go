package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

func init() {
	registerRoute("misc", func(router fiber.Router) {
		router.Get("/ping", func(c *fiber.Ctx) error {
			return c.JSON(CommonResp{
				Code:      200,
				Message:   "pong",
				Timestamp: time.Now().Unix(),
			})
		})

		router.Get("/healthz", func(c *fiber.Ctx) error {
			resp, err := http.Get("https://ipinfo.io")
			if err != nil {
				return c.Status(500).SendString("Health check failed")
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return c.Status(500).SendString("Failed to read response body")
			}

			// 创建一个映射来接收 JSON 数据
			ipinfo := make(map[string]interface{})
			err = json.Unmarshal(body, &ipinfo)
			if err != nil {
				return c.Status(500).SendString("Failed to parse JSON")
			}

			// 删除 readme 字段
			delete(ipinfo, "readme")

			return c.JSON(map[string]interface{}{
				"status": "ok",
				"health": "passed",
				"ipinfo": ipinfo, // 直接使用映射作为值
			})
		})

		router.Get("/info", func(c *fiber.Ctx) error {
			devices, _ := db.CountAll()
			return c.JSON(map[string]interface{}{
				"version": version,
				"build":   buildDate,
				"author":  "binaryYuki <noreply.tzpro.xyz>",
				"arch":    runtime.GOOS + "/" + runtime.GOARCH,
				"commit":  commitID,
				"devices": devices,
			})
		})
	})
}
