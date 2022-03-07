package models

import (
	"github.com/buliqioqiolibusdo/demp-core/interfaces"
	"github.com/buliqioqiolibusdo/demp-core/utils/binders"
)

func GetModelColName(id interfaces.ModelId) (colName string) {
	return binders.NewColNameBinder(id).MustBindString()
}
