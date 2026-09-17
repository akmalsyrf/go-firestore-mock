# TODO

Pinned SDK: **cloud.google.com/go/firestore v1.25.0** (`go.mod`).

PR gate = job **`gate`** (lint + unit/apicheck + generate-check + **accuracy**).

Accuracy = `scripts/check-accuracy.sh`: apicheck + required integration/parity tests; **skip = fail**.

## Follow-ups

- [ ] Tag `v2.0.0` after `gate` is green on default branch
- [ ] Branch protection: require status check named `gate`
- [ ] `go get cloud.google.com/go/firestore@latest && go generate && make gate` on SDK bumps
