package models

type Token struct {
	TokenPrimaryKey
	IsActive bool
}

type TokenPrimaryKey struct {
	IssueAt   int64  `json:"issueAt"`
	Subject   string `json:"subject"`
	TokenType string `json:"tokenType"`
}

var Tokens map[TokenPrimaryKey]*Token

func IsTokenValid(issueAt int64, subject string, tokenType string) bool {
	token, isExists := Tokens[TokenPrimaryKey{IssueAt: issueAt, Subject: subject, TokenType: tokenType}]
	if !isExists {
		panic("If token is allocated then token should exists in database")
	}
	if token.IsActive {
		return true
	}
	return false
}

func InvalidateToken(issueAt int64, subject string, tokenType string) {
	token, isExists := Tokens[TokenPrimaryKey{IssueAt: issueAt, Subject: subject, TokenType: tokenType}]
	if !isExists {
		return
	}
	if token.IsActive {
		token.IsActive = false
	}
	accessToken, accessTokenIsExists := Tokens[TokenPrimaryKey{IssueAt: issueAt, Subject: subject, TokenType: "access"}]
	if tokenType == "refresh" && accessTokenIsExists && accessToken.IsActive {
		accessToken.IsActive = false
	}

}
func InsertToken(issueAt int64, subject string, tokenType string) {
	tokenPrimaryKey := TokenPrimaryKey{IssueAt: issueAt, Subject: subject, TokenType: tokenType}
	token := &Token{TokenPrimaryKey: tokenPrimaryKey, IsActive: true}
	Tokens[tokenPrimaryKey] = token
	return
}
