package integration

import (
	"basic-go/webook/internal/integration/startup"
	"basic-go/webook/internal/repository/dao"
	ijwt "basic-go/webook/internal/web/jwt"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ArticleHandlerSuite struct {
	suite.Suite
	db     *gorm.DB
	server *gin.Engine
}

func (s *ArticleHandlerSuite) SetupSuite() {
	s.db = startup.InitDB()
	hdl := startup.InitArticleHandler()
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("user", ijwt.UserClaims{
			Uid: 123,
		})
	})
	hdl.RegisterRoutes(server)
	s.server = server
}

func (s *ArticleHandlerSuite) TearDownTest() {
	s.db.Exec("truncate table `articles`")
}

func (s *ArticleHandlerSuite) Test_Edit() {
	t := s.T()

	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		// 前端传过来的肯定是一个json
		art Article

		wantCode int
		wantRes  Result[int64]
	}{
		{
			name:   "新建帖子",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				// 验证是否保存到了数据库
				var art dao.Article
				err := s.db.Where("author_id=?", 123).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				assert.True(t, art.Id > 0)
				assert.Equal(t, "我的标题", art.Title)
				assert.Equal(t, "我的内容", art.Content)
				assert.Equal(t, int64(123), art.AuthorId)
			},
			art: Article{
				Title:   "我的标题",
				Content: "我的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64](Result[int]{
				// 我希望ID是1
				Data: 1,
			}),
		},
		{
			name: "修改帖子",
			before: func(t *testing.T) {
				err := s.db.Create(dao.Article{
					Id:       2,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证是否保存到了数据库
				var art dao.Article
				err := s.db.Where("id=?", 2).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Utime > 789)
				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       2,
					Title:    "新的标题",
					Content:  "新的内容",
					AuthorId: 123,
					Ctime:    456,
					Status:   1,
				}, art)
			},
			art: Article{
				Id:      2,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64](Result[int]{
				Data: 2,
			}),
		},
		{
			name: "修改帖子 - 修改别人的帖子",
			before: func(t *testing.T) {
				err := s.db.Create(dao.Article{
					Id:       3,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 234,
					Status:   1,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据没有变
				var art dao.Article
				err := s.db.Where("id=?", 3).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, dao.Article{
					Id:       3,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 234,
					Status:   1,
					Ctime:    456,
					Utime:    789,
				}, art)
			},
			art: Article{
				Id:      3,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg: "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/articles/edit",
				bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			s.server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantCode != http.StatusOK {
				return
			}
			var res Result[int64]
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)

		})
	}

}

func TestArticleHandler(t *testing.T) {
	suite.Run(t, &ArticleHandlerSuite{})
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
