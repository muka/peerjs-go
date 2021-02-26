FROM golang as base
ARG ARCH=amd64
ADD ./ /build
WORKDIR /build
RUN ARCH=${ARCH} make build

FROM scratch
ARG ARCH=amd64
COPY --from=base /build/build/${ARCH} /${ARCH}
ENTRYPOINT [ "${ARCH}" ]