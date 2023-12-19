package app

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	
	"R_I_P_labs/internal/app/ds"
	"R_I_P_labs/internal/app/role"
	"R_I_P_labs/internal/app/schemes"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)


// @Summary		Регистрация
// @Tags		Авторизация
// @Description	Регистрация нового пользователя
// @Accept		json
// @Produce		json
// @Param		user_credentials body schemes.RegisterReq true "login and password"
// @Success		200 {object} schemes.RegisterResp
// @Router		/api/user/sign_up [post]
func (app *Application) Register(c *gin.Context) {
	request := &schemes.RegisterReq{}
	if err := c.ShouldBind(request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if request.Password == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("password is empty"))
		return
	}

	if request.Login == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("login is empty"))
		return
	}

	if err := app.repo.AddUser(&ds.User{
		Role:     role.Customer,
		Login:    request.Login,
		Password: generateHashString(request.Password),
	}); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &schemes.RegisterResp{
		Ok: true,
	})
}

// @Summary		Авторизация
// @Tags		Авторизация
// @Description	Авторизует пользователя по логиню, паролю и отдаёт jwt токен для дальнейших запросов
// @Accept		json
// @Produce		json
// @Param		user_credentials body schemes.LoginReq true "login and password"
// @Success		200 {object} schemes.SwaggerLoginResp
// @Router		/api/user/login [post]
// @Consumes    json
func (app *Application) Login(c *gin.Context) {
	JWTConfig := app.config.JWT
	request := &schemes.LoginReq{}
	if err := c.ShouldBind(request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := app.repo.GetUserByLogin(request.Login)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}


	if user.Password != generateHashString(request.Password) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	token := jwt.NewWithClaims(JWTConfig.SigningMethod, &ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(JWTConfig.ExpiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "bitop-admin",
		},
		UserUUID: user.UUID,
		Role:	  user.Role,
	})
	if token == nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token is nil"))
		return
	}

	strToken, err := token.SignedString([]byte(JWTConfig.Token))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cant create str token"))
		return
	}

	c.JSON(http.StatusOK, schemes.LoginResp{
		ExpiresIn:   JWTConfig.ExpiresIn,
		AccessToken: strToken,
		TokenType:   "Bearer",
	})
}

// @Summary		Выйти из аккаунта
// @Tags		Авторизация
// @Description	Выход из аккаунта
// @Accept		json
// @Produce		json
// @Success		200
// @Router		/api/user/loguot [post]
func (app *Application) Logout(c *gin.Context) {
	jwtStr := c.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, jwtPrefix) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	jwtStr = jwtStr[len(jwtPrefix):]

	_, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.config.JWT.Token), nil
	})
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		log.Println(err)
		return
	}

	err = app.redisClient.WriteJWTToBlacklist(c.Request.Context(), jwtStr, app.config.JWT.ExpiresIn)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}