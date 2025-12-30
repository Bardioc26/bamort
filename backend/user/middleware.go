package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole is a middleware that checks if the user has the required role
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by auth middleware)
		userInterface, exists := c.Get("user")
		if !exists {
			respondWithError(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		user, ok := userInterface.(*User)
		if !ok {
			respondWithError(c, http.StatusInternalServerError, "Invalid user context")
			c.Abort()
			return
		}

		// Check if user has required role
		switch requiredRole {
		case RoleAdmin:
			if !user.IsAdmin() {
				respondWithError(c, http.StatusForbidden, "Admin role required")
				c.Abort()
				return
			}
		case RoleMaintainer:
			if !user.IsMaintainer() {
				respondWithError(c, http.StatusForbidden, "Maintainer role required")
				c.Abort()
				return
			}
		case RoleStandardUser:
			if !user.IsStandardUser() {
				respondWithError(c, http.StatusForbidden, "Insufficient permissions")
				c.Abort()
				return
			}
		default:
			respondWithError(c, http.StatusInternalServerError, "Invalid role requirement")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin is a convenience middleware for admin-only endpoints
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(RoleAdmin)
}

// RequireMaintainer is a convenience middleware for maintainer-or-higher endpoints
func RequireMaintainer() gin.HandlerFunc {
	return RequireRole(RoleMaintainer)
}
