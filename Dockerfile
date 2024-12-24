FROM golang as build-stage

WORKDIR /dist

COPY ./ /dist/

RUN go build -o server ./main.go

FROM ubuntu as runtime

COPY --from=build-stage /dist/ /usr/share

WORKDIR /usr/share

ENV GIN_MODE=release

EXPOSE 80 

CMD [ "./server","run","server","--port","0.0.0.0:80" ]