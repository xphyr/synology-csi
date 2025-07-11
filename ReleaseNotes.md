# [v1.3.0](https://github.com/xphyr/synology-csi/releases/tag/v1.3.0)
Initial test release of version 1.3.0 which includes the following changes:
 - Update Non-Kubernetes dependencies in build
 - Update logic to handle PV requests that do not meet the minimum allocatable size and bump it up to 1Gb.
 - Create basic CI Pipeline
   - Support Multi-arch (x86_64 and ARM) containers
   - Use [GoReleaser](https://goreleaser.com/) to handle multiarch builds
   - Update docs and deployment files to use new ghcr.io repo
 - gofmt to apply proper formatting and simplification to the entire codebase
 - apply go-staticcheck across all files and address issues.
 - Bump to modern version of go compiler (1.24)
 - Merged [PR85 - Predefined tool paths](https://github.com/SynologyOpenSource/synology-csi/pull/85) from upstream
 - Merged [PR75 - Support for devAttribs](https://github.com/SynologyOpenSource/synology-csi/pull/75) from upstream
 - Merged [PR48 - return extra lun info and allow direct_io_pattern](https://github.com/SynologyOpenSource/synology-csi/pull/48) from upstream
 - Partially removed deprecated function calls in k8s.io/mount-utils from code

# [v1.3.1](https://github.com/xphyr/synology-csi/releases/tag/v1.3.1)
- updated golang modules to remove known critical and high vulnerabilities
- updated install script to support installing on OpenShift clusters
- updated install script to support installing on Talos clusters

# [v1.3.2](https://github.com/xphyr/synology-csi/releases/tag/v1.3.2)
- addressed issues with snapshotter not properly creating snapshots
  - **NOTE**: You will need to apply the new `deploy\kubernetes\v1.20\snapshotter\snapshotter.yaml` file to your cluster as this updates the external-csi snapshotter application as well as the ClusterRole required to make snapshotting work on clusters based on K8s V1.20 and higher

# [v1.4.0]
This release introduces potentially **BREAKING CHANGES** If you are converting from the Upstream synology/synology-csi driver or a prior version of xphyr/synology-csi driver please read the following.
  - In order to start code cleanup, I have remvoed code that is duplicated by go built in packages such as `strconv.ParseBool()`. This means that use of the term "yes" in any configuration file is NO LONGER SUPPORTED. If your storageClass definition contains the word "yes" or "no" for settings, you need to update to "true" or "false". This also aligns with kubernetes defaults for boolean values
  - By default the Synology RecycleBin is DISABLED for NFS and SMB shares. If you want to use the #RecycleBin in your NFS or SMB shares, you will need to add 