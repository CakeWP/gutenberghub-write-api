## Gutenberghub Write API

This repository contains go api that extends the pocketbase as a framework. By extending we're adding the following enhancements in the pocketbase

1. **Rate Limiter**: We will be limiting users who creates more than 50 records in a minute.
2. **Fields Query Parameter**: Adds a new ``fields`` query parameter support, That basically allows explicitly selecting fields in each response. This will make the response much light-weight, and purposeful.

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