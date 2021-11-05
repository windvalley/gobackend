package user

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	srvv1 "go-web-backend/internal/app/apiserver/service/v1"
	"go-web-backend/internal/app/apiserver/store"
)

func TestNewUserController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFactory := store.NewMockFactory(ctrl)

	type args struct {
		store store.Factory
	}
	tests := []struct {
		name string
		args args
		want *Controller
	}{
		{
			name: "default",
			args: args{
				store: mockFactory,
			},
			want: &Controller{
				srv: srvv1.NewService(mockFactory),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserController(tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserController() = %v, want %v", got, tt.want)
			}
		})
	}
}
