## Summary

-

## Checklist

- [ ] `make gate` passes locally (or CI `gate` is green)
- [ ] Mocks regenerated if interfaces changed (`make generate-check`)
- [ ] New wrapper methods have `fstest` coverage **or** a waiver with Reason+Issue
- [ ] apicheck clean (no new unexplained SDK methods / signature drift)
- [ ] If Firestore minor changed: `version.go` + `COMPATIBILITY.md` updated
- [ ] Docs (`README` / `CHANGELOG` / `AGENTS.md`) updated when behavior changes

## Test plan

- [ ] Unit / apicheck
- [ ] Accuracy (emulator)
