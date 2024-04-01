package main

import (
	"encoding/json"
	"github.com/glebarez/sqlite"
	_ "github.com/glebarez/sqlite"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type Suite struct {
	suite.Suite
	db *gorm.DB
}

func (s *Suite) SetupTest() {
	s.db, _ = gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	_ = s.db.AutoMigrate(&User{}, &Subject{}, &SubjectComment{}, &UserScore{})

	s.db.Create(&User{
		Id:       1,
		Account:  "zhuzhu",
		Password: "123456",
	})
	s.db.Create(&Subject{
		Id:        1,
		Name:      "apple 13",
		CreatedBy: 2,
	})
	s.db.Create(&[]UserScore{
		{
			UserId: 2,
			Score:  1200,
		},
		{
			UserId: 1,
			Score:  1100,
		},
	})
}

func TestBaseSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestSignUp() {
	const body = `{"account":"biubiu","password":"123456"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth/sign_up", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctl := &controller{db: s.db}

	if assert.NoError(s.T(), ctl.signUp(c)) {
		var res AuthToken
		assert.NoError(s.T(), json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(s.T(), http.StatusCreated, rec.Code)
		assert.Equal(s.T(), "ey", res.Token[0:2])
	}
}

func (s *Suite) TestSignIn() {
	const body = `{"account":"zhuzhu","password":"123456"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth/sign_in", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctl := &controller{db: s.db}

	if assert.NoError(s.T(), ctl.signIn(c)) {
		var res AuthToken
		assert.NoError(s.T(), json.Unmarshal(rec.Body.Bytes(), &res))
		assert.Equal(s.T(), http.StatusOK, rec.Code)
		assert.Equal(s.T(), "ey", res.Token[0:2])
	}
}

func (s *Suite) TestCreateSubject() {
	const body = `{"name":"iphone 15p"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userId", int64(1))
	ctl := &controller{db: s.db}

	if assert.NoError(s.T(), ctl.createSubject(c)) {
		assert.Equal(s.T(), http.StatusOK, rec.Code)
		assert.Equal(s.T(), "{\"id\":2,\"name\":\"iphone 15p\",\"created_by\":1}\n", rec.Body.String())
	}
}

func (s *Suite) TestGetSubjectByRandom() {
	const expected = `{"id":1,"name":"apple 13","created_by":2}` + "\n"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userId", int64(1))
	ctl := &controller{db: s.db}

	if assert.NoError(s.T(), ctl.getSubjectByRandom(c)) {
		assert.Equal(s.T(), http.StatusOK, rec.Code)
		assert.Equal(s.T(), expected, rec.Body.String())
	}
}

func (s *Suite) TestCreateSubjectComment() {
	const body = `{"score":1}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userId", int64(1))
	c.SetPath("/subjects/:subjectId/comments")
	c.SetParamNames("subjectId")
	c.SetParamValues("2")
	ctl := &controller{db: s.db}

	if assert.NoError(s.T(), ctl.createSubjectComment(c)) {
		assert.Equal(s.T(), http.StatusOK, rec.Code)
		assert.Equal(s.T(), "{\"id\":1,\"user_id\":1,\"subject_id\":2,\"score\":1}\n", rec.Body.String())
	}
}

func (s *Suite) TestGetTrendingKings() {
	const expected = `[{"user_id":2,"score":1200},{"user_id":1,"score":1100}]` + "\n"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctl := &controller{db: s.db}

	if assert.NoError(s.T(), ctl.getTrendingKings(c)) {
		assert.Equal(s.T(), http.StatusOK, rec.Code)
		assert.Equal(s.T(), expected, rec.Body.String())
	}
}

func (s *Suite) TestGetTrendingQueens() {
	const expected = `[{"user_id":2,"score":1200},{"user_id":1,"score":1100}]` + "\n"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctl := &controller{db: s.db}

	if assert.NoError(s.T(), ctl.getTrendingQueens(c)) {
		assert.Equal(s.T(), http.StatusOK, rec.Code)
		assert.Equal(s.T(), expected, rec.Body.String())
	}
}

func TestGetTrendingSubjects(t *testing.T) {
	const expected = `[{"subject_id":1,"score":120}]` + "\n"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctl := &controller{db: nil}

	if assert.NoError(t, ctl.getTrendingSubjects(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, rec.Body.String())
	}
}
