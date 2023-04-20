package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

type Url struct {
	LongUrl  string `json:"long_url"`
	ShortUrl string `json:"short_url"`
}

func main() {
	e := echo.New()

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		DB:   0,
	})

	// Work endpoints
	e.POST("/shorten", func(c echo.Context) error {
		ctx := c.Request().Context()

		var url Url
		if err := c.Bind(&url); err != nil {
			return err
		}

		// Search if already exist
		_, err := rdb.Get(ctx, url.LongUrl).Result()
		if err == redis.Nil {
			// Generate a hash
			h := sha256.New()
			h.Write([]byte(url.LongUrl))
			hash := hex.EncodeToString(h.Sum(nil))

			shortUrl := fmt.Sprintf(hash[:10])
			if err := rdb.Set(ctx, shortUrl, url.LongUrl, 24*time.Hour).Err(); err != nil {
				return err
			}
			url.ShortUrl = fmt.Sprintf("http://localhost:8080/%s", hash[:10])

		} else if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, url)
	})

	e.GET("/:short_url", func(c echo.Context) error {
		ctx := c.Request().Context()
		shortUrl := c.Param("short_url")

		// Health endpoints
		if shortUrl == "redis" {
			info, err := rdb.Info(ctx).Result()
			if err != nil {
				return c.String(http.StatusInternalServerError, "Redis is not available")
			}

			return c.String(http.StatusOK, info)
		} else if shortUrl == "health" {
			return c.String(http.StatusOK, "Shorten App is healthy!")
		}

		longUrl, err := rdb.Get(ctx, shortUrl).Result()
		if err == redis.Nil {
			//redirect to meli's 404 page
			return c.Redirect(http.StatusMovedPermanently, "https://www.mercadolibre.com.mx/dsagdfasgfdagdfgdfg")
		} else if err != nil {
			return err
		}

		return c.Redirect(http.StatusMovedPermanently, longUrl)
	})

	e.DELETE("/:short_url", func(c echo.Context) error {
		ctx := c.Request().Context()
		shortUrl := c.Param("short_url")

		num, err := rdb.Del(ctx, shortUrl).Result()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Redis is not available")
		}

		if num == 1 {
			return c.String(http.StatusOK, "Short url deleted")
		}

		return c.String(http.StatusOK, "URL not found")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
