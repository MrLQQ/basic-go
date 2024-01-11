package integration

import (
	"basic-go/webook/internal/integration/startup"
	"basic-go/webook/internal/repository/dao"
	ijwt "basic-go/webook/internal/web/jwt"
	"bytes"
	"context"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type ArticleMongoDBHandlerSuite struct {
	suite.Suite
	mdb     *mongo.Database
	col     *mongo.Collection
	liveCol *mongo.Collection
	server  *gin.Engine
}

func (s *ArticleMongoDBHandlerSuite) SetupSuite() {
	s.mdb = startup.InitMongoDB()
	s.col = s.mdb.Collection("articles")
	s.liveCol = s.mdb.Collection("publish")
	node, err := snowflake.NewNode(1)
	assert.NoError(s.T(), err)
	hdl := startup.InitArticleHandler(dao.NewMongoDBAritcleDAO(s.mdb, node))
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("user", ijwt.UserClaims{
			Uid: 123,
		})
	})
	hdl.RegisterRoutes(server)
	s.server = server
}

func (s *ArticleMongoDBHandlerSuite) TearDownTest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := s.col.DeleteMany(ctx, bson.D{})
	assert.NoError(s.T(), err)
	_, err = s.liveCol.DeleteMany(ctx, bson.D{})
	assert.NoError(s.T(), err)
}

func (s *ArticleMongoDBHandlerSuite) Test_Edit() {
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
				ctx, canncel := context.WithTimeout(context.Background(), time.Second)
				defer canncel()

				// 验证是否保存到了数据库
				var art dao.Article
				err := s.col.FindOne(ctx, bson.D{bson.E{"author_id", 123}}).Decode(&art)
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				assert.True(t, art.Id > 0)
				art.Ctime = 0
				art.Utime = 0
				art.Id = 0
				assert.Equal(t, dao.Article{
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
					Status:   1,
				}, art)
			},
			art: Article{
				Title:   "我的标题",
				Content: "我的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 1,
			},
		},
		{
			name: "修改帖子",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_, err := s.col.InsertOne(ctx, &dao.Article{
					Id:       11,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
					// 假设这是一个已经发表了的帖子
					Status: 2,
					Ctime:  456,
					Utime:  789,
				})
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				// 验证是否保存到了数据库
				var art dao.Article
				err := s.col.FindOne(ctx, bson.D{bson.E{"id", 11}}).
					Decode(&art)
				assert.NoError(t, err)
				assert.True(t, art.Utime > 789)
				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       11,
					Title:    "新的标题",
					Content:  "新的内容",
					AuthorId: 123,
					Ctime:    456,
					// 更新之后，是为发表状态
					Status: 1,
				}, art)
			},
			art: Article{
				Id:      11,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 2,
			},
		},
		{
			name: "修改帖子 - 修改别人的帖子",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_, err := s.col.InsertOne(ctx, &dao.Article{
					Id:       22,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 1024,
					Status:   2,
					Ctime:    456,
					Utime:    789,
				})
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				// 验证数据没有变
				var art dao.Article
				err := s.col.FindOne(ctx, bson.D{bson.E{"id", 22}}).Decode(&art)
				assert.NoError(t, err)
				assert.Equal(t, dao.Article{
					Id:       22,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 1024,
					Status:   2,
					Ctime:    456,
					Utime:    789,
				}, art)
			},
			art: Article{
				Id:      22,
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
			if tc.wantRes.Data > 0 {
				// 你只能断言有ID
				assert.True(t, res.Data > 0)

			}

		})
	}

}

func TestArticleMongoDBHandler(t *testing.T) {
	suite.Run(t, &ArticleMongoDBHandlerSuite{})
}
