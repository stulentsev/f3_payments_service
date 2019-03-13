# F3 Payments Service

Simple REST API with a standard CQRS/ES design.

The API spec can be found in `swagger/swagger.yaml`

## Usage

#### Generating the API from swagger spec

```sh
make swagger
```

#### Testing

Local:
```sh
make test
```

With docker:
```sh
make docker_test
```

##### Running

```sh
docker-compose up server
```
