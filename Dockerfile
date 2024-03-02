FROM golang:alpine3.19 as build

WORKDIR /app

COPY api .

RUN go build -o /bin/rinha


FROM alpine:3.19
RUN apk add libc6-compat
COPY --from=build /bin/rinha /bin/rinha
COPY --from=build /app/.env .
EXPOSE 5000

CMD ["/bin/rinha"]