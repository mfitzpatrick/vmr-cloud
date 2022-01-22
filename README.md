![test status](https://github.com/mfitzpatrick/vmr-cloud/actions/workflows/test.yml/badge.svg)

# Cloud-based API service for VMR Taskings
This project implements an API service which can be used for VMR vessel taskings. It exposes endpoints like
`/voyage`, `/risk`, and `/assist` which can be used to add, update, or retrieve information about current
taskings and voyages. It is distributed via a docker image.

## Getting Started
### With Docker
```
docker run --rm -p "80:80" ghcr.io/mfitzpatrick/vmr-cloud:<version>
```

### Natively
Clone the project, then:
```
go mod download
cd src/
go run .
```

## Tests and Examples
Integration tests and API examples are provided via [Postman](https://www.postman.com/) collections.
To make use of it, import the `VMR Cloud API` collection file into Postman to view each of the available
API endpoints. From there, existing endpoints can be modified to interact with the service manually.

The service location (hostname) is referenced via collection variables. The `{{server}}` collection variable
contains the hostname and port (default `localhost:80`), and the `{{schema}}` collection variable contains
the URL schema (by default this is `http`).

An example of the full API usage is in `VMR Cloud Working Example`. This suite is run as part of the
integration tests in CI and can be used as an example of how to interact with the API endpoints. The endpoints
in this collection have Pre-request Scripts and Tests scripts configured which manage the data received.
Ultimately, the purpose of the scripts is to allow this test suite to validate the data received from the
service.

NB: Currently there is a limitation in the list endpoint validation. When running this collection as a whole,
ensure that the service has been freshly started otherwise the final voyage-list test endpoint will fail.
This may be fixed in a future test update.

