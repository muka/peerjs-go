FROM golang as base
ARG ARCH=amd64
ARG ARM=
ADD ./ /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} GOARM=${ARM} go build -o ./peer-server ./cmd/server/main.go

FROM scratch
COPY --from=base /build/peer-server /peer-server
ENTRYPOINT [ "/peer-server" ]