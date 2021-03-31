package datalayer_test

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/funkeyfreak/msgraph-openapi/internal/datalayer"
)

func Test_DataModelType_IsValid(t *testing.T) {
	t.Run("Validates_DataModelTypeByteArray", func(t *testing.T) {

		require.True(t, datalayer.DataModelTypeByteArray.IsValid())
	})

	t.Run("Validates_DataModelTypeList", func(t *testing.T) {
		require.True(t, datalayer.DataModelTypeList.IsValid())
	})

	t.Run("Validates_InvalidDataModelType", func(t *testing.T) {
		require.False(t, InvalidDataModelType.IsValid())
	})
}

func Test_NewDataModel(t *testing.T) {
	dm := datalayer.NewDataModel()
	require.NotNil(t, dm)
	require.IsType(t, datalayer.DataModel{}, dm)
}

func Test_DataModel_AddDataModel(t *testing.T) {
	dm := datalayer.NewDataModel()
	t.Run("Succeeds_WhenAddingDataModelTypeList", func(t *testing.T) {
		require.True(t, dm.AddDataModel(ListDataModel))
	})

	t.Run("Succeeds_WhenAddingDataModelTypeByteArray", func(t *testing.T) {
		require.True(t, dm.AddDataModel(ByteArrayDataModel))
	})

	t.Run("Fails_WhenDataModelTypeList_AlreadyAdded", func(t *testing.T) {
		require.False(t, dm.AddDataModel(ListDataModel))
	})

	t.Run("Fails_WhenDataModelTypeByteArray_AlreadyAdded", func(t *testing.T) {
		require.False(t, dm.AddDataModel(ByteArrayDataModel))
	})

	t.Run("Fails_AttemptingToAddInvalidDataModel", func(t *testing.T) {
		require.False(t, dm.AddDataModel("foo bar"))
	})
}

func Test_DataModel_RemoveModel(t *testing.T) {
	dm := datalayer.NewDataModel()
	dm.AddDataModel(ListDataModel)
	dm.AddDataModel(ByteArrayDataModel)

	t.Run("Succeeds_WhenRemovingListTypeDataModel", func(t *testing.T) {
		require.True(t, dm.RemoveModel(datalayer.DataModelTypeList))
	})

	t.Run("Succeeds_WhenRemovingByteArrayTypeDataModel", func(t *testing.T) {
		require.True(t, dm.RemoveModel(datalayer.DataModelTypeByteArray))
	})

	t.Run("Fails_WhenItemIsAlreadyDeleted", func(t *testing.T) {
		require.False(t, dm.RemoveModel(datalayer.DataModelTypeList))
		require.False(t, dm.RemoveModel(datalayer.DataModelTypeByteArray))
	})

	t.Run("Fails_WhenRemovingAnyTimeWhichIsNotInTheDataModel", func(t *testing.T) {
		dm.AddDataModel(ListDataModel)
		dm.AddDataModel(ByteArrayDataModel)
		require.False(t, dm.RemoveModel(InvalidDataModelType))
	})
}

func Test_DataModel_FetchListModel(t *testing.T) {
	dm := datalayer.NewDataModel()
	dm.AddDataModel(ListDataModel)

	t.Run("Succeeds", func(t *testing.T) {
		data := dm.FetchListModel()
		require.NotNil(t, data)
		require.IsType(t, &list.List{}, data)
	})

	t.Run("Fails_WhenEmpty", func(t *testing.T) {
		dm.RemoveModel(datalayer.DataModelTypeList)
		data := dm.FetchListModel()
		require.Nil(t, data)
	})
}

func Test_DataModel_FetchByteArrayModel(t *testing.T) {
	dm := datalayer.NewDataModel()
	dm.AddDataModel(ByteArrayDataModel)

	t.Run("Succeeds", func(t *testing.T) {
		data := dm.FetchByteArrayModel()
		require.NotNil(t, data)
		require.IsType(t, make([]byte, 0), data)
	})

	t.Run("Fails_WhenEmpty", func(t *testing.T) {
		dm.RemoveModel(datalayer.DataModelTypeByteArray)
		data := dm.FetchByteArrayModel()
		require.Nil(t, data)
	})
}
