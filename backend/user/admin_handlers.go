package user

import (
	"bamort/database"
	"bamort/logger"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListUsers returns all users (admin only)
func ListUsers(c *gin.Context) {
	logger.Debug("Listing all users...")

	var users []User
	if err := database.DB.Find(&users).Error; err != nil {
		logger.Error("Failed to fetch users: %s", err.Error())
		respondWithError(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	// Remove password hashes from response
	for i := range users {
		users[i].PasswordHash = ""
		users[i].ResetPwHash = nil
		users[i].DisplayName = users[i].DisplayNameOrUsername()
	}

	logger.Info("Successfully fetched %d users", len(users))
	c.JSON(http.StatusOK, users)
}

// GetUser returns a specific user by ID (admin only, or own profile)
func GetUser(c *gin.Context) {
	logger.Debug("Fetching user by ID...")

	userIDParam := c.Param("id")
	targetUserID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		logger.Error("Invalid user ID: %s", userIDParam)
		respondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get requesting user from context
	requestingUserInterface, exists := c.Get("user")
	if !exists {
		logger.Error("User not found in context")
		respondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	requestingUser, ok := requestingUserInterface.(*User)
	if !ok {
		logger.Error("Invalid user context")
		respondWithError(c, http.StatusInternalServerError, "Invalid user context")
		return
	}

	// Allow users to view their own profile, or admins to view any profile
	if requestingUser.UserID != uint(targetUserID) && !requestingUser.IsAdmin() {
		logger.Warn("User %s attempted to access user %d without permission", requestingUser.Username, targetUserID)
		respondWithError(c, http.StatusForbidden, "Forbidden")
		return
	}

	var user User
	if err := user.FirstId(uint(targetUserID)); err != nil {
		logger.Error("User not found: %d", targetUserID)
		respondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	// Remove sensitive data
	user.PasswordHash = ""
	user.ResetPwHash = nil
	user.DisplayName = user.DisplayNameOrUsername()

	logger.Info("Successfully fetched user: %s (ID: %d)", user.Username, user.UserID)
	c.JSON(http.StatusOK, user)
}

// UpdateUserRole updates a user's role (admin only)
func UpdateUserRole(c *gin.Context) {
	logger.Debug("Updating user role...")

	userIDParam := c.Param("id")
	targetUserID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		logger.Error("Invalid user ID: %s", userIDParam)
		respondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var input struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Failed to parse role update data: %s", err.Error())
		respondWithError(c, http.StatusBadRequest, "Role is required")
		return
	}

	// Validate role
	if !ValidateRole(input.Role) {
		logger.Error("Invalid role: %s", input.Role)
		respondWithError(c, http.StatusBadRequest, fmt.Sprintf("Invalid role. Must be one of: %s, %s, %s", RoleStandardUser, RoleMaintainer, RoleAdmin))
		return
	}

	var user User
	if err := user.FirstId(uint(targetUserID)); err != nil {
		logger.Error("User not found: %d", targetUserID)
		respondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	// Get requesting user for logging
	requestingUserInterface, _ := c.Get("user")
	requestingUser, _ := requestingUserInterface.(*User)

	oldRole := user.Role
	user.Role = input.Role

	if err := user.Save(); err != nil {
		logger.Error("Failed to update user role for user %s: %s", user.Username, err.Error())
		respondWithError(c, http.StatusInternalServerError, "Failed to update user role")
		return
	}

	logger.Info("User role updated: %s (ID: %d) from %s to %s by %s", user.Username, user.UserID, oldRole, user.Role, requestingUser.Username)
	c.JSON(http.StatusOK, gin.H{
		"message": "User role updated successfully",
		"user": gin.H{
			"id":           user.UserID,
			"username":     user.Username,
			"display_name": user.DisplayNameOrUsername(),
			"role":         user.Role,
		},
	})
}

// DeleteUser deletes a user (admin only)
func DeleteUser(c *gin.Context) {
	logger.Debug("Deleting user...")

	userIDParam := c.Param("id")
	targetUserID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		logger.Error("Invalid user ID: %s", userIDParam)
		respondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get requesting user
	requestingUserInterface, exists := c.Get("user")
	if !exists {
		logger.Error("User not found in context")
		respondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	requestingUser, ok := requestingUserInterface.(*User)
	if !ok {
		logger.Error("Invalid user context")
		respondWithError(c, http.StatusInternalServerError, "Invalid user context")
		return
	}

	// Prevent self-deletion
	if requestingUser.UserID == uint(targetUserID) {
		logger.Warn("User %s attempted to delete themselves", requestingUser.Username)
		respondWithError(c, http.StatusBadRequest, "Cannot delete your own account")
		return
	}

	var user User
	if err := user.FirstId(uint(targetUserID)); err != nil {
		logger.Error("User not found: %d", targetUserID)
		respondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		logger.Error("Failed to delete user %s: %s", user.Username, err.Error())
		respondWithError(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	logger.Info("User deleted: %s (ID: %d) by %s", user.Username, user.UserID, requestingUser.Username)
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// ChangeUserPassword allows admin to change a user's password (admin only)
func ChangeUserPassword(c *gin.Context) {
	logger.Debug("Admin changing user password...")

	userIDParam := c.Param("id")
	targetUserID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		logger.Error("Invalid user ID: %s", userIDParam)
		respondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var input struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Failed to parse password data: %s", err.Error())
		respondWithError(c, http.StatusBadRequest, "New password (min 6 characters) is required")
		return
	}

	var user User
	if err := user.FirstId(uint(targetUserID)); err != nil {
		logger.Error("User not found: %d", targetUserID)
		respondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	// Get requesting user for logging
	requestingUserInterface, _ := c.Get("user")
	requestingUser, _ := requestingUserInterface.(*User)

	// Hash new password using MD5 (same as registration)
	hashedPassword := md5.Sum([]byte(input.NewPassword))
	user.PasswordHash = hex.EncodeToString(hashedPassword[:])

	if err := user.Save(); err != nil {
		logger.Error("Failed to update password for user %s: %s", user.Username, err.Error())
		respondWithError(c, http.StatusInternalServerError, "Failed to update password")
		return
	}

	logger.Info("Password changed for user %s (ID: %d) by admin %s", user.Username, user.UserID, requestingUser.Username)
	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}
