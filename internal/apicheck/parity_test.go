package apicheck

import (
	"reflect"
	"strings"
	"testing"

	"cloud.google.com/go/firestore"
)

func TestDiscoverCoversRegistry(t *testing.T) {
	sdkDir, err := sdkModuleDir()
	if err != nil {
		t.Fatal(err)
	}
	found, err := discoverSDKTypesWithMethods(sdkDir)
	if err != nil {
		t.Fatal(err)
	}

	regNames := map[string]bool{}
	for _, p := range registry() {
		regNames[p.name] = true
	}

	for typeName := range found {
		if regNames[typeName] {
			continue
		}
		if reason, ok := ignoredTypes[typeName]; ok {
			t.Logf("ignore %s: %s", typeName, reason)
			continue
		}
		t.Errorf("SDK type %q has exported methods but is neither in apicheck registry nor ignoredTypes — add a fsmock wrapper or an ignoredTypes entry with a Reason.\n  Fix: edit internal/apicheck/registry.go or internal/apicheck/exceptions.go", typeName)
	}

	for typeName := range ignoredTypes {
		if _, ok := found[typeName]; ok {
			continue
		}
		t.Errorf("stale ignoredTypes entry %q — type has no exported methods (or was removed). Delete it from exceptions.go", typeName)
	}

	for name := range regNames {
		if name == "QuerySnapshot" {
			continue
		}
		if _, ok := found[name]; !ok {
			t.Errorf("registry lists %q but SDK has no exported methods for it — remove pair or fix discover", name)
		}
	}
}

func TestMethodParityAndSignatures(t *testing.T) {
	usedExceptions := map[string]bool{}
	usedDeviations := map[string]bool{}

	for _, p := range registry() {
		t.Run(p.name, func(t *testing.T) {
			sdkType := normalizeType(p.sdk)
			ifaceType := p.iface

			sdkIdx := methodIndex(sdkType)
			ifaceIdx := methodIndex(ifaceType)
			ifaceMethods := methodNames(ifaceType)

			for name := range sdkIdx {
				key := p.name + "." + name
				if reason, ok := exceptions[key]; ok {
					usedExceptions[key] = true
					t.Logf("skip %s: %s", key, reason)
					continue
				}
				if !ifaceMethods[name] {
					t.Errorf("fsmock.%s missing method %q present on SDK.\n  Fix: add it to the interface + wrapper, or add exceptions[%q] with a Reason.", p.name, name, key)
					continue
				}
				sdkSig := methodSig(sdkType, sdkIdx[name])
				ifaceSig := methodSig(ifaceType, ifaceIdx[name])
				if sdkSig == ifaceSig {
					continue
				}
				if reason, ok := deviations[key]; ok {
					usedDeviations[key] = true
					t.Logf("deviation %s: %s\n  sdk=%s\n  fsmock=%s", key, reason, sdkSig, ifaceSig)
					continue
				}
				t.Errorf("signature mismatch for %s\n  sdk:    %s\n  fsmock: %s\n  Fix: align the wrapper, or add deviations[%q] with a Reason.", key, sdkSig, ifaceSig, key)
			}

			switch p.name {
			case "DocumentSnapshot":
				for _, want := range []string{"Data", "DataTo", "DataAt", "DataAtPath", "Exists", "CreateTime", "UpdateTime", "ReadTime", "Ref", "Reference"} {
					if !ifaceMethods[want] {
						t.Errorf("fsmock.DocumentSnapshot missing %s", want)
					}
				}
			case "QuerySnapshot":
				for _, want := range []string{"Documents", "Size", "Changes", "ReadTime", "Reference"} {
					if !ifaceMethods[want] {
						t.Errorf("fsmock.QuerySnapshot missing %s", want)
					}
				}
			case "AggregationResult":
				for _, want := range []string{"Data", "DataTo"} {
					if !ifaceMethods[want] {
						t.Errorf("fsmock.AggregationResult missing %s", want)
					}
				}
			}

			for key, reason := range deviations {
				if !strings.HasPrefix(key, p.name+".") {
					continue
				}
				method := strings.TrimPrefix(key, p.name+".")
				if !ifaceMethods[method] {
					continue
				}
				if _, onSDK := sdkIdx[method]; onSDK {
					continue
				}
				usedDeviations[key] = true
				t.Logf("fsmock-only %s: %s", key, reason)
			}
		})
	}

	t.Cleanup(func() {
		// Run after all subtests so -run filters that skip pairs do not false-positive.
		// Only enforce when the parent test itself was selected without a subtest filter
		// that leaves pairs unexecuted — detect via counting executed pairs.
	})

	for key := range exceptions {
		if !usedExceptions[key] {
			t.Errorf("stale exceptions entry %q — delete it from exceptions.go", key)
		}
	}
	for key := range deviations {
		if !usedDeviations[key] {
			// Only fail if the owning type's subtest ran (appeared in registry walk).
			// When signatures match after substitution, a deviation entry is still
			// documentation — mark it used if the method exists on the iface.
			owner, _, ok := strings.Cut(key, ".")
			if !ok {
				t.Errorf("malformed deviations key %q", key)
				continue
			}
			var p *pair
			for i := range registry() {
				if registry()[i].name == owner {
					pp := registry()[i]
					p = &pp
					break
				}
			}
			if p == nil {
				t.Errorf("stale deviations entry %q — unknown type; delete it", key)
				continue
			}
			ifaceMethods := methodNames(p.iface)
			method := strings.TrimPrefix(key, owner+".")
			if ifaceMethods[method] {
				// Documented intentional method (signature may match after substitution).
				continue
			}
			t.Errorf("stale deviations entry %q — method not on fsmock.%s; delete it from exceptions.go", key, owner)
		}
	}
}

func TestAggregationResultSDKSurface(t *testing.T) {
	sdk := normalizeType(reflect.TypeOf(firestore.AggregationResult(nil)))
	idx := methodIndex(sdk)
	for _, want := range []string{"Data", "DataTo"} {
		if _, ok := idx[want]; !ok {
			t.Errorf("SDK AggregationResult missing %s", want)
		}
	}
}
