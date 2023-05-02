# authentication

Service to handle JWT and Firebase services.
To be accessed from backend.

# Launching directly

```
go get .
go run . PORT
```

where `PORT` is the port where the service will be listening for requests.

# Tests

```
go test 
```

# Lauching through Docker

```
docker build --tag IMAGE_NAME .
docker run --rm -p 8082:8082 --name CONTAINER_NAME IMAGE_NAME
```
Where `IMAGE_NAME` is the name chosen in the previous step and
`CONTAINER_NAME` is a name to identify the container running the app  
Notice `--rm` tells docker to remove the container after it stops, and `-p 8082:8082` maps 
the port 8082 in the container to the port 8082 in the host.  
The specific port can be changed in the Dockerfile.