# launchpad-api

Migrations
==========

manual migration
```
migrate -path migrations --database "postgres://localhost/launchpad?sslmode=disable&user=lpad&password=secret" up
```

new migration
```
migrate create -ext sql -dir migrations create_something
```

Generate OAPI docs
==================

```
swag i -g cmd/lpadapi/main.go -o docs
```
