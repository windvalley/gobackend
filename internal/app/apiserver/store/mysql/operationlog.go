package mysql

import (
	"context"

	gorm "gorm.io/gorm"

	"gobackend/pkg/errors"
	"gobackend/pkg/fields"
	metav1 "gobackend/pkg/meta/v1"
	"gobackend/pkg/util/gormtool"

	"gobackend/internal/pkg/code"
	"gobackend/internal/pkg/entity/apiserver/operationlog"
)

type operationLogs struct {
	db *gorm.DB
}

func newOperationLogs(ds *datastore) *operationLogs {
	return &operationLogs{db: ds.db}
}

// Create creates a new OperationLog.
func (o *operationLogs) Create(
	ctx context.Context,
	operationLog *operationlog.OperationLog,
	opts metav1.CreateOptions,
) error {
	return o.db.Create(&operationLog).Error
}

// Delete an OperationLog record.
func (o *operationLogs) Delete(
	ctx context.Context,
	id string,
	opts metav1.DeleteOptions,
) error {
	if opts.Unscoped {
		o.db = o.db.Unscoped()
	}

	err := o.db.Where("id = ?", id).Delete(&operationlog.OperationLog{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

// List OperationLog records.
func (o *operationLogs) List(
	ctx context.Context,
	opts metav1.ListOptions,
) (*operationlog.List, error) {
	ret := &operationlog.List{}
	ol := gormtool.Unpointer(opts.Offset, opts.Limit)

	var (
		where    string
		selector fields.Selector
		err      error
	)

	// opt.FieldSelector e.g.:
	// https://.../?field_selector=req_method==PUT,req_path=/users
	// == means exact match, and = means fuzzy match.
	selector, err = fields.ParseSelector(opts.FieldSelector)
	if err != nil {
		return nil, err
	}

	for _, require := range selector.Requirements() {
		switch require.Field {
		case "req_method":
			where, err = buildWhere(require, where)
		case "req_path":
			where, err = buildWhere(require, where)
		case "http_status":
			where, err = buildWhere(require, where)
		}
	}

	if err != nil {
		return nil, err
	}

	r := o.db.Where(where).
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount)

	return ret, r.Error
}
