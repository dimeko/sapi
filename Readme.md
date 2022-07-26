## Sapi
### Overview
A simple wallet repository server with the basic CRUD operations for listing, modifying wallets. Wallets are stored inside a json database that is created upon server starting with every json file name representing the wallet's id. Every wallet has an `id`, a `name` and a `balance`. The available operations are `create` (the id is generated with every creation without the need for the user to pass it as a pareameter), `add/{id}` funds to increase the wallet's balance, `remove/{id}` funds to decrease a wallets balace and finally `list` for the listing of all wallets and `get/{id}` to retreive one wallet.

### Setup
With `Docker` (recommended):
```
make docker-build
make docker-run
```

On host machine:
```
make host-run
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
    "name": "wallet_name",
    "balance": 500
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

##### Add
```
curl --location --request PUT 'http://localhost:6028/add/{id}' \
--header 'Content-Type: application/json' \
--data-raw '{
    "amount": 50
}'
```

##### Remove
```
curl --location --request PUT 'http://localhost:6028/remove/{id}' \
--header 'Content-Type: application/json' \
--data-raw '{
    "amount": 50
}'
```

### Cleanup
For a propper cleanup of the `Docker` container and image you can use the below command:
```
make docker-cleanup
```

### To have done better
1. Run tests inside a Docker container.
2. In order to increase performance we could write to file, on updates, using a goroutine.
3. Shutdown app with error code 0 and wait for every running process to finish (using wait groups).
4. Restrict access to db so it could be modified only from the app.