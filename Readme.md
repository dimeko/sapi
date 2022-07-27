## Sapi
### Overview
A simple server with the basic CRUD operations for listing, modifying users. Initially developed to be used as a boilerplate Go api. Future plans: add dynamic model creation, add actions on models

### Setup
With `Docker` (recommended):
```
make run
```

### Tests
```
make test
```

### Manual testing
Use the below example commands for manual testing (`curl` required):

##### Create
```
curl --location --request POST 'http://localhost:6028/create' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "username1",
    "firstname": "firstname1",
    "lastname": "lastname1"
}'
```

##### List
```
curl --location --request GET 'http://localhost:6028/list'
```

##### Get
```
curl --location --request GET 'http://localhost:6028/get/{id}'
```

##### Update
```
curl --location --request PUT 'http://localhost:6028/update/{id}' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "updated_username1",
    "firstname": "updated_firstname1",
    "lastname": "updated_lastname1"
}'
```

##### Delete
```
curl --location --request PUT 'http://localhost:6028/delete/{id}' \
```

### Cleanup
For a propper cleanup of the `Docker` container and image you can use the below command:
```
make cleanup
```

### TODO
1. Shutdown app with error code 0 and wait for every running process to finish (using wait groups).
2. Implement live reload with `Docker`.
3. Improve http error responses.
4. Improve migration commands