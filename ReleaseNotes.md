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

# v1.4.0
This release introduces potentially **BREAKING CHANGES** If you are converting from the Upstream synology/synology-csi driver or a prior version of xphyr/synology-csi driver please read the following.
  - All go pkg imports are now pointed to xphyr/synology-csi
  - Support for Kubernetes prior to v1.25 has been dropped. Kubernetes 1.20 is over 5.5 years old, and its EoL was [February of 2022](https://kubernetes.io/releases/patch-releases/)
     - removed deployment files for kubernetes 1.19 and kubernetes 1.20
  - In order to start code cleanup, I have removed code that is duplicated by go built in packages such as `strconv.ParseBool()`. This means that use of the term "yes" in any configuration file is NO LONGER SUPPORTED. If your storageClass definition contains the word "yes" or "no" for settings, you need to update to "true" or "false". This also aligns with kubernetes defaults for boolean values
  - By default the Synology RecycleBin is DISABLED for NFS and SMB shares. 
    - **NOTE:** This is a "breaking change". Previous releases enabled the recyclebin by default. This causes issues with many apps in K8s, so going forward all NEW shares will have the recyclebin disabled by default. You can enable the recyclebin by setting `recycleBin: true` in the storageClass definition file.
    - **NOTE:** This is also "breaking change". Previous releases with the recyclebin enabled an option that only allowed "Administrator" access to the files in the recycleBin. This causes issues with many apps in K8s, so going forward all NEW shares that have the recyclebin enabled will allow anyone access to the recycleBin. You can change this recyclebin by setting `recycleBinAdminOnly: true` in the storageClass definition file.
  - removed the final deprecated function calls from k8s.io/mount-utils from code base

# v1.4.1
 - Release 1.4.0 had an issue with the cluster-role for the snapshotter. Changes made in the deployment file in v1.3.2 were lost in the change over to supporting v1.25 or higher.

# v1.4.2
 - Added support for using a dedicated subnet for storage access. In prior releases if you were connecting to a Synology array over a dedicated subnet, NFS would not work due to NFS ACLs. There is now a "clientsubnetoverride" as part of the `client-info.yaml` file which will allow you to override the NFS client permissions to a subnet as defined in [CIDR notation](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing). 

# v1.4.3
 - Updated description to include the kubernetes cluster name to help in resolving where a lun is in use if you have multiple clusters. This will update the description to have the following `<clusterName>/<namespace>/<pvcname>` This will help in identifying what cluster a LUN or share belongs to if you have multiple clusters. If you do not supply a clusterName the original description of `<namespace>/<pvcname>` will be used

# v1.4.4.
  - Updated the synocli to list out hosts. 
  - Added additional functionality to the webapi module for use in another project.
  - Updated the PersistentVolume VolumeAttributes field to include the LUN UUID