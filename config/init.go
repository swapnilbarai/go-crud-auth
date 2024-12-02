package config

import (
	"os"
	"strconv"
	"swapnilbarai/go-crud-auth/models"
	"swapnilbarai/go-crud-auth/utils"
)

func Initialise() {
	models.Users = make(map[string]models.User)
	models.Tokens = make(map[int64]*models.Token)
	models.IdGenerator = models.RandomNumberGenerator{}
	models.IdGenerator.Seed()
	if os.Getenv("ACCESS_SECRET") != "" {
		utils.AccessSecret = []byte(os.Getenv("ACCESS_SECRET"))
	}
	if os.Getenv("REFRFESH_SECRET") != "" {
		utils.RefreshSecret = []byte(os.Getenv("REFRESH_SECRET"))
	}
	if os.Getenv("ACCESS_TOKEN_DURATION") != "" {
		duration, err := strconv.ParseInt(os.Getenv("ACCESS_TOKEN_DURATION"), 10, 64)
		if err != nil {
			panic(err.Error())
		}
		utils.AccessTokenDuration = duration
	}
	if os.Getenv("REFRESH_TOKEN_DURATION") != "" {
		duration, err := strconv.ParseInt(os.Getenv("REFRESH_TOKEN_DURATION"), 10, 64)
		if err != nil {
			panic(err.Error())
		}
		utils.RefreshTokenDuration = duration
	}

}
