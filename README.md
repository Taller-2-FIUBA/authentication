# authentication

[![codecov](https://codecov.io/gh/Taller-2-FIUBA/authentication/branch/main/graph/badge.svg?token=MM8e5db9JF)](https://codecov.io/gh/Taller-2-FIUBA/authentication)

Service to handle JWT and Firebase services.
To be accessed from backend.

# Launching directly

```bash
go get .
go run . PORT
```

where `PORT` is the port where the service will be listening for requests.

# Tests

```bash
go test
```

# Launching through Docker

```bash
docker build --tag IMAGE_NAME .
docker run --rm -p 8002:80 --name CONTAINER_NAME IMAGE_NAME
```

Where `IMAGE_NAME` is the name chosen in the previous step and
`CONTAINER_NAME` is a name to identify the container running the app  
Notice `--rm` tells docker to remove the container after it stops, and `-p 8002:80` maps 
the port 80 in the container to the port 8002 in the host.  
The specific port can be changed in the Dockerfile.