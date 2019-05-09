FROM golang:1.12-alpine AS development

ENV PROJECT_PATH=/lora-packet-multiplexer
ENV PATH=$PATH:$PROJECT_PATH/build
ENV CGO_ENABLED=0
ENV GO_EXTRA_BUILD_ARGS="-a -installsuffix cgo"

RUN apk add --no-cache tzdata make git bash

RUN mkdir -p $PROJECT_PATH
COPY . $PROJECT_PATH
WORKDIR $PROJECT_PATH

RUN make dev-requirements
RUN make

FROM alpine:latest AS production

WORKDIR /root/
RUN apk --no-cache add tzdata
COPY --from=development /lora-packet-multiplexer .
ENTRYPOINT ["./lora-packet-multiplexer"]
