package api

import (
	"maintenance-tasks/storage"
	"sort"
	"time"
)

func convertTime(timeStr string) time.Time {
	var timeLayout string = "2006-01-02T15:04:05"
	time, _ := time.Parse(timeLayout, timeStr)
	return time
}

var mockData = map[int]storage.Task{
	1: {
		ID:             1,
		Summary:        "Test 1 of technician 777\n",
		Performed_date: convertTime("2023-05-14T00:00:00Z"),
		TechnicianID:   777,
		Role:           "technician",
	},
	2: {
		ID:             2,
		Summary:        "Test 2 of technician 777\n",
		Performed_date: convertTime("2023-05-14T00:00:00Z"),
		TechnicianID:   777,
		Role:           "technician",
	},
	3: {
		ID:             3,
		Summary:        "Test 3 of technician 777\n",
		Performed_date: convertTime("2023-05-14T00:00:00Z"),
		TechnicianID:   777,
		Role:           "technician",
	},
	4: {
		ID:             4,
		Summary:        "Test 4 of technician 888\n",
		Performed_date: convertTime("2023-05-14T00:00:00Z"),
		TechnicianID:   888,
		Role:           "technician",
	},
	5: {
		ID:             5,
		Summary:        "Test 5 of technician 888\n",
		Performed_date: convertTime("2023-05-14T00:00:00Z"),
		TechnicianID:   888,
		Role:           "technician",
	},
	6: {
		ID:             6,
		Summary:        "Test 6 of technician 999\n",
		Performed_date: convertTime("2023-05-14T00:00:00Z"),
		TechnicianID:   999,
		Role:           "technician",
	},
}

func GetMockDataSlice(taskId int) []storage.Task {
	var tasks = []storage.Task{}

	if taskId > 0 {
		data := mockData[taskId]
		tasks = append(tasks, data)
	} else {
		tasks = OrderedTasksSlice(mockData)
	}

	return tasks
}

func ModifyTaskMockData(task storage.Task) storage.Task {
	tID := task.ID
	mockData[tID] = task
	return mockData[tID]
}

func OrderedTasksSlice(map[int]storage.Task) []storage.Task {
	var tasks = []storage.Task{}
	keys := make([]int, 0, len(mockData))
	for key := range mockData {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	for _, key := range keys {
		tasks = append(tasks, mockData[key])
	}

	return tasks
}
