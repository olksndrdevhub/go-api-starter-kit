package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/olksndrdevhub/go-api-starter-kit/utils"
)

type customResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseSize int
}

type Middleware func(http.Handler) http.Handler

func CreateStuck(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := splitToken[1]

		claims, err := utils.ValidateJWTToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, utils.EmailKey, claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(utils.UserIDKey).(int64)
	return userID, ok
}

func GetEmailFromContext(r *http.Request) (string, bool) {
	email, ok := r.Context().Value(utils.EmailKey).(string)
	return email, ok
}

func (crw *customResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}
func (crw *customResponseWriter) Write(b []byte) (int, error) {
	size, err := crw.ResponseWriter.Write(b)
	crw.responseSize += size
	return size, err
}
func LogsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		crw := &customResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		startTime := time.Now()
		requestID := fmt.Sprintf("%d", time.Now().UnixNano())

		method := r.Method
		path := r.URL.Path
		clientIP := r.RemoteAddr
		userAgent := r.UserAgent()
		referer := r.Referer()

		next.ServeHTTP(crw, r)

		duration := time.Since(startTime)
		statusCode := crw.statusCode
		statusText := http.StatusText(statusCode)

		// Create structured log entry
		logEntry := fmt.Sprintf(
			"request_id=%s method=%s path=%s status=%d status_text=%s duration=%v ip=%s size=%d user_agent=%q referer=%q",
			requestID,
			method,
			path,
			statusCode,
			statusText,
			duration,
			clientIP,
			crw.responseSize,
			userAgent,
			referer,
		)

		if statusCode >= 500 {
			log.Printf("ERROR %s", logEntry)
		} else if statusCode >= 400 {
			log.Printf("WARN %s", logEntry)
		} else {
			log.Printf("INFO %s", logEntry)
		}
	})
}
