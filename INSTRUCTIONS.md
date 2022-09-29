This page covers setting up and executing my backend service.

## Setup
There are several applications that need to be installed for this project to run. 
- task
- postgres

### postgres
If you do not already have postgres installed, please go to https://www.postgresql.org/ to get a copy. The site 
had excellent instructions about setup.

### task
This project uses **task** as a build tool. Please go to https://taskfile.dev/installation/ for  instructions on
how to set up task on your environment. If you're using GoLand, you can also open the project there and use it to
get the service built and running.

## Building to service
To build the app, you will need to run the following commands:
- **task mod** : gets the vendor libraries this project depends on
- **task build** : creates the executable in the ./build/out subdirectory.

On my machine, here's what the process looks like:
```text
dgayle@localhost/cc_DavidGayle_BackendAPI % task mod
task: [mod] go mod tidy
task: [mod] go mod vendor
dgayle@localhost/cc_DavidGayle_BackendAPI % task build
task: [build] go build -o ./build/out/backend-svc cmd/main.go
dgayle@localhost/cc_DavidGayle_BackendAPI % 
```
After executing the task commands, you should have an executable named `backend-svc` in the `./build/out` directory.

## Configuring the service
In the `./config`, there is an example of the configuration file the service needs to run. The service uses the Viper
library to read in the config file and assumes the format is JSON. These are the contents of the file at the time this
page was written:
```json
{
  "db": {
    "type": "postgres",
    "host": "localhost",
    "port": 9001,
    "name": "tnp_config",
    "user": "postgres",
    "password": "mysecretpassword"
  }
}
```
The fields in the config file are as follows: 
- `db.type`: the type of database you are connecting to. In this case, *postgres*.
- `db.host`: The host Postgres is running on. In this case, *localhost*. If you are not running on localhost, you will
want to use the fully qualified host name or the ip address of the machine Postgres is running on.
- `db.port`: The port Postgres is listening on. In this case, *9001*. The default is *5432* and will probably be the 
one your instance uses.
- `db.name`: The name of the schema the users table will be created in. In this case, *tnp_config*. Please refer to the
Postgres documentation for instructions on setting up the schema. 
- `db.user`: The user account you use to connect to Postgres. In this case, *postgres*. Please refer to the Postgres
documentation for instructions on setting up a user.
- `db.password`: The password for the account referenced by `db.user`. 

To verify thew connection information in the config file, run `psql -h localhost -p 9001 -U postgres -d tnp_config` to
connect to the Postgres instance.

## Running the service
To run the service, enter `./build/out/backend-svc config/example-config.json` into a terminal in the project root
directory after building the service. Your screen should display:
```text
dgayle@localhost/cc_DavidGayle_BackendAPI % ./build/out/backend-svc config/example-config.json
Database connection successful.
Server started at http://localhost:8080

```

If you do not include a path to the config file, you will see the following:
```text
dgayle@localhost/cc_DavidGayle_BackendAPI % ./build/out/backend-svc
Path and config file name must be passed in as arguments
    example:    <executable> <file>
```

## Testing the interfaces
The interfaces were tested using Postman. For each interface, the URI, sample data, a sample response will be included.

### `POST /signup`
#### Request
```json
{
  "email": "test_20220908@axiomzen.co",
  "password": "axiomzen",
  "firstName": "Alex",
  "lastName": "Zimmerman"
}
```
#### Response
```json
{
  "jwt_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJlbWFpbCI6InRlc3RfMjAyMjA5MDhAYXhpb216ZW4uY28iLCJleHAiOjE2NjQ0ODM4OTh9.IYehGxKpWDYQku31RmuuK5BvP5bybtigfGDurPzpiuE"
}
```
#### Include JWT Token:
Not used by this interface.

### `POST /login`
#### Request
```json
{
  "email": "test_20220908@axiomzen.co",
  "password": "axiomzen"
}
```
#### Response
```json
{
  "jwt_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJlbWFpbCI6InRlc3RfMjAyMjA5MDhAYXhpb216ZW4uY28iLCJleHAiOjE2NjQ0ODQxNTR9.eoo4ly1WUjGloIeKjuJM2Wsi89U9Iikc2PaMDJo8q7k"
}
```
#### Include JWT Token:
Not used by this interface.

### `GET /users`
#### Request
None used by this request.

#### Response
```json
[
  {
    "email": "test_001@axiomzen.co",
    "firstName": "Alex",
    "lastName": "Zimmerman"
  },
  {
    "email": "test_002@axiomzen.co",
    "firstName": "Alex",
    "lastName": "Zimmerman"
  },
  {
    "email": "test_20220908@axiomzen.co",
    "firstName": "Alex",
    "lastName": "Zimmerman"
  }
]
```
#### Include JWT Token:
Used by the interface. In Postman, do the following:
- Go to the `Authorization` tab
- In the dropdown next to the `Type` label, select `API Key`
- In the `Key` field, enter `x-authentication-token`
- In the `Value` field, add the JWT token returned by calling either the `/signup` or `/login` interface.
- In the `Add To` field, select `Header`.

### `PUT /users`
#### Request
```json
{
  "firstName": "Sacha",
  "lastName": "Zimmerman"
}
```
#### Response
```json
{
    "response": "update succeeded"
}
```
#### Include JWT Token:
Used by the interface. See notes about using the token in the `GET /users` section.
