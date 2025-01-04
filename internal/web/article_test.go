package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	svcmocks "github.com/chenmuyao/go-bootcamp/internal/service/mocks"
	ijwt "github.com/chenmuyao/go-bootcamp/internal/web/jwt"
	"github.com/chenmuyao/go-bootcamp/pkg/ginx"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) service.ArticleService
		reqBody string

		wantCode int
		wantRes  ginx.Result
	}{
		{
			name: "Create new article and publish",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "my title",
					Content: "my content",
					Author: domain.Author{
						ID: 123,
					},
				}).Return(int64(1), nil)
				return svc
			},
			reqBody: `
{
    "title": "my title",
    "content": "my content"
}`,
			wantCode: http.StatusOK,
			wantRes: ginx.Result{
				Code: ginx.CodeOK,
				Data: float64(1),
			},
		},
		{
			name: "publish existed article",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					ID:      2,
					Title:   "my title",
					Content: "my content",
					Author: domain.Author{
						ID: 123,
					},
				}).Return(int64(2), nil)
				return svc
			},
			reqBody: `
{
    "id": 2,
    "title": "my title",
    "content": "my content"
}`,
			wantCode: http.StatusOK,
			wantRes: ginx.Result{
				Code: ginx.CodeOK,
				Data: float64(2),
			},
		},
		{
			name: "failed to publish",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "my title",
					Content: "my content",
					Author: domain.Author{
						ID: 123,
					},
				}).Return(int64(0), errors.New("error"))
				return svc
			},
			reqBody: `
{
    "title": "my title",
    "content": "my content"
}`,
			wantCode: http.StatusInternalServerError,
			wantRes:  ginx.InternalServerErrorResult,
		},
		{
			name: "Input error",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				return svc
			},
			reqBody: `
{
    "title": "my title",
    "content": "my content",,,
}`,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Prepare
			articleSvc := tc.mock(ctrl)
			hdl := NewArticleHandler(logger.NewZapLogger(zap.L()), articleSvc, nil)

			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("user", ijwt.UserClaims{
					UID: 123,
				})
			})
			hdl.RegisterRoutes(server)

			req, err := http.NewRequest(
				http.MethodPost,
				"/articles/publish",
				bytes.NewBufferString(tc.reqBody))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			// Run Test
			server.ServeHTTP(recorder, req)

			// Check Results
			assert.Equal(t, tc.wantCode, recorder.Code)
			var res ginx.Result
			if tc.wantRes == res {
				// If bind error, don't add wantRes and return ok
				return
			}
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
