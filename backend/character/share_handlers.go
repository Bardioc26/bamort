package character

import (
	"bamort/database"
	"bamort/models"
	"bamort/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCharacterShares returns the list of users a character is shared with
// GetCharacterShares godoc
// @Summary Get character shares
// @Description Returns list of users with whom the character is shared
// @Tags Characters
// @Produce json
// @Param id path int true "Character ID"
// @Success 200 {array} object "List of user shares"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Failure 404 {object} map[string]string "Character not found"
// @Security BearerAuth
// @Router /api/characters/{id}/shares [get]
func GetCharacterShares(c *gin.Context) {
	charID := c.Param("id")

	var character models.Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Character not found")
		return
	}

	// Check ownership
	if !checkCharacterOwnership(c, &character) {
		return
	}

	var shares []models.CharShare
	if err := database.DB.Where("character_id = ?", character.ID).Find(&shares).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve shares")
		return
	}

	// Get user details for each share
	type ShareWithUser struct {
		models.CharShare
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
	}

	var sharesWithUsers []ShareWithUser
	for _, share := range shares {
		var u user.User
		if err := u.FirstId(share.UserID); err == nil {
			sharesWithUsers = append(sharesWithUsers, ShareWithUser{
				CharShare:   share,
				Username:    u.Username,
				DisplayName: u.DisplayNameOrUsername(),
				Email:       u.Email,
			})
		}
	}

	c.JSON(http.StatusOK, sharesWithUsers)
}

// UpdateCharacterShares updates the list of users a character is shared with
// UpdateCharacterShares godoc
// @Summary Update character shares
// @Description Updates the list of users with whom the character is shared (owner only)
// @Tags Characters
// @Accept json
// @Produce json
// @Param id path int true "Character ID"
// @Param shares body object{user_ids=[]int} true "List of user IDs to share with"
// @Success 200 {object} map[string]string "Shares updated successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied - owner only"
// @Security BearerAuth
// @Router /api/characters/{id}/shares [put]
func UpdateCharacterShares(c *gin.Context) {
	charID := c.Param("id")

	var character models.Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Character not found")
		return
	}

	// Check ownership
	if !checkCharacterOwnership(c, &character) {
		return
	}

	type UpdateSharesRequest struct {
		UserIDs []uint `json:"user_ids" binding:"required"`
	}

	var request UpdateSharesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Delete existing shares
	if err := database.DB.Where("character_id = ?", character.ID).Delete(&models.CharShare{}).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to delete existing shares")
		return
	}

	// Create new shares
	for _, userID := range request.UserIDs {
		// Don't share with yourself
		if userID == character.UserID {
			continue
		}

		share := models.CharShare{
			CharacterID: character.ID,
			UserID:      userID,
			Permission:  "read",
		}

		if err := database.DB.Create(&share).Error; err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to create share")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shares updated successfully"})
}

// GetAvailableUsersForSharing returns a list of users (excluding the owner)
// GetAvailableUsersForSharing godoc
// @Summary Get available users for sharing
// @Description Returns list of users available to share character with (excludes owner and current shares)
// @Tags Characters
// @Produce json
// @Param id path int true "Character ID"
// @Success 200 {array} object "List of available users"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Security BearerAuth
// @Router /api/characters/{id}/available-users [get]
func GetAvailableUsersForSharing(c *gin.Context) {
	charID := c.Param("id")

	var character models.Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Character not found")
		return
	}

	// Check ownership
	if !checkCharacterOwnership(c, &character) {
		return
	}

	var users []user.User
	if err := database.DB.Where("user_id != ?", character.UserID).Find(&users).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	// Remove sensitive data
	type UserInfo struct {
		UserID      uint   `json:"user_id"`
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
	}

	var userInfos []UserInfo
	for _, u := range users {
		userInfos = append(userInfos, UserInfo{
			UserID:      u.UserID,
			Username:    u.Username,
			DisplayName: u.DisplayNameOrUsername(),
			Email:       u.Email,
		})
	}

	c.JSON(http.StatusOK, userInfos)
}
