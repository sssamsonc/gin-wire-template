# gin-wire-template (RESTful apis)
## Prerequisite
1. MongoDB, MySQL db, Redis
2. install gin-swagger
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```
3. install google/wire
```bash
go install github.com/google/wire/cmd/wire@latest
```

## env setup
1. Create a replica ".env.example" file and name it to ".env"
2. Modify the var

## start the DEMO
1. set IS_DEMO_MODE to true in your env, if you want to ignore all database connections for demo purpose
2. execute the project
3. visit http://localhost:8080/swagger/index.html and you will see all the api demo

## WIRE
1. After finished the modification for any dependencies injection, you will need to re-generate the wire_gen.go
```bash
wire
```
## SWAGGER
1. Please add the comment to your controller by following the guide:
https://github.com/swaggo/swag/blob/master/README.md#declarative-comments-format
2. After finished the modification, please re-generate the (swagger docs) "docs" folder
```bash
swag i
```
3. start the project and visit http://localhost:8080/swagger/index.html to see the changes