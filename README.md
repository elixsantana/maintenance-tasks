# Maintenance-tasks
Develop a software to account for maintenance tasks performed during a working day. This application has two types of users (Manager, Technician).
The technician performs tasks and is only able to see, create or update his own performed tasks.
The manager can see tasks from all the technicians, delete them, and should be notified when some tech performs a task.
A task has a summary (max: 2500 characters) and a date when it was performed, the summary from the task can contain personal information.
### Setup
1) Clone https://github.com/elixsantana/maintenance-tasks.git
2) cd to root of project path/maintenance-tasks
4) Run mysql service
5) ```go run main.go```
3) Make a HTTP request like the ones explain below

- Setup with microservice below

# API Examples
### Retrieve all tasks
header: Role: <manager|technician>
GET localhost:3000/tasks

### Create a task
header: Role: <manager|technician>
POST localhost:3000/task
```json
{
    "summary": "Test Te34313223sttt\n",
    "technician_id": 2,
    "role": "technician"
}
```
### Retrieve a task
-  Manager
        header Role: manager
GET localhost:3000/task?id=1
-  Technician
header: 
Role: technician
TechId: <techId>
GET localhost:3000/task?id=1

### Update a task  
header: Role: <manager|technician>
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

### Delete a task
header: Role: <manager|technician>
DELETE localhost:3000/task?id=1

# Setup with Kubernetes and Docker (not finished)
Installations
1. install golang
2. install docker desktop
3. install minikube
4. For Windows users: install git bash

Run the following commands in Bash:
1. ```go mod init project_test```
2. ```go mod tidy```
3. ```docker build -t maintenance-deployment:latest .```
4. ```docker build -t maintenance-deployment:latest .```
5. ```eval $(minikube docker-env)```
6. ```docker save -o maintenance-deployment.tar maintenance-deployment:latest```
7. ```docker load -i maintenance-deployment.tar```
8. ```kubectl create namespace maintenance```
9. ```kubectl apply -f config/deployment.yaml --namespace maintenance ```     
10. ```kubectl apply -f config/service.yaml --namespace maintenance```
11. ```kubectl apply -f config/mysql-deployment.yaml --namespace maintenance```

# Improvements

1. More unit tests
2. Hash summary information to protect sensitive data and not persist plain text in database
3. Implement a message broker with either RabbitMQ or Redis for the notification system
4. Finish up the Kubernetes setup
5. Implement an init container to create the database and tables instead of doing it through the code
6. Implement authentication logic. Right now, I am assuming the header of the HTTP request contains the role and techID. This is not safe.
7. Create script to automatically run all commands for the Kubernetes setup
8. Better logging