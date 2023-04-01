FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates
ENV TZ Asia/Shanghai
ENV GIN_MODE release

WORKDIR /app
COPY ./templates /app/templates
COPY ./app/main /app/main
EXPOSE 9999
CMD ["./main"]