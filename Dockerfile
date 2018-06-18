FROM docker.io/matrixdotorg/synapse:v0.31.2

RUN apk add --no-cache go curl ca-certificates dos2unix git jq

ENV SYNAPSE_SERVER_NAME="localhost"
ENV SYNAPSE_REPORT_STATS="no"
ENV SYNAPSE_ENABLE_REGISTRATION="true"
ENV SYNAPSE_ALLOW_GUEST="true"
ENV SYNAPSE_REGISTRATION_SHARED_SECRET="shared-secret-test1234"
ENV SYNAPSE_MACAROON_SECRET_KEY="macaroon-secret-test1234"

RUN mkdir -p /project/src/github.com/turt2live/terraform-provider-matrix
RUN mkdir -p /project/bin
ENV GOPATH="/project"
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && chmod +x /usr/local/bin/dep

COPY /.docker/run-tests.sh /run-tests.sh
RUN chmod +x /run-tests.sh && dos2unix /run-tests.sh

COPY . /project/src/github.com/turt2live/terraform-provider-matrix

ENTRYPOINT [ "/run-tests.sh" ]