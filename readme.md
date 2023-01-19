## Gutenberghub Write API

This repository contains go api that extends the pocketbase as a framework. By extending we're adding the following enhancements in the pocketbase

1. **Rate Limiter**: We will be limiting users who creates more than 50 records in a minute.
1. **Secured Connection Route**: Adds a new custom route for securely changing a project connection.

## Build

Run the following command for building go executable:

```
go build main.go
```

## Serve 

Run the following command for serving pocketbase locally:

```
go run main.go serve
```