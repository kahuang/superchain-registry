package superchain

import (
	"path"
	"strings"
	"testing"
)

func TestConfigs(t *testing.T) {
	n := 0
	for name, sch := range Superchains {
		if name != sch.Superchain {
			t.Errorf("superchain %q has bad key", name)
		}
		n += len(sch.ChainIDs)
	}
	for id, ch := range OPChains {
		if id != ch.ChainID {
			t.Errorf("chain %d has bad id", id)
		}
	}
	if len(OPChains) != n {
		t.Errorf("number of chains %d does not match chains in superchains %d", len(OPChains), n)
	}
	if len(OPChains) < 5 {
		t.Errorf("only got %d op chains, has everything loaded?", len(OPChains))
	}
	if len(Superchains) < 3 {
		t.Errorf("only got %d superchains, has everything loaded?", len(Superchains))
	}
	// All chains require extra addresses data until the L1 SystemConfig can support address mappings.
	if len(OPChains) != len(Addresses) {
		t.Errorf("got %d chains and %d address lists", len(OPChains), len(Addresses))
	}
	// All chains require extra genesis system config data until the
	// initial SystemConfig values can be read from the latest L1 chain state.
	if len(OPChains) != len(GenesisSystemConfigs) {
		t.Errorf("got %d chains and %d genesis system configs", len(OPChains), len(GenesisSystemConfigs))
	}
}

func TestGenesis(t *testing.T) {
	for id := range OPChains {
		_, err := LoadGenesis(id)
		if err != nil {
			t.Fatalf("failed to load genesis of chain %d: %v", id, err)
		}
	}
}

// TestImplementations ensures that the global Implementations
// map is populated.
func TestImplementations(t *testing.T) {
	if len(Implementations) == 0 {
		t.Fatal("no implementations found")
	}
}

// TestContractImplementations tests specific contracts implementations are set
// correctly.
func TestContractImplementations(t *testing.T) {
	impls, err := newContractImplementations("")
	if err != nil {
		t.Fatalf("failed to load contract implementations: %v", err)
	}
	if impls.L1CrossDomainMessenger.Get("1.6.0") != HexToAddress("0xf4d5682dA3ad1820ea83E1cEE5Fd92a3A7BabC30") {
		t.Fatal("wrong L1CrossDomainMessenger address")
	}
	if impls.L1ERC721Bridge.Get("1.3.0") != HexToAddress("0x8ADd7FB53A242e827373519d260EE3B8F7612Ba1") {
		t.Fatal("wrong L1ERC721Bridge address")
	}
	if impls.L1StandardBridge.Get("1.3.0") != HexToAddress("0x9c540e769B9453d174EdB683a90D9170e6559F16") {
		t.Fatal("wrong L1StandardBridge address")
	}
	if impls.L2OutputOracle.Get("1.5.0") != HexToAddress("0x7a811C9862ab54E677EEdA7e6F075aC86a1f551e") {
		t.Fatal("wrong L2OutputOracle address")
	}
	if impls.OptimismMintableERC20Factory.Get("1.4.0") != HexToAddress("0x135B9097A0e1e56190251c62f111B676Fb4Ec494") {
		t.Fatal("wrong OptimismMintableERC20 address")
	}
	if impls.OptimismPortal.Get("1.9.0") != HexToAddress("0x8Cfa294bD0c6F63cD65d492bdB754eAcf684D871") {
		t.Fatal("wrong OptimismPortal address")
	}
	if impls.SystemConfig.Get("1.7.0") != HexToAddress("0x09323D05868393c7EBa8190BAc173f843b82030a") {
		t.Fatal("wrong SystemConfig address")
	}
}

// TestContractVersionsCheck will fail if the superchain semver file
// is not read correctly.
func TestContractVersionsCheck(t *testing.T) {
	if err := SuperchainSemver.Check(); err != nil {
		t.Fatal(err)
	}
}

// TestContractVersionsResolve will test that the high lever interface used works.
func TestContractVersionsResolve(t *testing.T) {
	impls, err := newContractImplementations("goerli")
	if err != nil {
		t.Fatalf("failed to load contract implementations: %v", err)
	}

	if impls.L1CrossDomainMessenger.Get("1.6.0") == (Address{}) {
		t.Fatal("wrong L1CrossDomainMessenger address")
	}
	if impls.L1ERC721Bridge.Get("1.3.0") == (Address{}) {
		t.Fatal("wrong L1ERC721Bridge address")
	}
	if impls.L1StandardBridge.Get("1.3.0") == (Address{}) {
		t.Fatal("wrong L1StandardBridge address")
	}
	if impls.L2OutputOracle.Get("1.5.0") == (Address{}) {
		t.Fatal("wrong L2OutputOracle address")
	}
	if impls.OptimismMintableERC20Factory.Get("1.4.0") == (Address{}) {
		t.Fatal("wrong OptimismMintableERC20 address")
	}
	if impls.OptimismPortal.Get("1.9.0") == (Address{}) {
		t.Fatal("wrong OptimismPortal address")
	}
	if impls.SystemConfig.Get("1.7.0") == (Address{}) {
		t.Fatal("wrong SystemConfig address")
	}

	versions := ContractVersions{
		L1CrossDomainMessenger:       "1.6.0",
		L1ERC721Bridge:               "1.3.0",
		L1StandardBridge:             "1.3.0",
		L2OutputOracle:               "1.5.0",
		OptimismMintableERC20Factory: "1.4.0",
		OptimismPortal:               "1.9.0",
		SystemConfig:                 "1.7.0",
	}

	list, err := impls.Resolve(versions)
	if err != nil {
		t.Fatalf("unable to resolve: %s", err)
	}

	if list.L1CrossDomainMessenger.Version != "v1.6.0" {
		t.Fatalf("wrong L1CrossDomainMessenger version: %s", list.L1CrossDomainMessenger.Version)
	}
	if list.L1ERC721Bridge.Version != "v1.3.0" {
		t.Fatalf("wrong L1ERC721Bridge version: %s", list.L1ERC721Bridge.Version)
	}
	if list.L1StandardBridge.Version != "v1.3.0" {
		t.Fatalf("wrong L1StandardBridge version: %s", list.L1StandardBridge.Version)
	}
	if list.L2OutputOracle.Version != "v1.5.0" {
		t.Fatalf("wrong L2OutputOracle version: %s", list.L2OutputOracle.Version)
	}
	if list.OptimismMintableERC20Factory.Version != "v1.4.0" {
		t.Fatalf("wrong OptimismMintableERC20Factory version: %s", list.OptimismMintableERC20Factory.Version)
	}
	if list.OptimismPortal.Version != "v1.9.0" {
		t.Fatalf("wrong OptimismPortal version: %s", list.OptimismPortal.Version)
	}
	if list.SystemConfig.Version != "v1.7.0" {
		t.Fatalf("wrong SystemConfig version: %s", list.SystemConfig.Version)
	}
}

// TestResolve ensures that the low level resolve function works on semantic
// versioning correctly. It will return the highest version that matches the
// given semver string.
func TestResolve(t *testing.T) {
	cases := []struct {
		name    string
		set     AddressSet
		version string
		expect  string
	}{
		{
			name: "exact",
			set: AddressSet{
				"v1.0.0": HexToAddress("0x123"),
			},
			version: "v1.0.0",
			expect:  "v1.0.0",
		},
		{
			name: "largest-minor",
			set: AddressSet{
				"v1.2.0": HexToAddress("0x123"),
				"v1.1.0": HexToAddress("0x234"),
			},
			version: "^1.0.0",
			expect:  "v1.2.0",
		},
		{
			name: "largest-patch",
			set: AddressSet{
				"v1.0.2": HexToAddress("0x123"),
				"v1.0.1": HexToAddress("0x234"),
			},
			version: "^1.0.0",
			expect:  "v1.0.2",
		},
		{
			name: "x-patch",
			set: AddressSet{
				"v3.0.5": HexToAddress("0x123"),
				"v3.0.2": HexToAddress("0x234"),
			},
			version: "v3.0.x",
			expect:  "v3.0.5",
		},
		{
			name: "x-minor",
			set: AddressSet{
				"v2.5.1": HexToAddress("0x456"),
				"v2.5.0": HexToAddress("0x123"),
				"v2.2.2": HexToAddress("0x234"),
			},
			version: "v2.x",
			expect:  "v2.5.1",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			resolved, err := resolve(test.set, test.version)
			if err != nil {
				t.Fatal(err)
			}
			if resolved.Version != test.expect {
				t.Fatalf("wrong version: %s", resolved.Version)
			}
		})
	}
}

// TestAddressSet ensures that the AddressSet.Get method works with
// both the "v" prefix and without the "v" prefix.
func TestAddressSet(t *testing.T) {
	set := AddressSet{
		"v1.0.0": HexToAddress("0x123"),
		"1.1.0":  HexToAddress("0x234"),
	}

	if set.Get("v1.0.0") != HexToAddress("0x123") {
		t.Fatal("wrong address")
	}
	if set.Get("1.0.0") != HexToAddress("0x123") {
		t.Fatal("wrong address")
	}

	if set.Get("v1.1.0") != HexToAddress("0x234") {
		t.Fatal("wrong address")
	}
	if set.Get("1.1.0") != HexToAddress("0x234") {
		t.Fatal("wrong address")
	}
}

// TestContractBytecodes verifies that all bytecodes can be loaded successfully,
// and hash to the code-hash in the name.
func TestContractBytecodes(t *testing.T) {
	entries, err := extraFS.ReadDir(path.Join("extra", "bytecodes"))
	if err != nil {
		t.Fatalf("failed to open bytecodes dir: %v", err)
	}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".bin.gz") {
			t.Fatalf("bytecode file has missing suffix: %q", name)
		}
		name = strings.TrimSuffix(name, ".bin.gz")
		var expected Hash
		if err := expected.UnmarshalText([]byte(name)); err != nil {
			t.Fatalf("bytecode filename %q failed to parse as hash: %v", e.Name(), err)
		}
		value, err := LoadContractBytecode(expected)
		if err != nil {
			t.Fatalf("failed to load contract code of %q: %v", e.Name(), err)
		}
		computed := keccak256(value)
		if expected != computed {
			t.Fatalf("expected bytecode hash %s but computed %s", expected, computed)
		}
	}
}

// TestCanyonTimestampOnBlockBoundary asserts that Canyon will activate on a block's timestamp.
// This is critical because the create2Deployer only activates on a block's timestamp.
func TestCanyonTimestampOnBlockBoundary(t *testing.T) {
	for superchainName, superchainConfig := range Superchains {
		if superchainConfig.Config.CanyonTime == nil {
			continue
		}
		ct := *superchainConfig.Config.CanyonTime
		for _, id := range superchainConfig.ChainIDs {
			chainCfg := OPChains[id]
			canyonOffset := ct - chainCfg.Genesis.L2Time
			// Block time is hardcoded in op-node/rollup/superchain.go
			if canyonOffset%2 != 0 {
				t.Fatalf("Canyon time on superchain %v for %v is not on the block time. canyon time: %v. L2 start time: %v, block time: %v",
					superchainName, id, ct, chainCfg.Genesis.L2Time, 2)
			}
		}

		// Aevo specific test, 10 second block time
		if superchainName == "mainnet" {
			const aevoGenesisL2Time = 1679193011
			const aevoBlockTime = 10

			canyonOffset := ct - aevoGenesisL2Time
			if canyonOffset%aevoBlockTime != 0 {
				t.Fatalf("Canyon time on superchain %v for %v is not on the block time. canyon time: %v. L2 start time: %v, block time: %v",
					superchainName, "aevo", ct, aevoGenesisL2Time, aevoBlockTime)
			}
		}
	}
}
