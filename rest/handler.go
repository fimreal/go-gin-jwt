package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fimreal/go-gin-jwt/database"
	"github.com/fimreal/go-gin-jwt/database/dblayer"
	"github.com/fimreal/go-gin-jwt/database/models"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type Handler struct {
	db dblayer.DBLayer
}

type HandlerInterface interface {
	AddUser(c *gin.Context)
	Login(c *gin.Context)
	IsValid(c *gin.Context)
	RenewToken(c *gin.Context)
}

func NewHandler() (HandlerInterface, error) {
	db, err := database.NewORM(database.DSN)
	if err != nil {
		return nil, err
	}
	return &Handler{
		db: db,
	}, nil
}

// @BasePath /account

// Go API godoc
// @Summary User signin
// @Schemes
// @Description add user
// @Tags Account
// @Accept json
// @Produce json
// @Param data body models.AddUserData true "新录入账号"
// @Success 200 {object} models.LoginResult
// @Router /signin [post]
func (h *Handler) AddUser(c *gin.Context) {
	var userData models.AddUserData
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	user.Username = userData.Username
	user.Password = EncodePassword(userData.Password)
	user.Email = NullStr2Str(userData.Email)
	user, err := h.db.AddUser(user)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "该邮箱已注册 " + userData.Email})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户 " + userData.Username + "遇到问题"})
		return
	}

	userDesc, err := h.db.GetUser(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.UserID = userDesc.UserID

	accessToken, refreshToken, err := CreateTokens(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	LoginResult := models.LoginResult{
		UserID:       user.UserID,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, LoginResult)
}

// Go API godoc
// @Summary User login
// @Schemes
// @Description user login with authrization
// @Tags Account
// @Accept json
// @Produce json
// @Param data body models.LoginRequest true "登入账号"
// @Success 200 {object} models.LoginResult "access token & refresh token"
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	var err error
	var loginRequest models.LoginRequest
	if err = c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// if loginRequest.Username != "" {
	// 	user, err = h.db.GetUser(loginRequest.Username)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户名或者密码错误"})
	// 		return
	// 	}
	// } else if loginRequest.UserID != 0 {
	// 	user, err = h.db.GetUser(loginRequest.UserID)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户名或者密码错误"})
	// 		return
	// 	}
	// }
	user, err = h.db.GetUser(loginRequest.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户名或者密码错误"})
		return
	}

	encodedPassword := EncodePassword(loginRequest.Password)
	if user.Password == encodedPassword {
		// do login
		accessToken, refreshToken, err := CreateTokens(user.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		loginResult := models.LoginResult{
			UserID:       user.UserID,
			Username:     user.Username,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		c.JSON(http.StatusOK, loginResult)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户名或者密码错误"})
	}

}

// Go API godoc
// @Summary check access_token
// @Schemes
// @Description check access_token
// @Tags Account
// @Accept json
// @Produce json
// @Param data body models.LoginRequest true "access_token"
// @Success 200 {object} bool "access_token 是否有效"
// @Router /valid [post]
func (h *Handler) IsValid(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "access_token 不合法"})
	}
	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "access_token 不合法"})
	}

	accessToken := data["access_token"]
	accessSecret := "access" + os.Getenv("SECRET")
	valid, _, err := DecodeToken(accessToken, accessSecret)

	if valid {
		c.JSON(http.StatusOK, valid)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "access_token 不合法"})
		fmt.Println(err)
	}
}

type refreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

// Go API godoc
// @Summary renew access_token
// @Schemes
// @Description renew access_token with refresh_token
// @Tags Account
// @Accept json
// @Produce json
// @Param data body refreshToken true "refresh_token"
// @Success 200 {object} string "refresh_token"
// @Router /renew [post]
func (h *Handler) RenewToken(c *gin.Context) {
	// 解析请求
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "access_token 不合法"})
	}

	// 解码 refreshToken
	var data map[string]string
	json.Unmarshal([]byte(body), &data)
	refreshToken := data["refresh_token"]
	refreshSecret := "refresh" + Secret
	valid, atClaims, _ := DecodeToken(refreshToken, refreshSecret)

	userId, _ := atClaims["user_id"].(int)

	if valid {
		accessSecret := "access" + Secret
		exp := time.Now().Add(ExpTime).Unix()
		accessToken, _ := CreateToken(userId, exp, accessSecret)
		c.JSON(http.StatusOK, accessToken)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error:": "无效的 refresh_token"})
	}
}
