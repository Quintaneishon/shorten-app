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

	// middleware to get simple stats
	middle := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			endpoint := c.Path()
			rdb.Incr(c.Request().Context(), "stats:"+endpoint)
			return next(c)
		}
	}

	e.Use(middle)

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
			url.ShortUrl = fmt.Sprintf("http://localhost/%s", hash[:10])

		} else if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, url)
	})

	e.GET("/:short_url", func(c echo.Context) error {
		ctx := c.Request().Context()
		shortUrl := c.Param("short_url")

		switch shortUrl {
		// Health endpoints
		case "redis":
			info, err := rdb.Info(ctx).Result()
			if err != nil {
				return c.String(http.StatusInternalServerError, "Redis is not available")
			}

			return c.String(http.StatusOK, info)
		case "health":
			return c.String(http.StatusOK, "Shorten App is healthy!")
		case "stats":
			keys, err := rdb.Keys(ctx, "stats:*").Result()
			if err != nil {
				return err
			}

			stats := make(map[string]int64)

			for _, key := range keys {
				count, err := rdb.Get(ctx, key).Int64()
				if err != nil {
					return err
				}

				endpoint := key[len("stats:"):]
				stats[endpoint] = count
			}

			return c.JSON(http.StatusOK, stats)
		default:

			longUrl, err := rdb.Get(ctx, shortUrl).Result()
			if err == redis.Nil {
				//redirect to 404 page
				return c.File("404.html")
			} else if err != nil {
				return err
			}

			return c.Redirect(http.StatusMovedPermanently, longUrl)
		}
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
