package storage

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"gotest.tools/assert"
)

func TestConnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockDatabase(ctrl)
	mockMysqlMetadata := Create(&MysqlConfig{})

	mockDB.EXPECT().Ping().Return(nil).AnyTimes()
	mockDB.EXPECT().Close().Return(nil).AnyTimes()
	mockDB.EXPECT().Exec(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mockDB.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mockDB.EXPECT().Prepare(gomock.Any()).Return(nil, nil).AnyTimes()

	err := mockMysqlMetadata.Connect()
	assert.NilError(t, err)
}

// func TestGetTask(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockDB := NewMockDatabase(ctrl)
// 	mockMysqlMetadata := Create(&MysqlConfig{})
// 	mockMysqlMetadata.db = mockDB

// 	taskID := 1
// 	techID := 1
// 	manager := false

// 	mockDB.EXPECT().Prepare(gomock.Any()).Return(nil, nil).AnyTimes()
// 	mockDB.EXPECT().Query(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

// 	_, err := mockMysqlMetadata.GetTask(taskID, techID, manager)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}
// }
