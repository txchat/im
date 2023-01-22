## We'll choose the incredibly lightweight
## Go alpine image to work with
FROM golang:1.17.8 AS builder

## the lightweight scratch image we'll
## run our application within
FROM alpine:latest
## We have to copy the output from our
## builder stage to our production stage
ARG server_name
ENV GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn
ENV SERVERNAME=${server_name}
WORKDIR /usr/local/bin
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip
COPY ${server_name} .
COPY ${server_name}.* /etc/${server_name}/config/
## we can then kick off our newly compiled
## binary exectuable!!
CMD ["sh", "-c", "./${SERVERNAME} -f /etc/${SERVERNAME}/config/${SERVERNAME}.yaml"]