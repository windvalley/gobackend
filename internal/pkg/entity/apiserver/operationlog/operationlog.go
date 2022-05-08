package operationlog

import (
	"time"

	metav1 "gobackend/pkg/meta/v1"
)

// OperationLog user operation audit.
type OperationLog struct {
	metav1.ObjectMetaBase `json:"metadata,omitempty"`

	Username   string    `json:"username" gorm:"index;column:username"`
	UserAgent  string    `json:"user_agent" gorm:"column:user_agent"`
	ClientIP   string    `json:"client_ip" gorm:"column:client_ip"`
	ReqMethod  string    `json:"req_method" gorm:"index;column:req_method"`
	ReqPath    string    `json:"req_path" gorm:"index;column:req_path"`
	ReqBody    string    `json:"req_body" gorm:"column:req_body"`
	ReqReferer string    `json:"req_referer" gorm:"column:req_referer"`
	ReqTime    time.Time `json:"req_time" gorm:"column:req_time"`
	ReqLatency float64   `json:"req_latency" gorm:"column:req_latency"`
	HTTPStatus int       `json:"http_status" gorm:"index;column:http_status"`
	ResData    string    `json:"res_data" gorm:"column:res_data"`
}

// List user operation audit list.
type List struct {
	metav1.ListMeta `json:",inline"`

	Items []*OperationLog `json:"items"`
}

// TableName maps to mysql table name.
func (u *OperationLog) TableName() string {
	return "operation_log"
}
