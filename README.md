# Maintenance-tasks
Develop a software to account for maintenance tasks performed during a working day. This application has two types of users (Manager, Technician).
The technician performs tasks and is only able to see, create or update his own performed tasks.
The manager can see tasks from all the technicians, delete them, and should be notified when some tech performs a task.
A task has a summary (max: 2500 characters) and a date when it was performed, the summary from the task can contain personal information.

### Setup
1) Clone https://github.com/elixsantana/maintenance-tasks.git
2) cd /path/to/project  i.e. ~/maintenance-tasks
4) Run mysql service listening to port 3306
5) ```go run main.go```
6) Use localhost on port 3000 for HTTP requests
7) 'Role' and/or 'TechId' header key-values are needed (This is mocking an authentication and authorization system. This is not safe for production environments)
8) HTTP requests examples below


- Setup with microservice below

# API Examples
### Retrieve all tasks
GET localhost:3000/tasks

Header requirements: Role

    Role: manager


### Retrieve a task
GET localhost:3000/task?id=1

Header requirements: Role, TechId (only for technician)

    Role: manager OR Role: technician

    TechId: 1 (Only for technician)

QueryParams requirement: id (task id)

### Create a task
POST localhost:3000/task

Header requirements: Role

    Role: technician

Body:
```json
{
    "summary": "Test Te34313223sttt\n",
    "technician_id": 2,
    "role": "technician"
}
```


### Update a task  
PUT localhost:3000/task

Header requirements: Role and TechId

    Role: technician
    TechId: 2

Body:
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
DELETE localhost:3000/task?id=1

Header requirements: Role

    Role: manager

# Setup with Kubernetes and Docker (not finished)
Installations
1. Install golang
2. Install docker desktop
3. Install minikube
4. For Windows users: install git bash

Run the following commands in Bash:
1. cd /path/to/project
1. ```go mod init maintenance-tasks```
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
6. Implement authentication logic. Right now, I am assuming the header of the HTTP request contains the Role and TechId. This is not safe.
7. Create script to automatically run all commands for the Kubernetes setup
8. Better logging