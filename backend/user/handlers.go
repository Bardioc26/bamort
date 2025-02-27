/*
User Handlers

Add handlers for user registration and login:
*/
package user

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//fmt.Printf("User input: '%s'", user.PasswordHash)
	//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	hashedPassword := md5.Sum([]byte(user.PasswordHash))
	user.PasswordHash = hex.EncodeToString(hashedPassword[:])
	//fmt.Printf("pwdh: %s", user.PasswordHash)
	if err := user.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create user: %s", err)})
		return
	}
	//fmt.Printf(" ___ pwdh2: %s", user.PasswordHash)
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully:"})
}

func GenerateToken(u *User) string {
	//u.Username + "lkiuztrew" + u.CreatedAt.String()
	tx := md5.Sum([]byte(u.Username + u.CreatedAt.String()))
	// Convert hash to raw string
	hashString := hex.EncodeToString(tx[:])
	pos := 7
	idm := "." + fmt.Sprintf("%d", u.UserID) + ":"
	// Insert the character
	token := hashString[:pos] + string(idm) + hashString[pos:]
	return token
}
func CheckToken(token string) *User {
	//fmt.Print("CheckToken1: " + token)
	var u User
	var err error
	pos := 7 + len("Bearer ")
	userid := 0
	// Check if a `.` is at position 7 (zero-indexed)
	if len(token) > pos && token[pos] == '.' {
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
				//fmt.Print("CheckToken4: " + err.Error() + "\n")
				return nil
			}
		} else {
			//fmt.Print("CheckToken5: not found\n")
			return nil
		}
	} else {
		//fmt.Print("CheckToken6: not found\n")
		return nil
	}

	if userid > 0 {
		//fmt.Printf("CheckToken6-1: userid %v\n", userid)
		//fmt.Printf("CheckToken6-1: userid %v\n", uint(userid))
		err := u.FirstId(uint(userid))
		if err != nil {
			//fmt.Printf("CheckToken7: not found error %s\n", err.Error())
			return nil
		}
		//fmt.Printf("CheckToken8: found:%s \n", u.Username)
		return &u
	}
	//fmt.Print("CheckToken9: not found\n")
	return nil
}

func LoginUser(c *gin.Context) {
	var user User
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
	if err := user.First(input.Username); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid username. or password %v", input)})
		return
	}

	hashedPassword := md5.Sum([]byte(input.Password))
	fmt.Printf("pwdh: %s", hex.EncodeToString(hashedPassword[:]))
	if user.PasswordHash != hex.EncodeToString(hashedPassword[:]) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid username. or password. %s %s", input.Password, hex.EncodeToString(hashedPassword[:]))})
		return
	}
	/*
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password."})
			return
		}
	*/

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": GenerateToken(&user)})
}

// Apply middleware to protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		if CheckToken(token) == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized."})
			c.Abort()
			return
		}
		// Add token validation logic here

		c.Next()
	}
}
