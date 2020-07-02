
FROM golang:1.9.2

COPY . .

RUN make

FROM scratch

MAINTAINER Alex

COPY --from=0 ./_output/bin/descheduler-controller /bin/descheduler-controller

CMD ["/bin/descheduler-controller", "--help"]
