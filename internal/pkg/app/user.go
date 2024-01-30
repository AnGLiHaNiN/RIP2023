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
// @Tags		Пользователь
// @Description	Регистрация нового пользователя
// @Accept		json
// @Param		user_credentials body schemes.RegisterReq true "login and password"
// @Success		200
// @Router		/api/user/sign_up [post]
func (app *Application) Register(c *gin.Context) {
	request := &schemes.RegisterReq{}
	if err := c.ShouldBind(request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	existing_user, err := app.repo.GetUserByLogin(request.Login)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if existing_user != nil {
		c.AbortWithStatus(http.StatusBadRequest)
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

	c.Status(http.StatusOK)
}

// @Summary		Авторизация
// @Tags		Пользователь
// @Description	Авторизует пользователя по логиню, паролю и отдаёт jwt токен для дальнейших запросов
// @Accept		json
// @Produce		json
// @Param		user_credentials body schemes.LoginReq true "login and password"
// @Success		200 {object} schemes.AuthResp
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
		Role:     user.Role,
		Login:    user.Login,
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

	c.JSON(http.StatusOK, schemes.AuthResp{
		AccessToken: strToken,
		TokenType:   "Bearer",
	})
}

// @Summary		Выйти из аккаунта
// @Tags		Пользователь
// @Description	Выход из аккаунта
// @Success		200
// @Router		/api/user/loguot [get]
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

// @Summary		Профиль
// @Tags		Пользователь
// @Description	Получить информацию о профиле пользователя
// @Success		200
// @Router		/api/user [get]
func (app *Application) Profile(c *gin.Context) {
	userId := getUserId(c)
	user, err := app.repo.GetUserById(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, *user)
}

// @Summary		Изменить профиль
// @Tags		Пользователь
// @Description	Изменить все или часть данных профиля пользователя
// @Access		json
// @Param		new_fields body schemes.ChangeUserReq true "Новые значения"
// @Success		200
// @Router		/api/user [put]
func (app *Application) UpdateUser(c *gin.Context) {
	var request schemes.ChangeUserReq
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userId := getUserId(c)
	user, err := app.repo.GetUserById(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if user == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if request.Email != nil {
		user.Email = request.Email
	}
	if request.Name != nil {
		user.Name = request.Name
	}
	if request.Password != nil {
		user.Password = generateHashString(*request.Password)
	}
	if err := app.repo.SaveUser(user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
