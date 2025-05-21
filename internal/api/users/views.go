package users

import (
	"net/http"
	api "shortener/internal/api"
	"shortener/internal/api/urls"
	conf "shortener/internal/config"
	"shortener/internal/database"
	"time"

	_ "shortener/docs"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserDTO struct {
	Username string `json:"username" validate:"required,min=5,max=16"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type UserLoginDTO struct {
	Username string `json:"username" validate:"required,min=5,max=16"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type TokenDTO struct {
	Token string `json:"token" validate:"required,jwt"`
}

type JWTData struct {
	ID       uint   `json:"id" validate:"required,gt=0"`
	Username string `json:"username" validate:"required,min=5,max=16"`
	Email    string `json:"email" validate:"required,email"`
}

type RegenerateCodeDTO struct {
	ID    int    `json:"id" validate:"required,gt=0"`
	Email string `json:"email" validate:"required,email"`
}

type VerificationDTO struct {
	ID   int    `json:"id" validate:"required,gt=0"`
	Code string `json:"code" validate:"len=5,required,number"`
}

type ResetPasswordDTO struct {
	Username string `json:"username" validate:"required,min=5,max=16"`
}

type ResetPasswordConfirmDTO struct {
	ID       int    `json:"id" validate:"required,gt=0"`
	Token    string `json:"token" validate:"len=64,required"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type UpdateAccountDTO struct {
	Email string `json:"email" validate:"required,email"`
}

// @Summary Регистрация пользователя
// @Description Создает нового пользователя
// @Accept json
// @Produce json
// @Param user body users.UserDTO true "Пользователь"
// @Success 200 {object} api.Response "Успешная регистрация"
// @Failure 400 {object} api.Response "Некорректные данные"
// @Failure 500 {object} api.Response "Внутренняя ошибка сервера"
// @Router /api/register [post]
func Register(c echo.Context) error {
	dto := UserDTO{}
	if err := api.DecodeRequest(c, &dto); err != nil {
		return err
	}

	hashedPassword, err := GetHashedPassword(dto.Password)
	if err != nil {
		return err
	}

	r := Repository{db: database.GetDB()}
	user, err := r.SaveUser(dto.Username, dto.Email, string(hashedPassword))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.Response{Msg: "Username or email is already exists"})
	}

	if err := CreateAndSendToken(user.ID, user.Email); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, api.Response{Msg: "success", Data: user})
}

// @Summary Авторизация пользователя
// @Description Авторизует пользователя
// @Accept json
// @Produce json
// @Param user body users.UserLoginDTO true "Пользователь"
// @Success 200 {object} api.Response "Успешная авторизация"
// @Failure 401 {object} api.Response "Некорректные данные"
// @Failure 500 {object} api.Response "Внутренняя ошибка сервера"
// @Router /api/login [post]
func Login(c echo.Context) error {
	dto := UserLoginDTO{}
	if err := api.DecodeRequest(c, &dto); err != nil {
		return err
	}

	r := Repository{db: database.GetDB()}
	user, err := r.GetUser(UserInfo{Username: dto.Username})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, api.Response{Msg: "User is not found or not verified"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	config := conf.GetConfig()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": JWTData{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
		"exp": time.Now().Add(time.Minute * time.Duration(config.JWT.JWT_TTL)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.JWT.SECRET))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not generate token")
	}

	return c.JSON(http.StatusOK, api.Response{Msg: "success", Data: TokenDTO{Token: tokenString}})
}

// @Summary Регенерация кода подтверждения
// @Description Регенерирует код подтверждения для пользователя
// @Accept json
// @Produce json
// @Param user body users.RegenerateCodeDTO true "Пользователь"
// @Success 200 {object} api.Response "Успешная регенерация"
// @Failure 400 {object} api.Response "Некорректные данные"
// @Failure 500 {object} api.Response "Внутренняя ошибка сервера"
// @Router /api/regenerate-code [post]
func RegenerateCode(c echo.Context) error {
	dto := RegenerateCodeDTO{}
	if err := api.DecodeRequest(c, &dto); err != nil {
		return err
	}
	if err := CreateAndSendToken(uint(dto.ID), dto.Email); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, api.Response{Msg: "success"})
}

// @Summary Подтверждение почты
// @Description Подтверждает почту пользователя
// @Accept json
// @Produce json
// @Param user body users.VerificationDTO true "Пользователь"
// @Success 200 {object} api.Response "Успешное подтверждение"
// @Failure 400 {object} api.Response "Некорректные данные"
// @Failure 500 {object} api.Response "Внутренняя ошибка сервера"
// @Router /api/verification [post]
func ActivateAccount(c echo.Context) error {
	dto := VerificationDTO{}
	if err := api.DecodeRequest(c, &dto); err != nil {
		return err
	}

	config := conf.GetConfig()
	if err := CheckRedisToken(dto.ID, dto.Code, config.OTP.RedisName); err != nil {
		return err
	}

	r := Repository{db: database.GetDB(), UserID: dto.ID}
	if err := r.VerificateUser(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, api.Response{Msg: "Internal server error. Please try again"})
	}

	return c.JSON(http.StatusAccepted, api.Response{Msg: "success"})
}

// @Summary Сброс пароля
// @Description Сбрасывает пароль пользователя
// @Accept json
// @Produce json
// @Param user body users.ResetPasswordDTO true "Пользователь"
// @Success 200 {object} api.Response "Успешный сброс пароля"
// @Failure 400 {object} api.Response "Некорректные данные"
// @Failure 500 {object} api.Response "Внутренняя ошибка сервера"
// @Router /api/reset-password [post]
func ResetPassword(c echo.Context) error {
	dto := ResetPasswordDTO{}
	if err := api.DecodeRequest(c, &dto); err != nil {
		return err
	}

	r := Repository{db: database.GetDB()}
	user, err := r.GetUser(UserInfo{Username: dto.Username})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.Response{Msg: "User not found"})
	}

	if err := CreateAndSendResetPasswordLink(user.ID, user.Email); err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, api.Response{Msg: "success"})
}

// @Summary Подтверждение сброса пароля
// @Description Подтверждает сброс пароля пользователя
// @Accept json
// @Produce json
// @Param user body users.ResetPasswordConfirmDTO true "Пользователь"
// @Success 200 {object} api.Response "Успешное подтверждение сброса пароля"
// @Failure 400 {object} api.Response "Некорректные данные"
// @Failure 500 {object} api.Response "Внутренняя ошибка сервера"
// @Router /api/reset-password-confirm [post]
func ResetPasswordConfirm(c echo.Context) error {
	dto := ResetPasswordConfirmDTO{}
	if err := api.DecodeRequest(c, &dto); err != nil {
		return err
	}

	config := conf.GetConfig()
	if err := CheckRedisToken(dto.ID, dto.Token, config.ResetToken.RedisName); err != nil {
		return err
	}

	hashedPassword, err := GetHashedPassword(dto.Password)
	if err != nil {
		return err
	}

	r := Repository{db: database.GetDB(), UserID: dto.ID}
	if err := r.ChangeUserPassword(string(hashedPassword)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, api.Response{Msg: "Internal server error. Please try again"})
	}

	return c.JSON(http.StatusOK, api.Response{Msg: "success"})
}

// @Summary Получение информации о текущем пользователе
// @Description Возвращает информацию о текущем пользователе
// @Accept json
// @Produce json
// @Success 200 {object} api.Response "Успешный запрос"
// @Failure 401 {object} api.Response "Неавторизованный запрос"
// @Failure 500 {object} api.Response "Внутренняя ошибка сервера"
// @Router /api/account [get]
func GetMe(c echo.Context) error {
	userID := urls.GetUserID(c)
	r := Repository{db: database.GetDB()}
	user, err := r.GetUser(UserInfo{ID: userID})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, api.Response{Msg: "User not found"})
	}
	return c.JSON(http.StatusOK, api.Response{Msg: "success", Data: user})
}

// @Summary Обновление аккаунта
// @Description Обновляет информацию о пользователе
// @Accept json
// @Produce json
// @Param user body users.UpdateAccountDTO true "Пользователь"
// @Success 200 {object} api.Response "Успешное обновление"
// @Failure 400 {object} api.Response "Некорректные данные"
// @Failure 500 {object} api.Response "Внутренняя ошибка сервера"
// @Router /api/account [patch]
func UpdateAccount(c echo.Context) error {
	userID := urls.GetUserID(c)
	dto := UpdateAccountDTO{}
	if err := api.DecodeRequest(c, &dto); err != nil {
		return err
	}
	r := Repository{db: database.GetDB(), UserID: userID}
	user, err := r.UpdateAccount(dto.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, api.Response{Msg: "Internal server error. Please try again"})
	}
	return c.JSON(http.StatusOK, api.Response{Msg: "success", Data: user})
}
