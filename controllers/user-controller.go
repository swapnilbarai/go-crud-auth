package controllers

import (
	"net/http"
	"strings"
	"swapnilbarai/go-crud-auth/models"
	"swapnilbarai/go-crud-auth/utils"

	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware(c *gin.Context) {
	splitTokens := strings.Split(c.GetHeader("Authorization"), " ")
	authorizationToken := ""
	if len(splitTokens) == 2 {
		authorizationToken = splitTokens[1]
	}
	if authorizationToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access-token"})
		c.Abort()
		return
	}
	jwtToken, err := utils.VerifyJWTToken(authorizationToken, "access")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	if !jwtToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
		c.Abort()
		return
	}
	c.Next()
}

func GetUserDetails(c *gin.Context) {
	userName := c.Param("username")

	if len(userName) == 0 {
		invalidUserMessage := "Please provide valid username or password"
		c.JSON(http.StatusNotAcceptable, gin.H{"error": invalidUserMessage})
		return
	}
	if _, exists := models.Users[userName]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found.Please enter valid user"})
		return
	}
	user := models.Users[userName]
	c.JSON(http.StatusOK, gin.H{"email": user.Email, "mobileNumber": user.MobileNo})

}

func RevokeUserToken(c *gin.Context) {
	var primaryTokenKey models.TokenPrimaryKey
	if err := c.BindJSON(&primaryTokenKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, exists := models.Tokens[primaryTokenKey]; !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token doesn't exsists"})
		return
	}
	models.InvalidateToken(primaryTokenKey.IssueAt, primaryTokenKey.Subject, primaryTokenKey.TokenType)
	c.JSON(http.StatusOK, gin.H{"message": "Token succesfully revoked"})
}
