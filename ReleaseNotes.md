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
