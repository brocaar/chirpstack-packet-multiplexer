FROM balenalib/raspberrypi3-golang:1.13-sid-build as pktmux-builder

ENV PROJECT_PATH=/chirpstack-packet-multiplexer
ENV PATH=$PATH:$PROJECT_PATH/build
ENV CGO_ENABLED=0
ENV GO_EXTRA_BUILD_ARGS="-a -installsuffix cgo"

RUN apt-get update
RUN apt-get install tzdata make git bash

RUN mkdir -p $PROJECT_PATH
COPY . $PROJECT_PATH
WORKDIR $PROJECT_PATH

RUN make dev-requirements
RUN make

FROM balenalib/raspberrypi3-golang:1.13-sid-run as pktmux-runner

WORKDIR /root/
COPY --from=pktmux-builder /chirpstack-packet-multiplexer/build .
ENTRYPOINT ["./chirpstack-packet-multiplexer"]
