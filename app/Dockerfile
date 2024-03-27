FROM golang:1.18-alpine
WORKDIR /build
COPY entrypoint.sh .
EXPOSE ${APP_PORT}
CMD [ "sh", "entrypoint.sh"]
