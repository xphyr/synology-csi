package sanitytest

import (
	"fmt"
	"os"
	"testing"

	sanity "github.com/kubernetes-csi/csi-test/v4/pkg/sanity"

	"github.com/xphyr/synology-csi/pkg/driver"
	"github.com/xphyr/synology-csi/pkg/dsm/common"
	"github.com/xphyr/synology-csi/pkg/dsm/service"
	"github.com/xphyr/synology-csi/pkg/utils/hostexec"
)

const (
	ConfigPath      = "./../../config/client-info-template.yml"
	SecretsFilePath = "./sanity-test-secret-file.yaml"
)

func TestSanity(t *testing.T) {
	nodeID := "CSINode"

	endpointFile, err := os.CreateTemp("", "csi-gcs.*.sock")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(endpointFile.Name())

	stagingPath, err := os.MkdirTemp("", "csi-gcs-staging")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(stagingPath)

	targetPath, err := os.MkdirTemp("", "csi-gcs-target")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(targetPath)

	info, err := common.LoadConfig(ConfigPath)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	dsmService := service.NewDsmService()

	for _, client := range info.Clients {
		err := dsmService.AddDsm(client)
		if err != nil {
			fmt.Printf("Failed to add DSM: %s, error: %v\n", client.Host, err)
		}
	}

	if dsmService.GetDsmsCount() == 0 {
		t.Fatal("No available DSM.")
	}
	defer dsmService.RemoveAllDsms()

	cmdExecutor, err := hostexec.New(nil, "")
	if err != nil {
		t.Fatalf("Failed to create command executor: %v\n", err)
	}
	tools := driver.NewTools(cmdExecutor)

	endpoint := "unix://" + endpointFile.Name()
	drv, err := driver.NewControllerAndNodeDriver(nodeID, endpoint, dsmService, tools)
	if err != nil {
		t.Fatalf("Failed to create driver: %v\n", err)
	}
	drv.Activate()

	// Set configuration options as needed
	testConfig := sanity.NewTestConfig()
	testConfig.TargetPath = targetPath
	testConfig.StagingPath = stagingPath
	testConfig.Address = endpoint
	testConfig.SecretsFile = SecretsFilePath

	// Set Input parameters for test
	testConfig.TestVolumeParameters = map[string]string{
		"protocol": "smb",
	}

	// testConfig.TestVolumeAccessType = "block" // raw block

	// Run test
	sanity.Test(t, testConfig)
}
