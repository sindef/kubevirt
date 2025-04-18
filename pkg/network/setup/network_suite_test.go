package network

import (
	"os"
	"sync"
	"testing"

	"kubevirt.io/kubevirt/pkg/libvmi"
	"kubevirt.io/kubevirt/pkg/network/cache"
	kfs "kubevirt.io/kubevirt/pkg/os/fs"

	v1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/testutils"

	"kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/api"
)

func TestNetwork(t *testing.T) {
	testutils.KubeVirtTestSuiteSetup(t)
}

func newVMIBridgeInterface(namespace string, name string) *v1.VirtualMachineInstance {
	vmi := libvmi.New(
		libvmi.WithNamespace(namespace),
		libvmi.WithName(name),
		libvmi.WithNetwork(v1.DefaultPodNetwork()),
		libvmi.WithInterface(*v1.DefaultBridgeNetworkInterface()),
	)
	v1.SetObjectDefaults_VirtualMachineInstance(vmi)
	return vmi
}

func NewDomainWithBridgeInterface() *api.Domain {
	domain := &api.Domain{}
	domain.Spec.Devices.Interfaces = []api.Interface{{
		Model: &api.Model{
			Type: v1.VirtIO,
		},
		Type: "bridge",
		Source: api.InterfaceSource{
			Bridge: api.DefaultBridgeName,
		},
		Alias: api.NewUserDefinedAlias("default"),
	},
	}
	return domain
}

type tempCacheCreator struct {
	once   sync.Once
	tmpDir string
}

func (c *tempCacheCreator) New(filePath string) *cache.Cache {
	c.once.Do(func() {
		tmpDir, err := os.MkdirTemp("", "temp-cache")
		if err != nil {
			panic("Unable to create temp cache directory")
		}
		c.tmpDir = tmpDir
	})
	return cache.NewCustomCache(filePath, kfs.NewWithRootPath(c.tmpDir))
}
