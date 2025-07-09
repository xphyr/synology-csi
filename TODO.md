# To Do list

The following is a list of tasks to bring the Synology CSI driver up to date.

- [x] Update Non-Kubernetes dependencies in build
- [x] Update logic to handle PV requests that do not meet the minimum allocatable size and bump it up to 1Gb.
- [x] Merge existing PR Requests from upstream project
   - [x] [PR85 - Predefined tool paths](https://github.com/SynologyOpenSource/synology-csi/pull/85)
       - NOTE: Change was pulled in but set to work in old mode. Must override chroot if no chroot is required
   - [X] [PR75 - Support for devAttribs](https://github.com/SynologyOpenSource/synology-csi/pull/75)
   - [X] [PR48 - return extra lun info and allow direct_io_pattern](https://github.com/SynologyOpenSource/synology-csi/pull/48)
- Add basic CI Pipeline
     - [x] Support Multi-arch (x86_64 and ARM) containers
       - [x] Use [GoReleaser](https://goreleaser.com/) to handle multiarch builds
     - [x] Update docs and deployment files to use new ghcr.io repo
- [x] Address known security vulnerabilities in go modules
- [x] `gofmt` to apply proper formatting and simplification to the entire codebase
  - [ ] ensure gofmt is applied as part of the build process/ci 
- [x] apply `go-staticcheck` across all files and address issues.
- [x] Bump to modern version of go compiler (1.24)
- [ ] Add support for OpenShift SCC profiles
  - [x] Add SCC definition file to repo
  - [ ] Update deployment documentation
  - [ ] Update deployment scripts
  - [ ] Update Helm Chart
- [ ] Add configuration to disable RecycleBin on NFS/SMB shares
- [ ] Add configuration to allow non-admin access to RecycleBin
- [ ] Override NFS mount Permissions for multi-homed servers
- [ ] Address Issues in upstream project
- [ ] Update code base to use current CSI releases
- [ ] Update testing framework
- [ ] Remove deprecated function calls in k8s.io/mount-utils
- [ ] Support Windows
