package main

import (
	"encoding/binary"
	"fmt"
	"github.com/glebarez/sqlite"
	_ "github.com/glebarez/sqlite"
	"github.com/stretchr/testify/suite"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type BussSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *BussSuite) SetupTest() {
	s.db, _ = gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	_ = s.db.AutoMigrate(&User{}, &Subject{}, &SubjectComment{}, &UserScore{})
}

func TestBussSuite(t *testing.T) {
	suite.Run(t, new(BussSuite))
}

func testPost(e *echo.Echo, body string) echo.Context {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}

func auth(c echo.Context, userId int64) {
	c.Set("userId", userId)
}

func (s *BussSuite) Test() {
	e := echo.New()
	ctl := &controller{db: s.db}

	// 生成用户
	for i := 1; i <= 100; i++ {
		body := fmt.Sprintf(`{"account":"biubiu%s","password":"123456"}`, strconv.Itoa(i))
		c := testPost(e, body)

		err := ctl.signUp(c)
		if err != nil {
			assert.FailNow(s.T(), err.Error())
		}
	}

	// 1~3号用户发表一篇subject
	for i := 1; i <= 3; i++ {
		body := fmt.Sprintf(`{"name":"iphone%s"}`, strconv.Itoa(i))
		c := testPost(e, body)

		auth(c, int64(i))
		err := ctl.createSubject(c)
		if err != nil {
			assert.FailNow(s.T(), err.Error())
		}
	}

	// 4~50号用户 对 2号用户 吐槽
	for i := 4; i <= 50; i++ {
		body := fmt.Sprintf(`{"score":%d}`, rand.Intn(9)+1)
		c := testPost(e, body)
		c.SetPath("/subjects/:subjectId/comments")
		c.SetParamNames("subjectId")
		c.SetParamValues(strconv.Itoa(2))

		auth(c, int64(i))
		err := ctl.createSubjectComment(c)
		if err != nil {
			assert.FailNow(s.T(), err.Error())
		}
	}

	// 51~80号用户 对 12号用户 吐槽
	for i := 51; i <= 80; i++ {
		body := fmt.Sprintf(`{"score":%d}`, rand.Intn(9)+1)
		c := testPost(e, body)
		c.SetPath("/subjects/:subjectId/comments")
		c.SetParamNames("subjectId")
		c.SetParamValues(strconv.Itoa(12))

		auth(c, int64(i))
		err := ctl.createSubjectComment(c)
		if err != nil {
			assert.FailNow(s.T(), err.Error())
		}
	}

	// 81~98号用户 对 67号用户 吐槽
	for i := 81; i <= 98; i++ {
		body := fmt.Sprintf(`{"score":%d}`, rand.Intn(9)+1)
		c := testPost(e, body)
		c.SetPath("/subjects/:subjectId/comments")
		c.SetParamNames("subjectId")
		c.SetParamValues(strconv.Itoa(67))

		auth(c, int64(i))
		err := ctl.createSubjectComment(c)
		if err != nil {
			assert.FailNow(s.T(), err.Error())
		}
	}

	// 99~100号用户 对 83号用户 吐槽
	for i := 99; i <= 100; i++ {
		body := fmt.Sprintf(`{"score":%d}`, rand.Intn(9)+1)
		c := testPost(e, body)
		c.SetPath("/subjects/:subjectId/comments")
		c.SetParamNames("subjectId")
		c.SetParamValues(strconv.Itoa(83))

		auth(c, int64(i))
		err := ctl.createSubjectComment(c)
		if err != nil {
			assert.FailNow(s.T(), err.Error())
		}
	}

	// 事实计算
	var userScore UserScore
	s.db.Table("user_score").Where("user_id = ?", 2).Find(&userScore)

	// 直接调用底层计算
	var comments []SubjectComment
	s.db.Find(&comments)
	var users []User
	s.db.Table("user").Find(&users)
	m := calcScore(&users, &comments)

	assert.Equal(s.T(), userScore.Score, int64(m[2]*10000))
}

func (s *BussSuite) TestBaseGraph() {
	dg := simple.NewDirectedGraph()
	// 假设我们有一些唯一标识符作为节点
	node1 := graph.Node(simple.Node(1))
	node2 := graph.Node(simple.Node(2))
	node3 := graph.Node(simple.Node(3))
	node4 := graph.Node(simple.Node(4))
	node5 := graph.Node(simple.Node(5))
	node6 := graph.Node(simple.Node(6))
	node7 := graph.Node(simple.Node(7))
	node8 := graph.Node(simple.Node(8))

	// 将节点添加到图中
	dg.AddNode(node1)
	dg.AddNode(node2)
	dg.AddNode(node3)
	dg.AddNode(node4)
	dg.AddNode(node5)
	dg.AddNode(node6)
	dg.AddNode(node7)
	dg.AddNode(node8)

	// 创建一条边，并指定两个端点
	edge1 := ScoreEdge{F: node4, T: node2}
	edge2 := ScoreEdge{F: node3, T: node2}
	edge3 := ScoreEdge{F: node2, T: node1}
	edge4 := ScoreEdge{F: node5, T: node1}
	edge5 := ScoreEdge{F: node6, T: node5}
	edge6 := ScoreEdge{F: node1, T: node7}
	edge7 := ScoreEdge{F: node7, T: node8}

	// 将边添加到图中
	dg.SetEdge(edge1)
	dg.SetEdge(edge2)
	dg.SetEdge(edge3)
	dg.SetEdge(edge4)
	dg.SetEdge(edge5)
	dg.SetEdge(edge6)
	dg.SetEdge(edge7)

	var asdd = toFullGraph(dg, node1.ID())

	var expected = "[{FromId:1 ToId:7 Deep:2 Score:0} {FromId:7 ToId:8 Deep:1 Score:0} {FromId:2 ToId:1 Deep:3 Score:0} {FromId:4 ToId:2 Deep:4 Score:0} {FromId:3 ToId:2 Deep:4 Score:0} {FromId:5 ToId:1 Deep:3 Score:0} {FromId:6 ToId:5 Deep:4 Score:0}]"
	assert.Equal(s.T(), expected, fmt.Sprintf("%+v", asdd))
}

func ssss(i uint64) []byte {
	n1 := make([]byte, 8)
	binary.BigEndian.PutUint64(n1, i)
	return n1
}
