package api

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"time"
)

func (svc *APIService) LoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			req := ctx.Request()
			res := ctx.Response()

			var bodyBytes []byte

			log := svc.log.WithFields(logrus.Fields{
				"user_agent":  req.UserAgent(),
				"host":        req.Host,
				"uri":         req.RequestURI,
				"http_method": req.Method,
			})

			start := time.Now()
			if err = next(ctx); err != nil {
				ctx.Error(err)
			}

			stop := time.Now()
			if res.Status >= 400 {
				if len(bodyBytes) > 4096 {
					bodyBytes = bodyBytes[:4096]
				}
				log = log.WithFields(logrus.Fields{
					"body": string(bodyBytes),
				})
			}

			log = log.WithFields(logrus.Fields{
				"execution_time": stop.Sub(start).String(),
				"status":         res.Status,
			})

			if res.Status >= 400 {
				log.Infof("[error]: %v", err)
			} else {
				log.Info("[success]")
			}

			return nil
		}
	}
}
