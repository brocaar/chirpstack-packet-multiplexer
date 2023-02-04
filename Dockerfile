FROM --platform=linux/arm64 balenalib/raspberrypi3-golang:1-bullseye-build as pktmux-builder

ENV PROJECT_PATH=/chirpstack-packet-multiplexer
ENV PATH=$PATH:$PROJECT_PATH/build
ENV CGO_ENABLED=0
ENV GO_EXTRA_BUILD_ARGS="-a -installsuffix cgo"
ENV OUTPUT_DIR=/opt/output/pktmux-dependencies
ENV INPUT_DIR=/opt/input

RUN apt-get update
RUN apt-get install tzdata make git bash python3 python3-pip

# Deal with Python deps
WORKDIR "$INPUT_DIR"
COPY requirements.txt requirements.txt
RUN pip3 install --target="$OUTPUT_DIR" --no-cache-dir -r requirements.txt

RUN mkdir -p $PROJECT_PATH
COPY . $PROJECT_PATH
WORKDIR $PROJECT_PATH

RUN make dev-requirements
RUN make

FROM --platform=linux/arm64 balenalib/raspberrypi3-golang:1-bullseye-run as pktmux-runner

# Python stuff
ENV ROOT_DIR=/opt
ENV PYTHON_DEPENDENCIES_DIR="$ROOT_DIR/pktmux-dependencies"
COPY --from=pktmux-builder /opt/output/pktmux-dependencies /opt/pktmux-dependencies
RUN apt-get update
RUN apt-get install python3
ENV PYTHONPATH=/opt/pktmux-dependencies

WORKDIR /root/
COPY --from=pktmux-builder /chirpstack-packet-multiplexer/build .
COPY --from=pktmux-builder /chirpstack-packet-multiplexer/start_multiplexer.py .
RUN mkdir -p /etc/chirpstack-packet-multiplexer
RUN chmod 755 start_multiplexer.py
COPY --from=pktmux-builder /chirpstack-packet-multiplexer/config/chirpstack-packet-multiplexer.toml /etc/chirpstack-packet-multiplexer/chirpstack-packet-multiplexer.toml
ENTRYPOINT ["/usr/bin/python3", "start_multiplexer.py"]
