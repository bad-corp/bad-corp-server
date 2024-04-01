package main

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strconv"
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
	user := User{
		Account:  body.Account,
		Password: body.Password,
	}
	result := ctl.db.Create(&user)
	if result.Error != nil {
		return echo.ErrInternalServerError
	}

	// 登录
	res, err := loginService(user.Id)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, res)
}

func getUserId(c echo.Context) (uint64, bool) {
	if c.Get("userId") == nil {
		return 0, false
	}
	return uint64(c.Get("userId").(uint64)), true
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
	res, err := loginService(user.Id)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, res)
}

func (ctl *controller) createSubject(c echo.Context) error {
	var userId, ok = getUserId(c)
	if !ok {
		return echo.ErrUnauthorized
	}

	var body CreateSubjectDTO
	err := c.Bind(&body)
	if err != nil {
		return echo.ErrBadRequest
	}

	var subject = Subject{
		Name:      body.Name,
		CreatedBy: userId,
	}
	result := ctl.db.Create(&subject)
	if result.Error != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, subject)
}

func (ctl *controller) getSubjectByRandom(c echo.Context) error {
	var userId, ok = getUserId(c)
	if !ok {
		return echo.ErrUnauthorized
	}

	bloom := getInitBloom(ctl.db, userId)
	done := false
	offset := 0
	for !done {
		var subjects []Subject
		result := ctl.db.Where("created_by <> ?", userId).Order("id desc").Limit(100).Offset(offset).Find(&subjects)
		if len(subjects) == 0 {
			return c.JSON(http.StatusNoContent, nil)
		}
		if result.Error != nil {
			return echo.ErrNotFound
		}
		for _, subject := range subjects {
			var sId = uint64ToBytes(uint64(subject.Id))
			var read = bloom.Test(sId)
			if !read {
				bloom.Add(sId)
				// saveDB
				upsertBloom(ctl.db, bloom, userId)
				done = true
				return c.JSON(http.StatusOK, subject)
			}
		}
		offset += 100
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (ctl *controller) createSubjectComment(c echo.Context) error {
	var userId, ok = getUserId(c)
	if !ok {
		return echo.ErrUnauthorized
	}

	var body CreateCommentDTO
	err := c.Bind(&body)
	if err != nil {
		return echo.ErrBadRequest
	}
	subjectId, err := strconv.ParseInt(c.Param("subjectId"), 10, 64)
	if err != nil || subjectId <= 0 {
		return echo.ErrBadRequest
	}
	// TODO: 验证subjectId是否存在

	var subject Subject
	result := ctl.db.Where("id = ?", subjectId).First(&subject)
	if result.RowsAffected == 0 {
		return echo.ErrNotFound
	}

	var comment SubjectComment
	result2 := ctl.db.Where(&SubjectComment{UserId: userId, SubjectIdCreator: subject.CreatedBy}).First(&comment)
	if result2.RowsAffected > 0 {
		var t1 = comment.C + 1
		var t2 = int32(comment.Score) * comment.C
		var t3 = t2 + int32(body.Score)
		comment.Score = int8(t3 / t1)
		comment.C += 1
		ctl.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "subject_id_creator"}},
			DoUpdates: clause.AssignmentColumns([]string{"score", "c"}),
		}).Create(&comment)
	} else {
		result3 := ctl.db.Create(&SubjectComment{
			UserId:           userId,
			SubjectId:        subjectId,
			SubjectIdCreator: subject.CreatedBy,
			Score:            body.Score,
			C:                1,
		})
		if result3.Error != nil {
			return echo.ErrInternalServerError
		}
	}

	calcAndSave(ctl.db)

	return c.JSON(http.StatusOK, comment)
}

type SubjectHot struct {
	SubjectId int64 `json:"subject_id"`
	Score     int64 `json:"score"`
}

func (ctl *controller) getTrendingKings(c echo.Context) error {
	var userScores []UserScore
	result := ctl.db.Order("score desc").Limit(100).Find(&userScores)
	if result.Error != nil {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, userScores)
}

func (ctl *controller) getTrendingQueens(c echo.Context) error {
	var userScores []UserScore
	result := ctl.db.Order("score desc").Limit(100).Find(&userScores)
	if result.Error != nil {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, userScores)
}

func (ctl *controller) getTrendingSubjects(c echo.Context) error {
	var res = []SubjectHot{
		{
			SubjectId: 1,
			Score:     120,
		},
	}
	return c.JSON(http.StatusOK, res)
}

func (ctl *controller) getUser(c echo.Context) error {
	var userId, ok = getUserId(c)
	if !ok {
		return echo.ErrUnauthorized
	}

	var user User
	result := ctl.db.Where(&User{Id: userId}).First(&user)
	if result.Error != nil {
		return echo.ErrNotFound
	}

	var (
		province = Resource{
			Id: uint64(user.Province),
		}
		city = Resource{
			Id: uint64(user.City),
		}
		district = Resource{
			Id: uint64(user.District),
		}
	)
	if user.District != 0 {
		var areaInfo = getAreaInfo(user.Province, user.City, user.District)
		province.Name = areaInfo[0]
		city.Name = areaInfo[1]
		district.Name = areaInfo[2]
	}

	return c.JSON(http.StatusOK, UserDTO{
		Id:        user.Id,
		Account:   user.Account,
		Gender:    user.Gender,
		Name:      user.Name,
		Province:  province,
		City:      city,
		District:  district,
		CreatedAt: user.CreatedAt,
	})
}
