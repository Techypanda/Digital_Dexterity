FROM golang:1.18-alpine as BUILD

WORKDIR /app/
COPY . /app/
RUN go build /app/cmd/api

FROM golang:1.18-alpine as RUN
WORKDIR /root/
COPY --from=0 /app/api ./
EXPOSE 8080
CMD ["./api"]