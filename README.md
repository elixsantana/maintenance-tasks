# Maintenance-tasks
Develop a software to account for maintenance tasks performed during a working day. This application has two types of users (Manager, Technician).
The technician performs tasks and is only able to see, create or update his own performed tasks.
The manager can see tasks from all the technicians, delete them, and should be notified when some tech performs a task.
A task has a summary (max: 2500 characters) and a date when it was performed, the summary from the task can contain personal information.

# Setup 

### Standalone (Wihtout K8s)
1) Clone https://github.com/elixsantana/maintenance-tasks.git
2) cd ```/path/to/project```  i.e. ~/maintenance-tasks
3) Download MySQL 5.7 or above and configure the root password and set it to ```test```
4) Start the MySQL server (port 3306).
5) ```go run main.go```
6) Use localhost on port 3000 for HTTP requests.
7) Since Auth was not a requirement, ```Role``` and/or ```TechId``` headers key-values are needed (This is playing the role of an authentication and authorization system. This is not safe for production environments).
8) HTTP requests examples below.

### Containerized. Setup with Kubernetes and Docker
Installations
1. Install golang
2. Install docker desktop
3. Install minikube
4. For Windows users: install git bash

Run the following commands in Bash:
1. Open Docker Desktop and wait for Docker Engine to start
2. cd /path/to/project
3. ```minikube start```
4. ```go mod init maintenance-tasks```
5. ```go mod tidy```
6. ```docker build -t maintenance-deployment:latest .```
7. ```docker build -t maintenance-deployment:latest .```
8. ```eval $(minikube docker-env)```
9. ```docker save -o maintenance-deployment.tar maintenance-deployment:latest``` OR (docker images -q maintenance-deployment:latest | xargs docker save | docker load)
10. ```docker load -i maintenance-deployment.tar``` (Skip if you ran the command with xargs in step 9) 
11. ```kubectl create namespace maintenance```
12. ```kubectl apply -f config/mysql-deployment.yaml --namespace maintenance```
13. ```kubectl apply -f config/deployment.yaml --namespace maintenance ```

Setting up requests from local machine:
1. ```kubectl config set-context --current --namespace=maintenance```
2. ```kubectl get services```
3. ```kubectl port-forward service/maintenance-deployment 3000:3000```


Now you can make requests from localhost:3000/ to the containarized app


If you need to restart the deployment:
- ```kubectl get deployments```
- ```kubectl rollout restart deployment maintenance-runner```

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

QueryParams requirement: ```id``` (task id)

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

# Improvements

1. More unit tests.
2. Hash summary information to protect sensitive data and not persist plain text in database.
3. Implement a message broker with either RabbitMQ or Redis for the notification system.
4. Implement an init container to create the database and tables instead of doing it through the code.
5. Implement authentication logic. Right now, I am assuming the header of the HTTP request contains the Role and TechId. This is not safe.
6. Create script to automatically run all commands for the Kubernetes setup
7. Load config for Vault or OS env. This will avoid having to rebuild the docker image if the creds are wrong or need to be updated.
8. Better logging.