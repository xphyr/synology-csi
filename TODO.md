# ToDO list

The following is a list of tasks to bring the Synology CSI driver up to date.

- [x] Update Non-Kubernetes dependencies in build
- [x] Update logic to handle PV requests that do not meet the minimum allocatable size and bump it up to 1Gb.
- [ ] Merge existing PR Requests from upstream project
   - [x] [PR85 - Predefined tool paths](https://github.com/SynologyOpenSource/synology-csi/pull/85)
   - [ ] [PR79 - Add basic CI Pipeline](https://github.com/SynologyOpenSource/synology-csi/pull/79)
     - [ ] Support Multi-arch (x86_64 and ARM) containers
     - [ ] Update docs and deployment files to use new ghcr.io repo
   - [X] [PR75 - Support for devAttribs](https://github.com/SynologyOpenSource/synology-csi/pull/75)
   - [X] [PR48 - return extra lun info and allow direct_io_pattern](https://github.com/SynologyOpenSource/synology-csi/pull/48)
- [ ] Add support for OpenShift SCC profiles
- [ ] Address Issues in upstream project
- [ ] Update code base to use current CSI releases
- [ ] Update testing framework
- [ ] Add configuration to disable RecycleBin on NFS/SMB shares
- [ ] Add configuration to allow non-admin access to RecycleBin
- [x] `gofmt` to apply proper formatting and simplification to the entire codebase
  - [ ] ensure gofmt is applied as part of the build process/ci 
- [x] apply `go-staticcheck` across all files and address issues.