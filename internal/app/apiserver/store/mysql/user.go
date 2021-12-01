package mysql

import (
	"context"
	"fmt"

	gorm "gorm.io/gorm"

	"gobackend/pkg/errors"
	"gobackend/pkg/fields"
	metav1 "gobackend/pkg/meta/v1"
	"gobackend/pkg/util/gormtool"

	"gobackend/internal/pkg/code"
	v1 "gobackend/internal/pkg/entity/apiserver/v1"
)

type users struct {
	db *gorm.DB
}

func newUsers(ds *datastore) *users {
	return &users{db: ds.db}
}

// Create creates a new user account.
func (u *users) Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error {
	return u.db.Create(&user).Error
}

// Update updates an user account information.
func (u *users) Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error {
	return u.db.Save(user).Error
}

// Delete deletes the user by the user identifier.
func (u *users) Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error {
	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}

	err := u.db.Where("name = ?", username).Delete(&v1.User{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

// DeleteCollection batch deletes the users.
func (u *users) DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
	if opts.Unscoped {
		u.db = u.db.Unscoped()
	}

	return u.db.Where("name in (?)", usernames).Delete(&v1.User{}).Error
}

// Get return an user by the user identifier.
func (u *users) Get(ctx context.Context, username string, opts metav1.GetOptions) (*v1.User, error) {
	user := &v1.User{}
	err := u.db.Where("name = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrUserNotFound, err.Error())
		}

		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return user, nil
}

// List users.
func (u *users) List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
	ret := &v1.UserList{}
	ol := gormtool.Unpointer(opts.Offset, opts.Limit)

	var (
		where string
		err   error
	)

	// opt.FieldSelector e.g.:
	// https://.../?field_selector=name==levin,email=n@gmail.com
	// == means exact match, and = means fuzzy match.
	selector, _ := fields.ParseSelector(opts.FieldSelector)

	for _, require := range selector.Requirements() {
		switch require.Field {
		case "name":
			where, err = buildWhere(require, where)
		case "email":
			where, err = buildWhere(require, where)
		}
	}

	if err != nil {
		return nil, err
	}

	d := u.db.Where(where).
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount)

	return ret, d.Error
}

func buildWhere(require fields.Requirement, where string) (string, error) {
	if where != "" {
		where += " and "
	}

	switch require.Operator {
	case "==":
		where += fmt.Sprintf("%s = '%v'", require.Field, require.Value)
	case "=":
		where += fmt.Sprintf("%s like '%%%v%%'", require.Field, require.Value)
	case "!=":
		where += fmt.Sprintf("%s %s '%v'", require.Field, require.Operator, require.Value)
	default:
		return "", fmt.Errorf("unknown operator '%s'", require.Operator)
	}

	return where, nil
}
