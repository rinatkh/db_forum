package api

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/rinatkh/db_forum/internal/constants"
	"github.com/rinatkh/db_forum/internal/model/core"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func (svc *APIService) XRequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			xRequestID := ctx.Request().Header.Get(constants.HeaderKeyRequestID)
			if len(xRequestID) == 0 {
				xRequestID, err := core.GenUUID()
				if err != nil {
					return err
				}
				ctx.Request().Header.Set(constants.HeaderKeyRequestID, xRequestID)
			}
			return next(ctx)
		}
	}
}

func (svc *APIService) LoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			req := ctx.Request()
			res := ctx.Response()

			var bodyBytes []byte
			if svc.debug {
				bodyBytes, err = ioutil.ReadAll(req.Body)
				if err != nil {
					ctx.Error(err)
				}
				err := req.Body.Close()
				if err != nil {
					return err
				}
				req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			log := svc.log.WithFields(logrus.Fields{
				"x_request_id": ctx.Request().Header.Get(constants.HeaderKeyRequestID),
				"user_agent":   req.UserAgent(),
				"host":         req.Host,
				"uri":          req.RequestURI,
				"http_method":  req.Method,
				"user_ip":      getIP(req),
			})

			userID := ctx.Request().Header.Get(constants.HeaderKeyUserID)
			if len(userID) != 0 {
				log = log.WithFields(logrus.Fields{
					"user_id": userID,
				})
			}

			start := time.Now()
			if err = next(ctx); err != nil {
				ctx.Error(err)
			}

			stop := time.Now()
			if res.Status >= 400 && svc.debug {
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

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	for _, ip := range strings.Split(r.Header.Get("X-FORWARDED-FOR"), ",") {
		if netIP := net.ParseIP(ip); netIP != nil {
			return ip
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	return ""
}
