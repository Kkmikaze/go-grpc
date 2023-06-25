# Go gRPC Documentation

## Table Of Content
- [Prerequisite](#prerequisites)
  - [Structure](#structure)
- [How To](#how-to)
- [References](#references)
  - [ORM](#orm)
  - [Architecture](#architecture)


### Prerequisites
What things you need to setup the application:
- [Air](https://github.com/cosmtrek/air) (not working yet)
- [Docker](https://www.docker.com/)
- [Golang](https://go.dev/doc/install)
- [Go Migrate](https://github.com/golang-migrate/migrate)
- Makefile
- MySQL

### Structure
```
.
└── api
|   └── embed.go
├── cmd
|   ├── server
|   │   └── main.go
|   └── root.go
├── common
|   ├── orm.go
|   ├── transform.go
|   └── validator.go
├── config
|   └── config.go
├── db
|   └── migration
|       └── {domain_dir}
|           └── *.sql
├── domain
|   └── {api_version_dir}
|       └── {domain_dir}
|           ├── internal
|           │   ├── entity
|           │   │   └── *.go
|           │   ├── handler
|           │   │   └── *.go
|           │   ├── repository
|           │   │   └── *.go
|           │   └── validator
|           │       └── *.go
|           ├── usecase
|           │   └── *.go
|           └── {domain}.go
├── interceptors
|   ├── rpc_client_interceptors.go
|   │
|   └── rpc_server_interceptors.go
├── internal
|   └── gateway
|       └── *.go
├── middleware
|   └── *go
├── pkg
|   ├── orm
|   │   └── *.go
|   ├── rcpclient
|   │   └── *.go
|   └── rcpserver
|       └── *.go
├── proto
|   └── {domain_dir}
|       └── {api_version_dir}
|           └── *.proto
├── scripts
|   └── *.sh
├── stubs
|   └── {domain_dir}
|       └── {api_version_dir}
|           └── *.pb.go
├── thid_party
|   ├── proto
|   └── swagger-ui
|   └── embed.go
├── .air.toml
├── .env.example
├── .gitignore
├── .golangci.yml
├── Dockerfile
├── go.mod
├── go.sum
├── LICENSE
├── Makefile
└── README.md
```

## How To
### Running The App (local)
- First get the dependencies with this command:
```shell
make dependency
```

- Copy the `.env.example` to `.env` with run this command:
```shell
cp .env-example .env
```

- Generate the stubs with this command:
```shell
make proto-gen
```

- and for running the application can use this command:
```shell
make debug
```

### Running The App (local) with Hot Reload
- For running the application with hot reload can use this command:
```shell
make air
```

## References
### GIT Style
For commit message style or git style guide, use this doc
- [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)


### ORM
ORM using GORM, follow this doc
- [GORM](https://gorm.io/docs/)


### Architecture
Architecture reference on go-grpc-microservices, follow this doc repo
- [go-grpc-microservices](https://github.com/Auxulry/go-grpc-microservices)