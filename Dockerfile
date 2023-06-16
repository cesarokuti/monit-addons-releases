# STEP 1 build executable binary
FROM golang:1.20 as builder
LABEL stage=intermediate

COPY . /releases-monitoring/
WORKDIR /releases-monitoring/

ARG SSH_PRIVATE_KEY
RUN mkdir /root/.ssh/
RUN ssh-keyscan github.com > /root/.ssh/known_hosts
RUN echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"
RUN go mod tidy
RUN go mod vendor

RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -o releases-monitoring main.go

# STEP 2 build a small image
# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder  /releases-monitoring/releases-monitoring .
USER 65532:65532


ENTRYPOINT [ "/releases-monitoring" ]
