FROM docker.io/matrixdotorg/synapse:v1.1.0

RUN apk add --no-cache gcc musl-dev openssl go curl ca-certificates dos2unix git jq

RUN wget -O go.tgz https://dl.google.com/go/go1.12.6.src.tar.gz
RUN tar -C /usr/local -xzf go.tgz
WORKDIR /usr/local/go/src/
RUN sh ./make.bash
ENV GOPATH="/opt/go"
ENV PATH="/usr/local/go/bin:$GOPATH/bin:$PATH"
RUN go version
RUN env

ENV SYNAPSE_CONFIG_DIR=/synapse
ENV SYNAPSE_CONFIG_PATH=/synapse/homeserver.yaml
ENV UID 991
ENV GID 991

RUN mkdir -p /synapse
COPY .docker/synapse /synapse
RUN chown -R 991:991 /synapse

RUN mkdir -p /project/src/github.com/turt2live/terraform-provider-matrix
ENV GOPATH="$GOPATH:/project:/project/src/github.com/turt2live/terraform-provider-matrix/vendor"
ENV GO111MODULE=on

COPY /.docker/run-tests.sh /run-tests.sh
RUN chmod +x /run-tests.sh && dos2unix /run-tests.sh

COPY . /project/src/github.com/turt2live/terraform-provider-matrix

ENTRYPOINT [ "/run-tests.sh" ]