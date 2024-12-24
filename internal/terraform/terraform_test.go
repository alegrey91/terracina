// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package terracina

import (
	"context"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terracina/internal/addrs"
	"github.com/hashicorp/terracina/internal/configs"
	"github.com/hashicorp/terracina/internal/configs/configload"
	"github.com/hashicorp/terracina/internal/initwd"
	"github.com/hashicorp/terracina/internal/plans"
	"github.com/hashicorp/terracina/internal/providers"
	testing_provider "github.com/hashicorp/terracina/internal/providers/testing"
	"github.com/hashicorp/terracina/internal/provisioners"
	"github.com/hashicorp/terracina/internal/registry"
	"github.com/hashicorp/terracina/internal/states"

	_ "github.com/hashicorp/terracina/internal/logging"
)

// This is the directory where our test fixtures are.
const fixtureDir = "./testdata"

func TestMain(m *testing.M) {
	flag.Parse()

	// We have fmt.Stringer implementations on lots of objects that hide
	// details that we very often want to see in tests, so we just disable
	// spew's use of String methods globally on the assumption that spew
	// usage implies an intent to see the raw values and ignore any
	// abstractions.
	spew.Config.DisableMethods = true

	os.Exit(m.Run())
}

func testModule(t *testing.T, name string) *configs.Config {
	t.Helper()
	c, _ := testModuleWithSnapshot(t, name)
	return c
}

func testModuleWithSnapshot(t *testing.T, name string) (*configs.Config, *configload.Snapshot) {
	t.Helper()

	dir := filepath.Join(fixtureDir, name)
	// FIXME: We're not dealing with the cleanup function here because
	// this testModule function is used all over and so we don't want to
	// change its interface at this late stage.
	loader, _ := configload.NewLoaderForTests(t)

	// We need to be able to exercise experimental features in our integration tests.
	loader.AllowLanguageExperiments(true)

	// Test modules usually do not refer to remote sources, and for local
	// sources only this ultimately just records all of the module paths
	// in a JSON file so that we can load them below.
	inst := initwd.NewModuleInstaller(loader.ModulesDir(), loader, registry.NewClient(nil, nil))
	_, instDiags := inst.InstallModules(context.Background(), dir, "tests", true, false, initwd.ModuleInstallHooksImpl{})
	if instDiags.HasErrors() {
		t.Fatal(instDiags.Err())
	}

	// Since module installer has modified the module manifest on disk, we need
	// to refresh the cache of it in the loader.
	if err := loader.RefreshModules(); err != nil {
		t.Fatalf("failed to refresh modules after installation: %s", err)
	}

	config, snap, diags := loader.LoadConfigWithSnapshot(dir)
	if diags.HasErrors() {
		t.Fatal(diags.Error())
	}

	return config, snap
}

// testModuleInline takes a map of path -> config strings and yields a config
// structure with those files loaded from disk
func testModuleInline(t testing.TB, sources map[string]string) *configs.Config {
	t.Helper()

	cfgPath := t.TempDir()

	for path, configStr := range sources {
		dir := filepath.Dir(path)
		if dir != "." {
			err := os.MkdirAll(filepath.Join(cfgPath, dir), os.FileMode(0777))
			if err != nil {
				t.Fatalf("Error creating subdir: %s", err)
			}
		}
		// Write the configuration
		cfgF, err := os.Create(filepath.Join(cfgPath, path))
		if err != nil {
			t.Fatalf("Error creating temporary file for config: %s", err)
		}

		_, err = io.Copy(cfgF, strings.NewReader(configStr))
		cfgF.Close()
		if err != nil {
			t.Fatalf("Error creating temporary file for config: %s", err)
		}
	}

	loader, cleanup := configload.NewLoaderForTests(t)
	defer cleanup()

	// We need to be able to exercise experimental features in our integration tests.
	loader.AllowLanguageExperiments(true)

	// Test modules usually do not refer to remote sources, and for local
	// sources only this ultimately just records all of the module paths
	// in a JSON file so that we can load them below.
	inst := initwd.NewModuleInstaller(loader.ModulesDir(), loader, registry.NewClient(nil, nil))
	_, instDiags := inst.InstallModules(context.Background(), cfgPath, "tests", true, false, initwd.ModuleInstallHooksImpl{})
	if instDiags.HasErrors() {
		t.Fatal(instDiags.Err())
	}

	// Since module installer has modified the module manifest on disk, we need
	// to refresh the cache of it in the loader.
	if err := loader.RefreshModules(); err != nil {
		t.Fatalf("failed to refresh modules after installation: %s", err)
	}

	config, diags := loader.LoadConfigWithTests(cfgPath, "tests")
	if diags.HasErrors() {
		t.Fatal(diags.Error())
	}

	return config
}

// testSetResourceInstanceCurrent is a helper function for tests that sets a Current,
// Ready resource instance for the given module.
func testSetResourceInstanceCurrent(module *states.Module, resource, attrsJson, provider string) {
	module.SetResourceInstanceCurrent(
		mustResourceInstanceAddr(resource).Resource,
		&states.ResourceInstanceObjectSrc{
			Status:    states.ObjectReady,
			AttrsJSON: []byte(attrsJson),
		},
		mustProviderConfig(provider),
	)
}

// testSetResourceInstanceTainted is a helper function for tests that sets a Current,
// Tainted resource instance for the given module.
func testSetResourceInstanceTainted(module *states.Module, resource, attrsJson, provider string) {
	module.SetResourceInstanceCurrent(
		mustResourceInstanceAddr(resource).Resource,
		&states.ResourceInstanceObjectSrc{
			Status:    states.ObjectTainted,
			AttrsJSON: []byte(attrsJson),
		},
		mustProviderConfig(provider),
	)
}

func testProviderFuncFixed(rp providers.Interface) providers.Factory {
	if p, ok := rp.(*testing_provider.MockProvider); ok {
		// make sure none of the methods were "called" on this new instance
		p.GetProviderSchemaCalled = false
		p.ValidateProviderConfigCalled = false
		p.ValidateResourceConfigCalled = false
		p.ValidateDataResourceConfigCalled = false
		p.UpgradeResourceStateCalled = false
		p.ConfigureProviderCalled = false
		p.StopCalled = false
		p.ReadResourceCalled = false
		p.PlanResourceChangeCalled = false
		p.ApplyResourceChangeCalled = false
		p.ImportResourceStateCalled = false
		p.ReadDataSourceCalled = false
		p.CloseCalled = false
	}

	return func() (providers.Interface, error) {
		return rp, nil
	}
}

func testProvisionerFuncFixed(rp *MockProvisioner) provisioners.Factory {
	// make sure this provisioner has has not been closed
	rp.CloseCalled = false

	return func() (provisioners.Interface, error) {
		return rp, nil
	}
}

func mustResourceInstanceAddr(s string) addrs.AbsResourceInstance {
	addr, diags := addrs.ParseAbsResourceInstanceStr(s)
	if diags.HasErrors() {
		panic(diags.Err())
	}
	return addr
}

func mustConfigResourceAddr(s string) addrs.ConfigResource {
	addr, diags := addrs.ParseAbsResourceStr(s)
	if diags.HasErrors() {
		panic(diags.Err())
	}
	return addr.Config()
}

func mustAbsResourceAddr(s string) addrs.AbsResource {
	addr, diags := addrs.ParseAbsResourceStr(s)
	if diags.HasErrors() {
		panic(diags.Err())
	}
	return addr
}

func mustAbsOutputValue(s string) addrs.AbsOutputValue {
	p, diags := addrs.ParseAbsOutputValueStr(s)
	if diags.HasErrors() {
		panic(diags.Err())
	}
	return p
}

func mustProviderConfig(s string) addrs.AbsProviderConfig {
	p, diags := addrs.ParseAbsProviderConfigStr(s)
	if diags.HasErrors() {
		panic(diags.Err())
	}
	return p
}

func mustReference(s string) *addrs.Reference {
	p, diags := addrs.ParseRefStr(s)
	if diags.HasErrors() {
		panic(diags.Err())
	}
	return p
}

func mustModuleInstance(s string) addrs.ModuleInstance {
	p, diags := addrs.ParseModuleInstanceStr(s)
	if diags.HasErrors() {
		panic(diags.Err())
	}
	return p
}

// HookRecordApplyOrder is a test hook that records the order of applies
// by recording the PreApply event.
type HookRecordApplyOrder struct {
	NilHook

	Active bool

	IDs    []string
	States []cty.Value
	Diffs  []*plans.Change

	l sync.Mutex
}

func (h *HookRecordApplyOrder) PreApply(id HookResourceIdentity, dk addrs.DeposedKey, action plans.Action, priorState, plannedNewState cty.Value) (HookAction, error) {
	if plannedNewState.RawEquals(priorState) {
		return HookActionContinue, nil
	}

	if h.Active {
		h.l.Lock()
		defer h.l.Unlock()

		h.IDs = append(h.IDs, id.Addr.String())
		h.Diffs = append(h.Diffs, &plans.Change{
			Action: action,
			Before: priorState,
			After:  plannedNewState,
		})
		h.States = append(h.States, priorState)
	}

	return HookActionContinue, nil
}

// Below are all the constant strings that are the expected output for
// various tests.

const testTerracinaInputProviderOnlyStr = `
aws_instance.foo:
  ID = 
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = us-west-2
  type = 
`

const testTerracinaApplyStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = bar
  type = aws_instance
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  num = 2
  type = aws_instance
`

const testTerracinaApplyDataBasicStr = `
data.null_data_source.testing:
  ID = yo
  provider = provider["registry.terracina.io/hashicorp/null"]
`

const testTerracinaApplyRefCountStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = 3
  type = aws_instance

  Dependencies:
    aws_instance.foo
aws_instance.foo.0:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance
aws_instance.foo.1:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance
aws_instance.foo.2:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance
`

const testTerracinaApplyProviderAliasStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"].bar
  foo = bar
  type = aws_instance
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  num = 2
  type = aws_instance
`

const testTerracinaApplyProviderAliasConfigStr = `
another_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/another"].two
  type = another_instance
another_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/another"]
  type = another_instance
`

const testTerracinaApplyEmptyModuleStr = `
<no state>
Outputs:

end = XXXX
`

const testTerracinaApplyDependsCreateBeforeStr = `
aws_instance.lb:
  ID = baz
  provider = provider["registry.terracina.io/hashicorp/aws"]
  instance = foo
  type = aws_instance

  Dependencies:
    aws_instance.web
aws_instance.web:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  require_new = ami-new
  type = aws_instance
`

const testTerracinaApplyCreateBeforeStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  require_new = xyz
  type = aws_instance
`

const testTerracinaApplyCreateBeforeUpdateStr = `
aws_instance.bar:
  ID = bar
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = baz
  type = aws_instance
`

const testTerracinaApplyCancelStr = `
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance
  value = 2
`

const testTerracinaApplyComputeStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = computed_value
  type = aws_instance

  Dependencies:
    aws_instance.foo
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  compute = value
  compute_value = 1
  num = 2
  type = aws_instance
  value = computed_value
`

const testTerracinaApplyCountDecStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = bar
  type = aws_instance
aws_instance.foo.0:
  ID = bar
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = foo
  type = aws_instance
aws_instance.foo.1:
  ID = bar
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = foo
  type = aws_instance
`

const testTerracinaApplyCountDecToOneStr = `
aws_instance.foo:
  ID = bar
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = foo
  type = aws_instance
`

const testTerracinaApplyCountDecToOneCorruptedStr = `
aws_instance.foo:
  ID = bar
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = foo
  type = aws_instance
`

const testTerracinaApplyCountDecToOneCorruptedPlanStr = `
DIFF:

DESTROY: aws_instance.foo[0]
  id:   "baz" => ""
  type: "aws_instance" => ""



STATE:

aws_instance.foo:
  ID = bar
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = foo
  type = aws_instance
aws_instance.foo.0:
  ID = baz
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance
`

const testTerracinaApplyCountVariableStr = `
aws_instance.foo.0:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = foo
  type = aws_instance
aws_instance.foo.1:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = foo
  type = aws_instance
`

const testTerracinaApplyCountVariableRefStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = 2
  type = aws_instance

  Dependencies:
    aws_instance.foo
aws_instance.foo.0:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance
aws_instance.foo.1:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance
`
const testTerracinaApplyForEachVariableStr = `
aws_instance.foo["b15c6d616d6143248c575900dff57325eb1de498"]:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = foo
  type = aws_instance
aws_instance.foo["c3de47d34b0a9f13918dd705c141d579dd6555fd"]:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = foo
  type = aws_instance
aws_instance.foo["e30a7edcc42a846684f2a4eea5f3cd261d33c46d"]:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = foo
  type = aws_instance
aws_instance.one["a"]:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance
aws_instance.one["b"]:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance
aws_instance.two["a"]:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance

  Dependencies:
    aws_instance.one
aws_instance.two["b"]:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance

  Dependencies:
    aws_instance.one`
const testTerracinaApplyMinimalStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance
`

const testTerracinaApplyModuleStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = bar
  type = aws_instance
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  num = 2
  type = aws_instance

module.child:
  aws_instance.baz:
    ID = foo
    provider = provider["registry.terracina.io/hashicorp/aws"]
    foo = bar
    type = aws_instance
`

const testTerracinaApplyModuleBoolStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = true
  type = aws_instance
`

const testTerracinaApplyModuleDestroyOrderStr = `
<no state>
`

const testTerracinaApplyMultiProviderStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = bar
  type = aws_instance
do_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/do"]
  num = 2
  type = do_instance
`

const testTerracinaApplyModuleOnlyProviderStr = `
<no state>
module.child:
  aws_instance.foo:
    ID = foo
    provider = provider["registry.terracina.io/hashicorp/aws"]
    type = aws_instance
  test_instance.foo:
    ID = foo
    provider = provider["registry.terracina.io/hashicorp/test"]
    type = test_instance
`

const testTerracinaApplyModuleProviderAliasStr = `
<no state>
module.child:
  aws_instance.foo:
    ID = foo
    provider = module.child.provider["registry.terracina.io/hashicorp/aws"].eu
    type = aws_instance
`

const testTerracinaApplyModuleVarRefExistingStr = `
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = bar
  type = aws_instance

module.child:
  aws_instance.foo:
    ID = foo
    provider = provider["registry.terracina.io/hashicorp/aws"]
    type = aws_instance
    value = bar

    Dependencies:
      aws_instance.foo
`

const testTerracinaApplyOutputOrphanModuleStr = `
<no state>
`

const testTerracinaApplyProvisionerStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance

  Dependencies:
    aws_instance.foo
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  compute = value
  compute_value = 1
  num = 2
  type = aws_instance
  value = computed_value
`

const testTerracinaApplyProvisionerModuleStr = `
<no state>
module.child:
  aws_instance.bar:
    ID = foo
    provider = provider["registry.terracina.io/hashicorp/aws"]
    type = aws_instance
`

const testTerracinaApplyProvisionerFailStr = `
aws_instance.bar: (tainted)
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  num = 2
  type = aws_instance
`

const testTerracinaApplyProvisionerFailCreateStr = `
aws_instance.bar: (tainted)
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance
`

const testTerracinaApplyProvisionerFailCreateNoIdStr = `
<no state>
`

const testTerracinaApplyProvisionerFailCreateBeforeDestroyStr = `
aws_instance.bar: (tainted) (1 deposed)
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  require_new = xyz
  type = aws_instance
  Deposed ID 1 = bar
`

const testTerracinaApplyProvisionerResourceRefStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  num = 2
  type = aws_instance
`

const testTerracinaApplyProvisionerSelfRefStr = `
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = bar
  type = aws_instance
`

const testTerracinaApplyProvisionerMultiSelfRefStr = `
aws_instance.foo.0:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = number 0
  type = aws_instance
aws_instance.foo.1:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = number 1
  type = aws_instance
aws_instance.foo.2:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = number 2
  type = aws_instance
`

const testTerracinaApplyProvisionerMultiSelfRefSingleStr = `
aws_instance.foo.0:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = number 0
  type = aws_instance
aws_instance.foo.1:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = number 1
  type = aws_instance
aws_instance.foo.2:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = number 2
  type = aws_instance
`

const testTerracinaApplyProvisionerDiffStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = bar
  type = aws_instance
`

const testTerracinaApplyProvisionerSensitiveStr = `
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance
`

const testTerracinaApplyDestroyStr = `
<no state>
`

const testTerracinaApplyErrorStr = `
aws_instance.bar: (tainted)
  ID = 
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = 2

  Dependencies:
    aws_instance.foo
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance
  value = 2
`

const testTerracinaApplyErrorCreateBeforeDestroyStr = `
aws_instance.bar:
  ID = bar
  provider = provider["registry.terracina.io/hashicorp/aws"]
  require_new = abc
  type = aws_instance
`

const testTerracinaApplyErrorDestroyCreateBeforeDestroyStr = `
aws_instance.bar: (1 deposed)
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  require_new = xyz
  type = aws_instance
  Deposed ID 1 = bar
`

const testTerracinaApplyErrorPartialStr = `
aws_instance.bar:
  ID = bar
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance

  Dependencies:
    aws_instance.foo
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  type = aws_instance
  value = 2
`

const testTerracinaApplyResourceDependsOnModuleStr = `
aws_instance.a:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  ami = parent
  type = aws_instance

  Dependencies:
    module.child.aws_instance.child

module.child:
  aws_instance.child:
    ID = foo
    provider = provider["registry.terracina.io/hashicorp/aws"]
    ami = child
    type = aws_instance
`

const testTerracinaApplyResourceDependsOnModuleDeepStr = `
aws_instance.a:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  ami = parent
  type = aws_instance

  Dependencies:
    module.child.module.grandchild.aws_instance.c

module.child.grandchild:
  aws_instance.c:
    ID = foo
    provider = provider["registry.terracina.io/hashicorp/aws"]
    ami = grandchild
    type = aws_instance
`

const testTerracinaApplyResourceDependsOnModuleInModuleStr = `
<no state>
module.child:
  aws_instance.b:
    ID = foo
    provider = provider["registry.terracina.io/hashicorp/aws"]
    ami = child
    type = aws_instance

    Dependencies:
      module.child.module.grandchild.aws_instance.c
module.child.grandchild:
  aws_instance.c:
    ID = foo
    provider = provider["registry.terracina.io/hashicorp/aws"]
    ami = grandchild
    type = aws_instance
`

const testTerracinaApplyTaintStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  num = 2
  type = aws_instance
`

const testTerracinaApplyTaintDepStr = `
aws_instance.bar:
  ID = bar
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = foo
  num = 2
  type = aws_instance

  Dependencies:
    aws_instance.foo
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  num = 2
  type = aws_instance
`

const testTerracinaApplyTaintDepRequireNewStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = foo
  require_new = yes
  type = aws_instance

  Dependencies:
    aws_instance.foo
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  num = 2
  type = aws_instance
`

const testTerracinaApplyOutputStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = bar
  type = aws_instance
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  num = 2
  type = aws_instance

Outputs:

foo_num = 2
`

const testTerracinaApplyOutputAddStr = `
aws_instance.test.0:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = foo0
  type = aws_instance
aws_instance.test.1:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = foo1
  type = aws_instance

Outputs:

firstOutput = foo0
secondOutput = foo1
`

const testTerracinaApplyOutputListStr = `
aws_instance.bar.0:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = bar
  type = aws_instance
aws_instance.bar.1:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = bar
  type = aws_instance
aws_instance.bar.2:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = bar
  type = aws_instance
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  num = 2
  type = aws_instance

Outputs:

foo_num = [bar,bar,bar]
`

const testTerracinaApplyOutputMultiStr = `
aws_instance.bar.0:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = bar
  type = aws_instance
aws_instance.bar.1:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = bar
  type = aws_instance
aws_instance.bar.2:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = bar
  type = aws_instance
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  num = 2
  type = aws_instance

Outputs:

foo_num = bar,bar,bar
`

const testTerracinaApplyOutputMultiIndexStr = `
aws_instance.bar.0:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = bar
  type = aws_instance
aws_instance.bar.1:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = bar
  type = aws_instance
aws_instance.bar.2:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  foo = bar
  type = aws_instance
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  num = 2
  type = aws_instance

Outputs:

foo_num = bar
`

const testTerracinaApplyUnknownAttrStr = `
aws_instance.foo: (tainted)
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  num = 2
  type = aws_instance
`

const testTerracinaApplyVarsStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  bar = override
  baz = override
  foo = us-east-1
aws_instance.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  bar = baz
  list.# = 2
  list.0 = Hello
  list.1 = World
  map.Baz = Foo
  map.Foo = Bar
  map.Hello = World
  num = 2
`

const testTerracinaApplyVarsEnvStr = `
aws_instance.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/aws"]
  list.# = 2
  list.0 = Hello
  list.1 = World
  map.Baz = Foo
  map.Foo = Bar
  map.Hello = World
  string = baz
  type = aws_instance
`

const testTerracinaRefreshDataRefDataStr = `
data.null_data_source.bar:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/null"]
  bar = yes
data.null_data_source.foo:
  ID = foo
  provider = provider["registry.terracina.io/hashicorp/null"]
  foo = yes
`
