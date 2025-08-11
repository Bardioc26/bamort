/*
User Handlers

Add handlers for user registration and login:
*/
package user

import (
	"bamort/logger"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func respondWithError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func RegisterUser(c *gin.Context) {
	logger.Debug("Starte Benutzerregistrierung...")

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Error("Fehler beim Parsen der Registrierungsdaten: %s", err.Error())
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debug("Registriere Benutzer: %s", user.Username)
	//fmt.Printf("User input: '%s'", user.PasswordHash)
	//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	hashedPassword := md5.Sum([]byte(user.PasswordHash))
	user.PasswordHash = hex.EncodeToString(hashedPassword[:])
	logger.Debug("Passwort-Hash erstellt für Benutzer: %s", user.Username)

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

		c.Next()
	}
}
