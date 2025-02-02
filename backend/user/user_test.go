package user

import (
	"bamort/database"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	database.SetupTestDB()
	usr := User{
		Username:     "bebe",
		PasswordHash: "osiris",
		Email:        "frank@wuenscheonline.de",
	}

	hashedPassword := md5.Sum([]byte(usr.PasswordHash))
	usr.PasswordHash = hex.EncodeToString(hashedPassword[:])
	err := usr.Create()
	assert.NoError(t, err, "no error expected when creating record")

	usr2 := User{
		Username:     "bubnu",
		PasswordHash: "osiris",
		Email:        "spacer@wuenscheonline.de",
	}
	hashedPassword = md5.Sum([]byte(usr2.PasswordHash))
	usr2.PasswordHash = hex.EncodeToString(hashedPassword[:])
	err = usr2.Create()
	assert.NoError(t, err, "no error expected when creating record")
}

func TestLoginUser(t *testing.T) {
	TestRegisterUser(t)
	var usr User
	input := struct {
		Username       string `json:"username"`
		Password       string `json:"password"`
		HashedPassword string
	}{
		Username: "bebe",
		Password: "osiris",
	}
	err := usr.First(input.Username)
	assert.NoError(t, err, "no error expected when finding record")

	hashedPassword := md5.Sum([]byte(input.Password))
	input.HashedPassword = hex.EncodeToString(hashedPassword[:])
	assert.Equal(t, input.HashedPassword, usr.PasswordHash)

}

func TestHshing(t *testing.T) {
	TestRegisterUser(t)
	var u1 User
	u1.First("bebe")
	assert.Equal(t, "", u1.Username+u1.CreatedAt.String())
	tx := md5.Sum([]byte(u1.Username + u1.CreatedAt.String()))
	assert.NotEmpty(t, tx)
	// Convert hash to raw string
	hashString := hex.EncodeToString(tx[:])
	assert.Equal(t, "", hashString)
	pos := 7
	idm := "." + strconv.Itoa(int(u1.UserID)) + ":"
	// Insert the character
	token := hashString[:pos] + idm + hashString[pos:]
	assert.Equal(t, "", token)

	// check
	var u User
	var err error
	userid := 0

	// Check if a `.` is at position 7 (zero-indexed)
	if len(token) > pos && token[pos] == '.' {
		assert.Equal(t, ". ", token[pos])
		// Find the next `:` after the `.`
		colonPos := strings.Index(token[pos+1:], ":") // Start searching after position 7
		if colonPos != -1 {
			// Extract the substring between `.` and `:`
			uu := token[pos+1 : pos+1+colonPos]
			assert.Equal(t, "1 ", uu)
			//fmt.Println("Extracted Substring:", result)
			userid, err = strconv.Atoi(uu)

			assert.NoError(t, err, "no error expexted when strconv")
			assert.Equal(t, 2, userid)
		}
	}
	if userid > 0 {
		err = u.FirstId(uint(userid))
		assert.NoError(t, err, "no error expexted when fetching user")

	}
}

func TestCors(t *testing.T) {
	database.SetupTestDB()
	us := User{
		Username:     "bebe",
		UserID:       1,
		PasswordHash: "5f29e63a3f26798930e5bf218445164f",
		//CreatedAt: "2025-01-04 09:01:44.911",
	}
	token := GenerateToken(&us)
	fmt.Print(token)
	usr := CheckToken("Bearer " + token)
	fmt.Print(usr)
}
