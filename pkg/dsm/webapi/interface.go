/*
 * Copyright 2021 Synology Inc.
 */

package webapi

// DSMClient defines the interface for interacting with Synology DSM.
type DSMClient interface {
	Login() error
	Logout() error

	// Volume operations
	VolumeList() ([]VolInfo, error)
	VolumeGet(path string) (VolInfo, error)

	// LUN operations
	LunCreate(spec LunCreateSpec) (string, error)
	LunGet(name string) (LunInfo, error)
	LunMapTarget(targetIds []string, lunUuid string) error
	LunClone(spec LunCloneSpec) (string, error)
	LunDelete(uuid string) error
	LunUpdate(spec LunUpdateSpec) error

	// Target operations
	TargetCreate(spec TargetCreateSpec) (int, error)
	TargetGet(name string) (TargetInfo, error)
	TargetSet(targetId string, maxSession int) error
	TargetDelete(targetId string) error
	TargetList() ([]TargetInfo, error)

	// Snapshot operations
	SnapshotClone(spec SnapshotCloneSpec) (string, error)
	SnapshotCreate(spec SnapshotCreateSpec) (string, error)

	// Share operations
	ShareDelete(name string) error
	SetShareQuota(shareInfo ShareInfo, newSizeInMB int64) error
	ShareSnapshotCreate(spec ShareSnapshotCreateSpec) (string, error)

	// Metadata/System operations
	DsmInfoGet() (DsmInfo, error)

	// NFS operations
	NfsGet() (NfsInfo, error)
	NfsSet(enabled bool, majorVer bool, minorVer int) error
}
