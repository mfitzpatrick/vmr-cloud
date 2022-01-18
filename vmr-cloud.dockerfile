FROM golang:1.16-alpine AS build

WORKDIR /app
#disabling cgo means compiled binary will be completely statically linked
ENV CGO_ENABLED=0
#setting OS to linux means we can work in a completely empty docker container (i.e. scratch)
ENV GOOS=linux

COPY /go.* ./
RUN go mod download

COPY /src/*.go ./

RUN go build -o /vmr-cloud

FROM scratch

COPY --from=build /vmr-cloud /vmr-cloud

CMD ["/vmr-cloud"]

