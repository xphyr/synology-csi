# Copyright 2021 Synology Inc.

############## Build stage ##############
FROM golang:1.20.4-alpine as builder
LABEL stage=synobuilder

RUN apk add --no-cache alpine-sdk
COPY . /go/src/github.com/SynologyOpenSource/synology-csi/
WORKDIR /go/src/github.com/SynologyOpenSource/synology-csi/

ARG TARGETPLATFORM

RUN env GOARCH=$(echo "$TARGETPLATFORM" | cut -f2 -d/) \
        GOARM=$(echo "$TARGETPLATFORM" | cut -f3 -d/ | cut -c2-) \
        make

############## Final stage ##############
FROM alpine:latest
LABEL maintainers="Synology Authors" \
      description="Synology CSI Plugin"

RUN apk add --no-cache e2fsprogs e2fsprogs-extra xfsprogs xfsprogs-extra blkid util-linux iproute2 bash btrfs-progs ca-certificates cifs-utils

# Create symbolic link for chroot.sh
WORKDIR /
RUN mkdir /csibin
COPY chroot/chroot.sh /csibin
RUN chmod 777 /csibin/chroot.sh \
        && ln -s /csibin/chroot.sh /csibin/iscsiadm \
        && ln -s /csibin/chroot.sh /csibin/multipath \
        && ln -s /csibin/chroot.sh /csibin/multipathd

ENV PATH="/csibin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

# Copy and run CSI driver
COPY --from=builder /go/src/github.com/SynologyOpenSource/synology-csi/bin/synology-csi-driver synology-csi-driver

ENTRYPOINT ["/synology-csi-driver"]
