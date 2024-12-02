package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"swapnilbarai/go-crud-auth/models"
	"swapnilbarai/go-crud-auth/utils"

	"github.com/gin-gonic/gin"
)

// middleware :protecting user routes
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

	//user able to see active tokens or revoke tokens
	//if he is admin
	//only for route /user/active/tokens or user/revoke/{tokenId}
	if utils.PathIsProtected(c.FullPath()) {
		userName, _ := jwtToken.Claims.GetSubject()
		user, exist := models.Users[userName]
		if !exist || user.Role != utils.Admin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not authurized for to see active or revoke tokens"})
			c.Abort()
			return
		}

	}
	c.Next()
}

// return details asked user
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

// Revoke the token
// if tokenType==refresh then it is also revoke access token if it is active and issue with given refresh token
func RevokeUserToken(c *gin.Context) {
	tokenID, err := strconv.ParseInt(c.Param("tokenID"), 10, 64)

	if err != nil {
		invalidTokenMessage := "Please provide  tokenID"
		c.JSON(http.StatusNotAcceptable, gin.H{"error": invalidTokenMessage})
		return
	}
	if _, exists := models.Tokens[tokenID]; !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token doesn't exsists"})
		return
	}
	models.InvalidateTokenByID(tokenID)
	c.JSON(http.StatusOK, gin.H{"message": "Token succesfully revoked"})
}

// return Json array of Active Tokens
func ShowActiveTokens(c *gin.Context) {
	tokens := models.GetActiveToken()
	c.JSON(http.StatusOK, tokens)

}
