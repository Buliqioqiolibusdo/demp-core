package models

import "github.com/buliqioqiolibusdo/demp-core/interfaces"

type modelRelation struct {
	d       interfaces.Model
	id      interfaces.ModelId
	colName string
}

var ModelRelations = []modelRelation{
	{d: &Job{}, id: interfaces.ModelIdJob, colName: interfaces.ModelColNameJob},
	{d: &Node{}, id: interfaces.ModelIdNode, colName: interfaces.ModelColNameNode},
}
