# Step-by-Step Setup Guide

## Step 1: Clone the Repository

First, clone the repository that contains the Docker Compose configuration.

```bash
git clone <repository-url>
cd <repository-directory>
```

## Step 2: Ensure Docker and Docker Compose are Installed

Make sure Docker and Docker Compose are installed on your machine. You can check by running:

```bash
docker --version
docker-compose --version
```

## Step 3: Build and Start the Services

Navigate to the directory containing the docker-compose.yml file and run the following command to build and start all services:

```bash
docker-compose up --build
```

This command will build the Docker images for each service and start them up according to the configuration.

## Step 4: Access the Services

Frontend: Access the frontend service at http://localhost:80.

Mailhog: Access the Mailhog web interface at http://localhost:8025.

Postgres: Connect to the Postgres database using a database client on localhost:5432.

MongoDB: Connect to the MongoDB database using a database client on localhost:27017.


## Step 5: Monitoring and Logs
You can monitor the logs for each service using:

```bash
docker-compose logs -f <service-name>
```

Replace <service-name> with the name of the service you want to monitor (e.g., frontend-service, broker-service).

### Step 6: Stopping the Services
To stop all running services, use:

```bash
docker-compose down
```
This will stop and remove all containers, networks, and volumes created by docker-compose up.

## Step 7: Scaling Services (Optional)
If you want to scale a service (e.g., increase the number of replicas), you can use:

```bash
docker-compose up --scale <service-name>=<number-of-instances>
```

Replace <service-name> with the service you want to scale and <number-of-instances> with the desired number of replicas.