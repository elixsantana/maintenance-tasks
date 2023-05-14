# maintenance-tasks
# Retrieve all tasks
header: Role: <mananger|technician>
GET localhost:3000/tasks

# Create a task
header: Role: <mananger|technician>
POST localhost:3000/task
Body -> x-www-form-urlencoded
    
    summary:Test Testtt↵
    techId:2
    role:technician

# Retrieve a task
## Manager
header: Role: mananger
GET localhost:3000/task?id=1
## Technician
header: 

Role: technician

TechId: <tech id>

GET localhost:3000/task?id=1

# Update a task  
header: Role: <mananger|technician>
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
header: Role: <mananger|technician>
DELETE localhost:3000/task?id=1
