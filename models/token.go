package models

import (
	"math/rand"
	"time"
)

const RefrehTokenType = "refresh"
const AccessTokenType = "access"

type Token struct {
	TokenID   int64  `json:"tokenID"`
	IssueAt   int64  `json:"issueAt"`
	Subject   string `json:"subject"`
	TokenType string `json:"tokenType"`
	IsActive  bool   `json:"isActive"`
}

func (t *Token) isTokenEqual(issueAt int64, subject string, tokenType string) bool {
	if t.IssueAt == issueAt && t.Subject == subject && t.TokenType == tokenType {
		return true
	}
	return false
}

type RandomNumberGenerator struct {
}

func (r RandomNumberGenerator) Seed() {
	rand.Seed(time.Now().UnixNano())
}
func (r RandomNumberGenerator) GenerateId() int64 {
	return rand.Int63()
}

var IdGenerator RandomNumberGenerator

var Tokens map[int64]*Token

func findToken(issueAt int64, subject string, tokenType string) *Token {
	for _, token := range Tokens {
		if token.isTokenEqual(issueAt, subject, tokenType) {
			return token
		}
	}
	return nil
}

func IsTokenValid(issueAt int64, subject string, tokenType string) bool {
	token := findToken(issueAt, subject, tokenType)
	if token == nil || !token.IsActive {
		return false
	}

	return true
}

func InvalidateToken(issueAt int64, subject string, tokenType string) bool {
	token := findToken(issueAt, subject, tokenType)
	if token == nil || !token.IsActive {
		return false
	}
	if token.IsActive {
		token.IsActive = false
	}
	accessToken := findToken(issueAt, subject, AccessTokenType)
	if tokenType == RefrehTokenType && accessToken != nil && accessToken.IsActive {
		accessToken.IsActive = false
	}
	return true

}
func InvalidateTokenByID(tokenId int64) bool {
	token, isExist := Tokens[tokenId]
	if !isExist || !token.IsActive {
		return false
	}
	token.IsActive = false
	accessToken := findToken(token.IssueAt, token.Subject, AccessTokenType)
	if token.TokenType == RefrehTokenType && accessToken != nil && accessToken.IsActive {
		accessToken.IsActive = false
	}
	return true
}
func InsertToken(issueAt int64, subject string, tokenType string) {
	tokenID := IdGenerator.GenerateId()
	token := &Token{TokenID: tokenID, IssueAt: issueAt, Subject: subject, TokenType: tokenType, IsActive: true}

	Tokens[tokenID] = token

}
func GetActiveToken() []*Token {
	var filteredTokens []*Token
	for _, token := range Tokens {
		if token.IsActive {
			filteredTokens = append(filteredTokens, token)
		}
	}
	return filteredTokens
}
