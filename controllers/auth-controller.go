package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"swapnilbarai/go-crud-auth/models"
	"swapnilbarai/go-crud-auth/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func SignUpUser(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// checking for valid email format
	if ok := utils.IsEmailValid(user.Email); !ok {
		invalidEmailMessage := utils.FormatInvalidMeesage("email", user.Email)
		c.JSON(http.StatusNotAcceptable, gin.H{"error": invalidEmailMessage})
		return
	}

	// checking for valid mobile number format
	if ok := utils.IsMobileNumberValid(user.MobileNo); !ok {
		invalidMobileNumberMessage := utils.FormatInvalidMeesage("mobile number", user.MobileNo)
		c.JSON(http.StatusNotAcceptable, gin.H{"error": invalidMobileNumberMessage})
		return
	}

	// checking for valid username
	//only accept user with length greater than 6
	if len(user.UserName) < 6 {
		invalidUserNameMessage := utils.FormatInvalidMeesage("user name", user.UserName)
		c.JSON(http.StatusNotAcceptable, gin.H{"error": invalidUserNameMessage})
		return
	}

	//checking for user role
	if user.Role != utils.Admin && user.Role != utils.Normal {
		invalidUserRoleMessage := utils.FormatInvalidMeesage("role", user.Role)
		c.JSON(http.StatusNotAcceptable, gin.H{"error": invalidUserRoleMessage})
		return
	}

	// checking for valid password
	//only accept user with length greater than 6
	if len(user.PassWord) < 6 {
		invalidPassportMessage := utils.FormatInvalidMeesage("password", user.PassWord)
		c.JSON(http.StatusNotAcceptable, gin.H{"error": invalidPassportMessage})
		return
	}

	//checking for already signup user
	if _, exists := models.Users[user.UserName]; exists {
		userAlreadyExists := fmt.Sprintf("User %s is already exists\n", user.UserName)
		c.JSON(http.StatusConflict, gin.H{"error": userAlreadyExists})
		return
	}

	hashPassord := utils.HashPassword(user.PassWord)
	user.PassWord = hashPassord
	models.Users[user.UserName] = user
	okayMeesage := fmt.Sprintf("User: %s is succesfully created.Please try login\n", user.UserName)
	c.JSON(http.StatusCreated, gin.H{"message": okayMeesage})

}

func SignInUser(c *gin.Context) {
	userName := c.PostForm("username")
	password := c.PostForm("password")

	if len(userName) == 0 || len(password) == 0 {
		invalidUserMessage := "Please provide valid username or password"
		c.JSON(http.StatusNotAcceptable, gin.H{"error": invalidUserMessage})
		return
	}

	if _, exists := models.Users[userName]; !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please signup before login"})
		return
	}

	//comparing password
	hashPassword := utils.HashPassword(password)
	storedPassword := models.Users[userName].PassWord
	if hashPassword != storedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please enter valid password"})
		return
	}
	//giving back token
	createNewToken(c, userName)
	c.JSON(http.StatusOK, gin.H{"message": "sucessfully logged in"})

}

func RefreshToken(c *gin.Context) {
	//token format "bearer x.y.z"
	splitTokens := strings.Split(c.GetHeader("Refresh-Authorization"), " ")
	authorizationToken := ""
	if len(splitTokens) == 2 {
		authorizationToken = splitTokens[1]
	}
	if authorizationToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing refresh-token"})
		return
	}
	//verifying jwt token
	jwtToken, err := utils.VerifyJWTToken(authorizationToken, models.RefrehTokenType)

	//other error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//expire case or any other error
	if !jwtToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	//giving new refresh, access token
	userName, userNameError := jwtToken.Claims.GetSubject()
	issueAt, issuerTimeError := jwtToken.Claims.GetIssuedAt()

	if userNameError != nil || issuerTimeError != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Malformed Token"})
		return
	}
	//invalidating previous token
	models.InvalidateToken(issueAt.Unix(), userName, models.RefrehTokenType)
	//creating new tokens
	createNewToken(c, userName)
	c.JSON(http.StatusOK, gin.H{"message": "sucessfully refresh tokens"})

}

func createNewToken(c *gin.Context, userName string) {
	issuerTime := time.Now().Unix() //for making sure issue time of access and refresh token should be same

	//access token
	accessToken, err := utils.CreateJWTToken(userName, false, "access", int(utils.AccessTokenDuration), issuerTime)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//refresh token
	refreshToken, err := utils.CreateJWTToken(userName, false, "refresh", int(utils.RefreshTokenDuration), issuerTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//adding access token and refresh token  header to logged in or refresh token request
	c.Header("Authorization", "Bearer "+accessToken)
	c.Header("Refresh-Authorization", "Bearer "+refreshToken)

}
