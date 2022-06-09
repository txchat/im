## the lightweight scratch image we'll
## run our application within
FROM alpine:latest
## We have to copy the output from our
## builder stage to our production stage
ENV GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn
WORKDIR /usr/local/bin
COPY logic .
COPY logic.* /etc/logic/config/
## we can then kick off our newly compiled
## binary exectuable!!
CMD ["./logic","-conf","/etc/logic/config/logic.toml"]
