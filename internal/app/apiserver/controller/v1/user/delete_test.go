package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"

	srvv1 "go-web-demo/internal/app/apiserver/service/v1"
)

func TestUserController_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := srvv1.NewMockService(ctrl)
	mockUserSrv := srvv1.NewMockUserSrv(ctrl)
	mockUserSrv.EXPECT().Delete(gomock.Any(), gomock.Eq("admin"), gomock.Any()).Return(nil)
	mockService.EXPECT().Users().Return(mockUserSrv)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("DELETE", "/v1/users/admin", nil)
	c.Params = []gin.Param{{Key: "name", Value: "admin"}}

	type fields struct {
		srv srvv1.Service
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "default",
			fields: fields{
				srv: mockService,
			},
			args: args{
				c: c,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Controller{
				srv: tt.fields.srv,
			}
			u.Delete(tt.args.c)
		})
	}
}
