package types

import (
	uuid "github.com/kthomas/go.uuid"
	"github.com/unibrightio/proxy-api/dbutil"
	"github.com/unibrightio/proxy-api/logger"
)

type Organization struct {
	Id               uuid.UUID
	OrganizationName string
}

func (t *Organization) Create() bool {
	if dbutil.Db.GetConn().NewRecord(t) {
		result := dbutil.Db.GetConn().Create(&t)
		rowsAffected := result.RowsAffected
		errors := result.GetErrors()

		if len(errors) > 0 {
			logger.Errorf("errors while creating new entry %v\n", errors)
			return false
		}

		return rowsAffected > 0
	}

	return false
}

func (t *Organization) Delete() bool {
	result := dbutil.Db.GetConn().Delete(&t)
	rowsAffected := result.RowsAffected
	errors := result.GetErrors()

	if len(errors) > 0 {
		logger.Errorf("errors while creating new entry %v\n", errors)
		return false
	}

	return rowsAffected > 0
}
