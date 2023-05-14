# maintenance-tasks
# Retrieve all tasks
GET localhost:3000/tasks

# Create a task
POST localhost:3000/task
Body -> x-www-form-urlencoded
    
    summary:Test Testtt↵
    techId:2
    role:technician

# Retrieve a task
GET localhost:3000/task?taskID=1&techID=2

# Update a task  
PUT localhost:3000/task
```json
{
    "id": 1,
    "summary": "Test Te34313223sttt\n",
    "performed_date": "2023-05-14T00:00:00Z",
    "technician_id": 2,
    "role": "technician"
}
```

# Delete a task
DELETE localhost:3000/task?id=1
