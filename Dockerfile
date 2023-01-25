FROM balenalib/raspberry-pi-alphine-golang:1.13-build as pktmux-builder

ENV PROJECT_PATH=/chirpstack-packet-multiplexer
ENV PATH=$PATH:$PROJECT_PATH/build
ENV CGO_ENABLED=0
ENV GO_EXTRA_BUILD_ARGS="-a -installsuffix cgo"

RUN apk add --no-cache tzdata make git bash

RUN mkdir -p $PROJECT_PATH
COPY . $PROJECT_PATH
WORKDIR $PROJECT_PATH

RUN make dev-requirements
RUN make

FROM balenalib/raspberry-pi-alphine-golang:1.13-run as pktmux-runner

WORKDIR /root/
RUN apk --no-cache add tzdata
COPY --from=pktmux-builder /chirpstack-packet-multiplexer/build .
ENTRYPOINT ["./chirpstack-packet-multiplexer"]
