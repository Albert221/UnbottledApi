FROM golang:alpine AS build

ADD . /src
RUN cd /src && go build -o unbottled .

FROM alpine

WORKDIR /app
COPY --from=build /src/unbottled /app/
ENTRYPOINT ./unbottled migrate && ./unbottled serve