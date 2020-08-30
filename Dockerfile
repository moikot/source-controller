FROM --platform=$BUILDPLATFORM golang:1.14 as builder

# xx wraps go to automatically configure $GOOS, $GOARCH, and $GOARM
COPY --from=tonistiigi/xx:golang / /

WORKDIR /workspace

# copy api submodule
COPY api/ api/

# copy modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# cache modules
RUN go mod download

# copy source code
COPY main.go main.go
COPY controllers/ controllers/
COPY pkg/ pkg/
COPY internal/ internal/

# build
ARG TARGETPLATFORM
RUN CGO_ENABLED=0 GO111MODULE=on go build -a -o source-controller main.go

FROM alpine:3.12

RUN apk add --no-cache ca-certificates tini

COPY --from=builder /workspace/source-controller /usr/local/bin/

RUN addgroup -S controller && adduser -S -g controller controller

USER controller

ENTRYPOINT [ "/sbin/tini", "--", "source-controller" ]
