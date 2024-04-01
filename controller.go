package main

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strconv"
	"strings"
)

type (
	controller struct {
		db *gorm.DB
	}
)

func (ctl *controller) signUp(c echo.Context) error {
	var body SignUpDTO
	err := c.Bind(&body)
	if err != nil {
		log.Error(err)
		return echo.ErrBadRequest
	}

	// 注册
	u, _ := uuid.NewV7()
	account := strings.ReplaceAll(u.String(), "-", "")
	user := User{
		Account:  account,
		Password: body.Password,
	}
	result := ctl.db.Create(&user)
	if result.Error != nil {
		return echo.ErrInternalServerError
	}

	// 登录
	token, expiredTime, err := loginService(user.Account)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, AuthToken{
		Token:       token,
		ExpiredTime: expiredTime,
		Account:     account,
	})
}

func getUserId(c echo.Context) (string, bool) {
	if c.Get("userId") == nil || c.Get("userId") == "" {
		return "", false
	}
	return c.Get("userId").(string), true
}

func (ctl *controller) signIn(c echo.Context) error {
	var body SignInDTO
	err := c.Bind(&body)
	if err != nil {
		return echo.ErrBadRequest
	}

	var user User
	result := ctl.db.Where(&User{Account: body.Account, Password: body.Password}).First(&user)
	if result.Error != nil {
		return echo.ErrUnauthorized
	}

	// 登录
	token, expiredTime, err := loginService(user.Account)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, AuthToken{
		Token:       token,
		ExpiredTime: expiredTime,
		Account:     user.Account,
	})
}

func (ctl *controller) searchCorp(c echo.Context) error {
	query := c.QueryParams()

	var corps []Corp
	_ = ctl.db.
		Where("name LIKE ?", query.Get("name")+"%").
		Scopes(Paginate(c.Request())).
		Find(&corps)

	return c.JSON(http.StatusOK, corps)
}

func (ctl *controller) getCorp(c echo.Context) error {
	corpId, err := strconv.ParseInt(c.Param("corp_id"), 10, 64)
	if err != nil || corpId <= 0 {
		return echo.ErrBadRequest
	}

	var corp Corp
	_ = ctl.db.
		Where("id = ?").
		First(&corp)

	return c.JSON(http.StatusOK, corp)
}

func (ctl *controller) createCorpComment(c echo.Context) error {
	var userId, ok = getUserId(c)
	if !ok {
		return echo.ErrUnauthorized
	}

	var body CreateCommentDTO
	err := c.Bind(&body)
	if err != nil {
		return echo.ErrBadRequest
	}
	corpId, err := strconv.ParseInt(c.Param("corp_id"), 10, 64)
	if err != nil || corpId <= 0 {
		return echo.ErrBadRequest
	}
	// TODO: 验证 corpId 是否存在

	var corp Corp
	result := ctl.db.Where("id = ?", corpId).First(&corp)
	if result.RowsAffected == 0 {
		return echo.ErrNotFound
	}

	// TODO: 锁
	var comment CorpComment
	result2 := ctl.db.Where(&CorpComment{UserId: userId, CorpId: corp.Id}).First(&comment)
	if result2.RowsAffected > 0 {
		var s, sc = score(comment.Score, body.Score, comment.ScoreCount)
		var s2, sc2 = score(comment.Score2, body.Score2, comment.ScoreCount2)
		var s3, sc3 = score(comment.Score3, body.Score3, comment.ScoreCount3)
		comment.Score = s
		comment.ScoreCount = sc
		comment.Score = s2
		comment.ScoreCount = sc2
		comment.Score = s3
		comment.ScoreCount = sc3
		ctl.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "corp_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"score", "c", "score2", "c2", "score3", "c3"}),
		}).Create(&comment)
	} else {
		comment = CorpComment{
			UserId:      userId,
			CorpId:      corp.Id,
			Score:       body.Score,
			ScoreCount:  1,
			Score2:      body.Score2,
			ScoreCount2: 1,
			Score3:      body.Score3,
			ScoreCount3: 1,
		}
		result3 := ctl.db.Create(&comment)
		if result3.Error != nil {
			return echo.ErrInternalServerError
		}
	}

	calcAndSave(ctl.db)

	return c.JSON(http.StatusOK, comment)
}

func (ctl *controller) getTrendingCorps(c echo.Context) error {
	queryTrendingType := c.QueryParam("type")
	var trendingType string
	if queryTrendingType == "asc" {
		trendingType = queryTrendingType
	} else {
		trendingType = "desc"
	}

	var corpScores []CorpScore
	result := ctl.db.Order("score " + trendingType).Limit(100).Find(&corpScores)
	if result.Error != nil {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, corpScores)
}

func (ctl *controller) getUser(c echo.Context) error {
	var userId, ok = getUserId(c)
	if !ok {
		return echo.ErrUnauthorized
	}

	var user User
	result := ctl.db.Where(&User{Account: userId}).First(&user)
	if result.Error != nil {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, UserDTO{
		Account:   user.Account,
		CreatedAt: user.CreatedAt,
	})
}
