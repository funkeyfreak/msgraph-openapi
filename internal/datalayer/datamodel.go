package datalayer

import (
	"container/list"
	"errors"
)

type DataModelType string

const (
	DataModelTypeByteArray DataModelType = "bytearray"
	DataModelTypeList      DataModelType = "listList"
)

var (
	ErrUnhandledDataModelType = errors.New("unhandled data model type")
)

func (dmt DataModelType) IsValid() bool {
	switch dmt {
	case DataModelTypeByteArray, DataModelTypeList:
		return true
	}
	return false
}

type DataModel struct {
	model map[DataModelType]interface{}
}

func NewDataModel() DataModel {
	return DataModel{
		model: make(map[DataModelType]interface{}),
	}
}

func (dm DataModel) AddDataModel(data interface{}) bool {
	var dmt DataModelType

	switch data.(type) {
	case *list.List:
		dmt = DataModelTypeList
	case []byte:
		dmt = DataModelTypeByteArray
	default:
		return false
	}

	if _, ok := dm.model[dmt]; !ok {
		switch v := data.(type) {
		case *list.List, []byte:
			dm.model[dmt] = v
			return true
		}
	}
	return false
}

func (dm DataModel) RemoveModel(dmt DataModelType) bool {
	if !dmt.IsValid() {
		return false
	}

	if _, ok := dm.model[dmt]; ok {
		delete(dm.model, dmt)
		return true
	}

	return false
}

func (dm DataModel) FetchListModel() *list.List {
	if t, ok := dm.model[DataModelTypeList]; ok {
		switch v := t.(type) {
		case *list.List:
			return v
		}
	}
	return nil
}

func (dm DataModel) FetchByteArrayModel() []byte {
	if t, ok := dm.model[DataModelTypeByteArray]; ok {
		switch v := t.(type) {
		case []byte:
			return v
		}
	}
	return nil
}
