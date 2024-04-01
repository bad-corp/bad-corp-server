package main

import (
	"encoding/json"
	"fmt"
	"github.com/glebarez/sqlite"
	_ "github.com/glebarez/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	db *gorm.DB
}

const (
	user1 = "018e988376b379bfb86ec9e1cf0b8ee1"
)

func (s *Suite) SetupTest() {
	s.db, _ = gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	_ = s.db.AutoMigrate(&User{}, &Corp{}, &CorpComment{}, &CorpScore{})

	s.db.Create(&User{
		Account:  user1,
		Password: "123456",
	})
	s.db.Create(&[]Corp{
		{
			Id:   1,
			Name: "苹果公司",
		},
		{
			Id:   2,
			Name: "苹公司",
		},
		{
			Id:   3,
			Name: "微软苹果公司",
		},
	})
	s.db.Create(&[]CorpScore{
		{
			CorpId: 2,
			Score:  1200,
		},
		{
			CorpId: 1,
			Score:  1100,
		},
	})
}

func TestBaseSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestSignUp() {
	const body = `{"password":"123456"}`

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
		assert.Equal(s.T(), 32, len(res.Account))
		assert.Equal(s.T(), "ey", res.Token[0:2])
	}
}

func (s *Suite) TestSignIn() {
	var body = fmt.Sprintf(`{"account":"%s","password":"123456"}`, user1)

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
		assert.Equal(s.T(), 32, len(res.Account))
		assert.Equal(s.T(), "ey", res.Token[0:2])
	}
}

func (s *Suite) TestSearchCorp() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/corps/search?name=苹", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctl := &controller{db: s.db}

	if assert.NoError(s.T(), ctl.searchCorp(c)) {
		var corps []Corp
		_ = json.Unmarshal(rec.Body.Bytes(), &corps)
		assert.Equal(s.T(), http.StatusOK, rec.Code)
		assert.Equal(s.T(), 2, len(corps))
	}
}

func (s *Suite) TestCreateCorpComment() {
	const body = `{"score":6,"score2":4,"score3":5}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userId", user1)
	c.SetPath("/corps/:corp_id/comments")
	c.SetParamNames("corp_id")
	c.SetParamValues("2")
	ctl := &controller{db: s.db}

	if assert.NoError(s.T(), ctl.createCorpComment(c)) {
		assert.Equal(s.T(), http.StatusOK, rec.Code)
		assert.Equal(s.T(), "{\"id\":1,\"user_id\":\"018e988376b379bfb86ec9e1cf0b8ee1\",\"corp_id\":2,\"score\":6,\"c\":1,\"score2\":4,\"c2\":1,\"score3\":5,\"c3\":1}\n", rec.Body.String())
	}
}

func (s *Suite) TestGetTrendingCorps() {
	const expected = `[{"corp_id":2,"score":1200},{"corp_id":1,"score":1100}]` + "\n"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctl := &controller{db: s.db}

	if assert.NoError(s.T(), ctl.getTrendingCorps(c)) {
		assert.Equal(s.T(), http.StatusOK, rec.Code)
		assert.Equal(s.T(), expected, rec.Body.String())
	}
}

func (s *Suite) TestGetUser() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userId", user1)
	ctl := &controller{db: s.db}

	if assert.NoError(s.T(), ctl.getUser(c)) {
		var user UserDTO
		_ = json.Unmarshal(rec.Body.Bytes(), &user)
		assert.Equal(s.T(), http.StatusOK, rec.Code)
		assert.Equal(s.T(), user1, user.Account)
		assert.IsType(s.T(), time.Time{}, user.CreatedAt)
	}
}
