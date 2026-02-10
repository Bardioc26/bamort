/*
User Handlers

Add handlers for user registration and login:
*/
package user

import (
	"bamort/logger"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

func respondWithError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Creates a new user account with username, email, and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body object{username=string,email=string,password_hash=string} true "User registration data"
// @Success 201 {object} map[string]string "User registered successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 500 {object} map[string]string "Failed to create user"
// @Router /register [post]
func RegisterUser(c *gin.Context) {
	logger.Debug("Starte Benutzerregistrierung...")

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Error("Fehler beim Parsen der Registrierungsdaten: %s", err.Error())
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate that username is not empty
	if user.Username == "" {
		logger.Error("Registrierung fehlgeschlagen - Benutzername ist leer")
		respondWithError(c, http.StatusBadRequest, "Username cannot be empty")
		return
	}

	// Validate that email is not empty
	if user.Email == "" {
		logger.Error("Registrierung fehlgeschlagen - E-Mail ist leer")
		respondWithError(c, http.StatusBadRequest, "Email cannot be empty")
		return
	}

	// Validate that password is not empty
	if user.PasswordHash == "" {
		logger.Error("Registrierung fehlgeschlagen - Passwort ist leer")
		respondWithError(c, http.StatusBadRequest, "Password cannot be empty")
		return
	}

	logger.Debug("Registriere Benutzer: %s", user.Username)
	//fmt.Printf("User input: '%s'", user.PasswordHash)
	//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	hashedPassword := md5.Sum([]byte(user.PasswordHash))
	user.PasswordHash = hex.EncodeToString(hashedPassword[:])
	logger.Debug("Passwort-Hash erstellt für Benutzer: %s", user.Username)

	// Set default role for new users
	if user.Role == "" {
		user.Role = RoleStandardUser
	}

	//fmt.Printf("pwdh: %s", user.PasswordHash)
	if err := user.Create(); err != nil {
		logger.Error("Fehler beim Erstellen des Benutzers %s: %s", user.Username, err.Error())
		respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %s", err))
		return
	}

	logger.Info("Benutzer erfolgreich registriert: %s (ID: %d)", user.Username, user.UserID)
	//fmt.Printf(" ___ pwdh2: %s", user.PasswordHash)
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully:"})
}

// GenerateToken creates a JWT token for the given user
func GenerateToken(u *User) string {
	logger.Debug("Generiere Token für Benutzer: %s (ID: %d)", u.Username, u.UserID)

	//u.Username + "lkiuztrew" + u.CreatedAt.String()
	tx := md5.Sum([]byte(u.Username + u.CreatedAt.String()))
	// Convert hash to raw string
	hashString := hex.EncodeToString(tx[:])
	pos := 7
	idm := "." + fmt.Sprintf("%d", u.UserID) + ":"
	// Insert the character
	token := hashString[:pos] + string(idm) + hashString[pos:]

	logger.Debug("Token erfolgreich generiert für Benutzer: %s", u.Username)
	return token
}
func CheckToken(token string) *User {
	logger.Debug("Prüfe Token-Gültigkeit...")

	//fmt.Print("CheckToken1: " + token)
	var u User
	var err error
	pos := 7 + len("Bearer ")
	userid := 0
	// Check if a `.` is at position 7 (zero-indexed)
	if len(token) > pos && token[pos] == '.' {
		logger.Debug("Token-Format erkannt, extrahiere Benutzer-ID...")
		//fmt.Print("CheckToken2: " + token + "\n")
		// Find the next `:` after the `.`
		colonPos := strings.Index(token[pos+1:], ":") // Start searching after position 7
		if colonPos != -1 {
			//fmt.Printf("CheckToken3: %v\n", colonPos)
			// Extract the substring between `.` and `:`
			uu := token[pos+1 : pos+1+colonPos]
			//fmt.Println("Extracted Substring:" + uu + "\n")
			userid, err = strconv.Atoi(uu)
			//fmt.Printf("Extracted UserID: %v \n", userid)
			if err != nil {
				logger.Error("Fehler beim Parsen der Benutzer-ID aus Token: %s", err.Error())
				//fmt.Print("CheckToken4: " + err.Error() + "\n")
				return nil
			}
			logger.Debug("Benutzer-ID aus Token extrahiert: %d", userid)
		} else {
			logger.Debug("Token-Format ungültig: Kein ':' nach '.' gefunden")
			//fmt.Print("CheckToken5: not found\n")
			return nil
		}
	} else {
		logger.Debug("Token-Format ungültig: Kein '.' an erwarteter Position")
		//fmt.Print("CheckToken6: not found\n")
		return nil
	}

	if userid > 0 {
		logger.Debug("Lade Benutzer mit ID: %d", userid)
		//fmt.Printf("CheckToken6-1: userid %v\n", userid)
		//fmt.Printf("CheckToken6-1: userid %v\n", uint(userid))
		err := u.FirstId(uint(userid))
		if err != nil {
			logger.Error("Benutzer mit ID %d nicht gefunden: %s", userid, err.Error())
			//fmt.Printf("CheckToken7: not found error %s\n", err.Error())
			return nil
		}
		logger.Debug("Benutzer gefunden und Token validiert: %s (ID: %d)", u.Username, u.UserID)
		//fmt.Printf("CheckToken8: found:%s \n", u.Username)
		return &u
	}
	logger.Debug("Token-Validierung fehlgeschlagen: Ungültige Benutzer-ID")
	//fmt.Print("CheckToken9: not found\n")
	return nil
}

// LoginUser godoc
// @Summary User login
// @Description Authenticates a user and returns a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body object{username=string,password=string} true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful with token and user data"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /login [post]
func LoginUser(c *gin.Context) {
	logger.Debug("Starte Benutzer-Anmeldung...")

	var user User
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Fehler beim Parsen der Login-Daten: %s", err.Error())
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debug("Login-Versuch für Benutzer: %s", input.Username)

	//if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
	if err := user.First(input.Username); err != nil {
		logger.Warn("Login fehlgeschlagen - Benutzer nicht gefunden: %s", input.Username)
		respondWithError(c, http.StatusUnauthorized, fmt.Sprintf("Invalid username. or password %v", input))
		return
	}

	logger.Debug("Benutzer gefunden, prüfe Passwort für: %s", input.Username)
	hashedPassword := md5.Sum([]byte(input.Password))
	fmt.Printf("pwdh: %s", hex.EncodeToString(hashedPassword[:]))
	if user.PasswordHash != hex.EncodeToString(hashedPassword[:]) {
		logger.Warn("Login fehlgeschlagen - Ungültiges Passwort für Benutzer: %s", input.Username)
		respondWithError(c, http.StatusUnauthorized, fmt.Sprintf("Invalid username. or password. %s %s", input.Password, hex.EncodeToString(hashedPassword[:])))
		return
	}
	/*
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
			respondWithError(c, http.StatusUnauthorized, "Invalid username or password.")
			return
		}
	*/

	logger.Info("Login erfolgreich für Benutzer: %s (ID: %d)", user.Username, user.UserID)
	token := GenerateToken(&user)
	logger.Debug("Login-Token generiert für Benutzer: %s", user.Username)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

// Apply middleware to protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug("Prüfe Authentifizierung für Request: %s %s", c.Request.Method, c.Request.URL.Path)

		token := c.GetHeader("Authorization")
		if token == "" {
			logger.Warn("Authentifizierung fehlgeschlagen - Kein Authorization-Header für %s %s", c.Request.Method, c.Request.URL.Path)
			respondWithError(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		logger.Debug("Authorization-Header gefunden, prüfe Token...")
		user := CheckToken(token)
		if user == nil {
			logger.Warn("Authentifizierung fehlgeschlagen - Ungültiger Token für %s %s", c.Request.Method, c.Request.URL.Path)
			respondWithError(c, http.StatusUnauthorized, "Unauthorized.")
			c.Abort()
			return
		}

		logger.Debug("Authentifizierung erfolgreich für Benutzer: %s (ID: %d) - %s %s", user.Username, user.UserID, c.Request.Method, c.Request.URL.Path)

		// Set user information in context
		c.Set("userID", user.UserID)
		c.Set("username", user.Username)
		c.Set("user", user)

		c.Next()
	}
}

// generateResetHash generiert einen sicheren Hash für Password-Reset
func generateResetHash() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// sendResetEmail simuliert das Senden einer E-Mail (hier nur Logging)
// In einer echten Implementierung würde hier ein E-Mail-Service verwendet
func sendResetEmail(email, username, resetHash, frontendURL string) error {
	// Verwende die mitgegebene Frontend-URL oder fallback auf Standard
	baseURL := frontendURL
	if baseURL == "" {
		baseURL = "http://localhost:3000" // Fallback, sollte aber nicht verwendet werden
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", baseURL, resetHash)

	logger.Info("=== PASSWORD RESET EMAIL ===")
	logger.Info("An: %s", email)
	logger.Info("Betreff: Passwort zurücksetzen für %s", username)
	logger.Info("Nachricht:")
	logger.Info("Hallo %s,", username)
	logger.Info("")
	logger.Info("Sie haben eine Passwort-Zurücksetzung angefordert.")
	logger.Info("Klicken Sie auf den folgenden Link, um Ihr Passwort zurückzusetzen:")
	logger.Info("")
	logger.Info("%s", resetLink)
	logger.Info("")
	logger.Info("Dieser Link ist 14 Tage gültig.")
	logger.Info("Falls Sie diese Anfrage nicht gestellt haben, ignorieren Sie diese E-Mail.")
	logger.Info("")
	logger.Info("=== END EMAIL ===")

	// TODO: Hier echte E-Mail-Integration hinzufügen
	// z.B. SendGrid, SMTP, etc.

	return nil
}

// RequestPasswordReset Handler für Passwort-Reset-Anfrage
// RequestPasswordReset godoc
// @Summary Request password reset
// @Description Sends a password reset email to the user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param email body object{email=string} true "User email address"
// @Success 200 {object} map[string]string "Password reset email sent"
// @Failure 400 {object} map[string]string "Invalid email format"
// @Router /password-reset/request [post]
func RequestPasswordReset(c *gin.Context) {
	logger.Debug("Starte Passwort-Reset-Anfrage...")

	var input struct {
		Email       string `json:"email" binding:"required,email"`
		RedirectURL string `json:"redirect_url,omitempty"` // Optionale Frontend-URL
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Fehler beim Parsen der Reset-Anfrage: %s", err.Error())
		respondWithError(c, http.StatusBadRequest, "Gültige E-Mail-Adresse erforderlich")
		return
	}

	// Frontend-URL aus Request verwenden
	redirectURL := input.RedirectURL
	if redirectURL == "" {
		// Fallback, sollte aber nicht verwendet werden, da Frontend die URL mitgeben sollte
		redirectURL = "http://localhost:3000"
	}

	logger.Debug("Reset-Anfrage für E-Mail: %s", input.Email)

	var user User
	if err := user.FindByEmail(input.Email); err != nil {
		// Aus Sicherheitsgründen keine Information preisgeben, ob die E-Mail existiert
		logger.Warn("Reset-Anfrage für nicht existierende E-Mail: %s", input.Email)
		c.JSON(http.StatusOK, gin.H{
			"message": "Falls ein Account mit dieser E-Mail-Adresse existiert, wurde eine Reset-E-Mail gesendet.",
		})
		return
	}

	// Generiere Reset-Hash
	resetHash, err := generateResetHash()
	if err != nil {
		logger.Error("Fehler beim Generieren des Reset-Hashes: %s", err.Error())
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Verarbeiten der Anfrage")
		return
	}

	// Speichere Reset-Hash in der Datenbank
	if err := user.SetPasswordResetHash(resetHash); err != nil {
		logger.Error("Fehler beim Speichern des Reset-Hashes für Benutzer %s: %s", user.Username, err.Error())
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Verarbeiten der Anfrage")
		return
	}

	// Sende Reset-E-Mail
	if err := sendResetEmail(user.Email, user.Username, resetHash, redirectURL); err != nil {
		logger.Error("Fehler beim Senden der Reset-E-Mail für Benutzer %s: %s", user.Username, err.Error())
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Senden der E-Mail")
		return
	}

	logger.Info("Reset-E-Mail erfolgreich für Benutzer %s (%s) gesendet", user.Username, user.Email)
	c.JSON(http.StatusOK, gin.H{
		"message": "Falls ein Account mit dieser E-Mail-Adresse existiert, wurde eine Reset-E-Mail gesendet.",
	})
}

// ResetPassword Handler für das Zurücksetzen des Passworts
// ResetPassword godoc
// @Summary Reset password
// @Description Resets user password using a valid reset token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param reset_data body object{token=string,new_password=string} true "Reset token and new password"
// @Success 200 {object} map[string]string "Password reset successful"
// @Failure 400 {object} map[string]string "Invalid token or password"
// @Router /password-reset/reset [post]
func ResetPassword(c *gin.Context) {
	logger.Debug("Starte Passwort-Reset...")

	var input struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Fehler beim Parsen der Reset-Daten: %s", err.Error())
		respondWithError(c, http.StatusBadRequest, "Token und neues Passwort (mind. 6 Zeichen) erforderlich")
		return
	}

	logger.Debug("Reset-Versuch mit Token: %s", input.Token[:10]+"...")

	var user User
	if err := user.FindByResetHash(input.Token); err != nil {
		logger.Warn("Ungültiger oder abgelaufener Reset-Token verwendet")
		respondWithError(c, http.StatusBadRequest, "Ungültiger oder abgelaufener Reset-Link")
		return
	}

	// Zusätzliche Validierung des Tokens
	if !user.IsResetHashValid(input.Token) {
		logger.Warn("Reset-Token-Validierung fehlgeschlagen für Benutzer: %s", user.Username)
		respondWithError(c, http.StatusBadRequest, "Ungültiger oder abgelaufener Reset-Link")
		return
	}

	// Neues Passwort hashen (gleiche Methode wie bei der Registrierung)
	hashedPassword := md5.Sum([]byte(input.NewPassword))
	user.PasswordHash = hex.EncodeToString(hashedPassword[:])

	// Reset-Hash entfernen
	if err := user.ClearPasswordResetHash(); err != nil {
		logger.Error("Fehler beim Entfernen des Reset-Hashes für Benutzer %s: %s", user.Username, err.Error())
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Aktualisieren des Accounts")
		return
	}

	// Passwort speichern
	if err := user.Save(); err != nil {
		logger.Error("Fehler beim Speichern des neuen Passworts für Benutzer %s: %s", user.Username, err.Error())
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Aktualisieren des Passworts")
		return
	}

	logger.Info("Passwort erfolgreich zurückgesetzt für Benutzer: %s", user.Username)
	c.JSON(http.StatusOK, gin.H{
		"message": "Passwort erfolgreich zurückgesetzt",
	})
}

// ValidateResetToken Handler zur Validierung eines Reset-Tokens
// ValidateResetToken godoc
// @Summary Validate reset token
// @Description Validates a password reset token
// @Tags Authentication
// @Produce json
// @Param token path string true "Reset token"
// @Success 200 {object} map[string]string "Token is valid"
// @Failure 400 {object} map[string]string "Invalid or expired token"
// @Router /password-reset/validate/{token} [get]
func ValidateResetToken(c *gin.Context) {
	logger.Debug("Validiere Reset-Token...")

	token := c.Param("token")
	//token := c.Query("token")
	if token == "" {
		respondWithError(c, http.StatusBadRequest, "Token erforderlich")
		return
	}

	var user User
	if err := user.FindByResetHash(token); err != nil {
		logger.Debug("Reset-Token nicht gefunden oder abgelaufen")
		respondWithError(c, http.StatusBadRequest, "Ungültiger oder abgelaufener Reset-Link")
		return
	}

	if !user.IsResetHashValid(token) {
		logger.Debug("Reset-Token-Validierung fehlgeschlagen")
		respondWithError(c, http.StatusBadRequest, "Ungültiger oder abgelaufener Reset-Link")
		return
	}

	logger.Debug("Reset-Token gültig für Benutzer: %s", user.Username)
	c.JSON(http.StatusOK, gin.H{
		"valid":        true,
		"username":     user.Username,
		"display_name": user.DisplayNameOrUsername(),
		"expires":      user.ResetPwHashExpires,
	})
}

// GetUserProfile Handler to get current user's profile information
// GetUserProfile godoc
// @Summary Get current user profile
// @Description Returns the profile information of the authenticated user
// @Tags User
// @Produce json
// @Success 200 {object} User "User profile data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "User not found"
// @Security BearerAuth
// @Router /api/user/profile [get]
func GetUserProfile(c *gin.Context) {
	logger.Debug("Lade Benutzerprofil...")

	// Get user ID from context (set by AuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		logger.Error("Benutzer-ID nicht im Context gefunden")
		respondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var user User
	if err := user.FirstId(userID.(uint)); err != nil {
		logger.Error("Benutzer mit ID %v nicht gefunden: %s", userID, err.Error())
		respondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	logger.Debug("Benutzerprofil geladen für: %s (ID: %d)", user.Username, user.UserID)
	c.JSON(http.StatusOK, gin.H{
		"id":                 user.UserID,
		"username":           user.Username,
		"display_name":       user.DisplayNameOrUsername(),
		"email":              user.Email,
		"role":               user.Role,
		"preferred_language": user.PreferredLanguage,
	})
}

// UpdateEmail Handler to update user's email address
// UpdateEmail godoc
// @Summary Update user email
// @Description Updates the email address of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Param email body object{email=string} true "New email address"
// @Success 200 {object} map[string]string "Email updated successfully"
// @Failure 400 {object} map[string]string "Invalid email format"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/user/email [put]
func UpdateEmail(c *gin.Context) {
	logger.Debug("Starte E-Mail-Aktualisierung...")

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		logger.Error("Benutzer-ID nicht im Context gefunden")
		respondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Fehler beim Parsen der E-Mail-Daten: %s", err.Error())
		respondWithError(c, http.StatusBadRequest, "Valid email address required")
		return
	}

	var user User
	if err := user.FirstId(userID.(uint)); err != nil {
		logger.Error("Benutzer mit ID %v nicht gefunden: %s", userID, err.Error())
		respondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	// Check if email is already in use by another user
	var existingUser User
	if err := existingUser.FindByEmail(input.Email); err == nil {
		if existingUser.UserID != user.UserID {
			logger.Warn("E-Mail-Aktualisierung fehlgeschlagen - E-Mail bereits vergeben: %s", input.Email)
			respondWithError(c, http.StatusConflict, "Email already in use")
			return
		}
	}

	user.Email = input.Email
	if err := user.Save(); err != nil {
		logger.Error("Fehler beim Speichern der E-Mail für Benutzer %s: %s", user.Username, err.Error())
		respondWithError(c, http.StatusInternalServerError, "Failed to update email")
		return
	}

	logger.Info("E-Mail erfolgreich aktualisiert für Benutzer: %s (ID: %d)", user.Username, user.UserID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Email updated successfully",
		"email":   user.Email,
	})
}

// UpdatePassword Handler to update user's password
// UpdatePassword godoc
// @Summary Update user password
// @Description Updates the password of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Param passwords body object{old_password=string,new_password=string} true "Old and new passwords"
// @Success 200 {object} map[string]string "Password updated successfully"
// @Failure 400 {object} map[string]string "Invalid password data"
// @Failure 401 {object} map[string]string "Unauthorized or wrong old password"
// @Security BearerAuth
// @Router /api/user/password [put]
func UpdatePassword(c *gin.Context) {
	logger.Debug("Starte Passwort-Aktualisierung...")

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		logger.Error("Benutzer-ID nicht im Context gefunden")
		respondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Fehler beim Parsen der Passwort-Daten: %s", err.Error())
		respondWithError(c, http.StatusBadRequest, "Current password and new password (min 6 characters) required")
		return
	}

	var user User
	if err := user.FirstId(userID.(uint)); err != nil {
		logger.Error("Benutzer mit ID %v nicht gefunden: %s", userID, err.Error())
		respondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	// Verify current password
	hashedCurrentPassword := md5.Sum([]byte(input.CurrentPassword))
	if user.PasswordHash != hex.EncodeToString(hashedCurrentPassword[:]) {
		logger.Warn("Passwort-Aktualisierung fehlgeschlagen - Aktuelles Passwort ungültig für Benutzer: %s", user.Username)
		respondWithError(c, http.StatusUnauthorized, "Current password is incorrect")
		return
	}

	// Hash new password
	hashedNewPassword := md5.Sum([]byte(input.NewPassword))
	user.PasswordHash = hex.EncodeToString(hashedNewPassword[:])

	if err := user.Save(); err != nil {
		logger.Error("Fehler beim Speichern des Passworts für Benutzer %s: %s", user.Username, err.Error())
		respondWithError(c, http.StatusInternalServerError, "Failed to update password")
		return
	}

	logger.Info("Passwort erfolgreich aktualisiert für Benutzer: %s (ID: %d)", user.Username, user.UserID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

// UpdateLanguage Handler to update user's preferred language
// UpdateLanguage godoc
// @Summary Update preferred language
// @Description Updates the preferred language setting for the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Param language body object{language=string} true "Language code (e.g., 'en', 'de')"
// @Success 200 {object} map[string]string "Language updated successfully"
// @Failure 400 {object} map[string]string "Invalid language code"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/user/language [put]
func UpdateLanguage(c *gin.Context) {
	logger.Debug("Starte Sprach-Aktualisierung...")

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		logger.Error("Benutzer-ID nicht im Context gefunden")
		respondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input struct {
		Language string `json:"language" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Fehler beim Parsen der Sprach-Daten: %s", err.Error())
		respondWithError(c, http.StatusBadRequest, "Language is required")
		return
	}

	// Validate language (only de and en supported)
	if input.Language != "de" && input.Language != "en" {
		logger.Warn("Ungültige Sprache angefordert: %s", input.Language)
		respondWithError(c, http.StatusBadRequest, "Invalid language. Supported languages: de, en")
		return
	}

	var user User
	if err := user.FirstId(userID.(uint)); err != nil {
		logger.Error("Benutzer mit ID %v nicht gefunden: %s", userID, err.Error())
		respondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	user.PreferredLanguage = input.Language
	if err := user.Save(); err != nil {
		logger.Error("Fehler beim Speichern der Sprache für Benutzer %s: %s", user.Username, err.Error())
		respondWithError(c, http.StatusInternalServerError, "Failed to update language")
		return
	}

	logger.Info("Sprache erfolgreich aktualisiert für Benutzer: %s (ID: %d) - Neue Sprache: %s", user.Username, user.UserID, input.Language)
	c.JSON(http.StatusOK, gin.H{
		"message":  "Language updated successfully",
		"language": user.PreferredLanguage,
	})
}

// UpdateDisplayName Handler to update user's display name
// UpdateDisplayName godoc
// @Summary Update display name
// @Description Updates the display name of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Param display_name body object{display_name=string} true "New display name"
// @Success 200 {object} map[string]string "Display name updated successfully"
// @Failure 400 {object} map[string]string "Invalid display name"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/user/display-name [put]
func UpdateDisplayName(c *gin.Context) {
	logger.Debug("Starte Anzeigenamen-Aktualisierung...")

	userID, exists := c.Get("userID")
	if !exists {
		logger.Error("Benutzer-ID nicht im Context gefunden")
		respondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input struct {
		DisplayName string `json:"display_name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Error("Fehler beim Parsen des Anzeigenamens: %s", err.Error())
		respondWithError(c, http.StatusBadRequest, "Display name is required")
		return
	}

	if utf8.RuneCountInString(input.DisplayName) > 30 {
		logger.Warn("Anzeigename zu lang: %d Zeichen", utf8.RuneCountInString(input.DisplayName))
		respondWithError(c, http.StatusBadRequest, "Display name must be at most 30 characters")
		return
	}

	var user User
	if err := user.FirstId(userID.(uint)); err != nil {
		logger.Error("Benutzer mit ID %v nicht gefunden: %s", userID, err.Error())
		respondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	user.DisplayName = input.DisplayName
	if err := user.Save(); err != nil {
		logger.Error("Fehler beim Speichern des Anzeigenamens für Benutzer %s: %s", user.Username, err.Error())
		respondWithError(c, http.StatusInternalServerError, "Failed to update display name")
		return
	}

	logger.Info("Anzeigename erfolgreich aktualisiert für Benutzer: %s (ID: %d)", user.Username, user.UserID)
	c.JSON(http.StatusOK, gin.H{
		"message":      "Display name updated successfully",
		"display_name": user.DisplayNameOrUsername(),
	})
}
