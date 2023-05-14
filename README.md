# maintenance-tasks
# Retrieve all tasks
GET localhost:3000/tasks

# Retrieve a task
GET localhost:3000/task?taskID=1&techID=2

# Create a task
POST localhost:3000/task
Body -> x-www-form-urlencoded
    
    summary:Test Testtt↵
    techId:2
    role:technician

# Update a task  
PUT localhost:3000/task
    Body:
    {
            "id": 1,
            "summary": "Test Te34313223sttt\n",
            "performed_date": "2023-05-14T00:00:00Z",
            "technician_id": 2,
            "role": "technician"
    }
