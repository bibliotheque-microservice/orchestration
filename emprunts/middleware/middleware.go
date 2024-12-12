package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware(log *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {

		log := logrus.New()

		// Configuration du logger
		log.SetFormatter(&logrus.JSONFormatter{})
		log.SetLevel(logrus.InfoLevel)

		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		statusCode := c.Response().StatusCode()

		switch statusCode {
		case 200:
			log.WithFields(logrus.Fields{
				"method":   c.Method(),
				"path":     c.Path(),
				"status":   statusCode,
				"latency":  duration.String(),
				"ip":       c.IP(),
				"host":     c.Hostname(),
				"protocol": c.Protocol(),
				"query":    c.OriginalURL(),
			}).Info("Request completed")

		case 404:
			log.WithFields(logrus.Fields{
				"method":   c.Method(),
				"path":     c.Path(),
				"status":   statusCode,
				"latency":  duration.String(),
				"host":     c.Hostname(),
				"protocol": c.Protocol(),
				"query":    c.OriginalURL(),
			}).Warn("URL not found")

		case 400:
			log.WithFields(logrus.Fields{
				"method":   c.Method(),
				"path":     c.Path(),
				"status":   statusCode,
				"latency":  duration.String(),
				"ip":       c.IP(),
				"host":     c.Hostname(),
				"protocol": c.Protocol(),
				"query":    c.OriginalURL(),
			}).Warn("Bad request")

		case 500:
			log.WithFields(logrus.Fields{
				"method":   c.Method(),
				"path":     c.Path(),
				"status":   statusCode,
				"latency":  duration.String(),
				"ip":       c.IP(),
				"host":     c.Hostname(),
				"protocol": c.Protocol(),
				"query":    c.OriginalURL(),
			}).Error("Error")

		}

		return err
	}
}
