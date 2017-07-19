FROM alpine:3.6

RUN apk update && apk add --no-cache libc6-compat

COPY cgnet /cgnet

EXPOSE 9101
ENTRYPOINT ["/cgnet", "serve"]
CMD ["--port", "9101"]
