package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const admin = "admin"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if u, p, ok := r.BasicAuth(); !ok || !(u == admin && p == admin) {
				http.Error(w, "unautorized", http.StatusUnauthorized)
				return
			}
			// r = r.WithContext(context.WithValue(r.Context(), 1, 0))
			next.ServeHTTP(w, r)
		},
	)
}

func GinAuthMW(c *gin.Context) {
	if u, p, ok := c.Request.BasicAuth(); !ok || !(u == admin && p == admin) {
		err := c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unautorized"))
		if err != nil {
			log.Println(err)
		}
		return
	}
	c.Next()
}
