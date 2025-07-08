# Copyright 2021 Synology Inc.
FROM alpine:latest
LABEL maintainers="Synology Authors" \
      description="Synology CSI Plugin"

RUN apk add --no-cache e2fsprogs e2fsprogs-extra xfsprogs xfsprogs-extra blkid util-linux iproute2 bash btrfs-progs ca-certificates cifs-utils nfs-utils

# Create symbolic link for chroot.sh
WORKDIR /

# Copy and run CSI driver
COPY bin/synology-csi-driver synology-csi-driver

ENTRYPOINT ["/synology-csi-driver"]
