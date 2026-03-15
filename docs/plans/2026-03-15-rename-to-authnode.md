# Rename Project to authnode Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Rename the project from "studBridge" to "authnode" across all source files, module paths, and documentation.

**Architecture:** All Go modules use `github.com/PavelShe11/studbridge/...` as the import path prefix. Every source file, go.mod, proto file, and generated code must be updated to `github.com/PavelShe11/authnode/...`. Display names ("studBridge", "StudBridge") must also be updated.

**Tech Stack:** Go modules, Protocol Buffers (buf), swag (Swagger codegen), Go workspace (go.work)

---

### Task 1: Replace Go module paths in go.mod files

**Files:**
- Modify: `common/go.mod`
- Modify: `authMicro/go.mod`
- Modify: `authMicro/grpcApi/go.mod`
- Modify: `userMicro/go.mod`

**Step 1: Replace in all go.mod files at once**

```bash
find /Users/pavel/GolandProjects/studBridge -name "go.mod" -exec sed -i '' 's|github.com/PavelShe11/studbridge|github.com/PavelShe11/authnode|g' {} \;
```

**Step 2: Verify the changes**

```bash
grep -r "studbridge\|studBridge" /Users/pavel/GolandProjects/studBridge --include="go.mod"
```
Expected: no output (all replaced)

---

### Task 2: Replace import paths in all Go source files

**Files:**
- Modify: all `*.go` files under `authMicro/`, `userMicro/`, `common/`

**Step 1: Replace in all .go files**

```bash
find /Users/pavel/GolandProjects/studBridge -name "*.go" -exec sed -i '' 's|github.com/PavelShe11/studbridge|github.com/PavelShe11/authnode|g' {} \;
```

**Step 2: Verify no studbridge references remain in Go files**

```bash
grep -r "studbridge\|studBridge" /Users/pavel/GolandProjects/studBridge --include="*.go"
```
Expected: no output

---

### Task 3: Update proto files and regenerate gRPC code

**Files:**
- Modify: `authMicro/grpcApi/proto/account_service.proto`
- Modify: `authMicro/grpcApi/proto/error.proto`
- Regenerate: `authMicro/grpcApi/account_service.pb.go`, `authMicro/grpcApi/error.pb.go`, etc.

**Step 1: Replace go_package option in proto files**

```bash
find /Users/pavel/GolandProjects/studBridge -name "*.proto" -exec sed -i '' 's|github.com/PavelShe11/studbridge|github.com/PavelShe11/authnode|g' {} \;
```

**Step 2: Regenerate gRPC code**

```bash
cd /Users/pavel/GolandProjects/studBridge/authMicro/grpcApi && buf generate
```

**Step 3: Verify generated files use new path**

```bash
grep "studbridge" /Users/pavel/GolandProjects/studBridge/authMicro/grpcApi/*.go
```
Expected: no output

---

### Task 4: Update YAML configuration files

**Files:**
- Modify: `authMicro/.mockery.yaml`
- Modify: `authMicro/docs/swagger.yaml`

**Step 1: Replace in YAML files**

```bash
find /Users/pavel/GolandProjects/studBridge -name "*.yaml" -exec sed -i '' 's|github.com/PavelShe11/studbridge|github.com/PavelShe11/authnode|g' {} \;
find /Users/pavel/GolandProjects/studBridge -name "*.yaml" -exec sed -i '' 's|studBridge|authnode|g' {} \;
find /Users/pavel/GolandProjects/studBridge -name "*.yaml" -exec sed -i '' 's|StudBridge|AuthNode|g' {} \;
```

**Step 2: Also update swagger.json if it exists**

```bash
find /Users/pavel/GolandProjects/studBridge -name "*.json" -path "*/docs/*" -exec sed -i '' 's|github.com/PavelShe11/studbridge|github.com/PavelShe11/authnode|g' {} \;
find /Users/pavel/GolandProjects/studBridge -name "*.json" -path "*/docs/*" -exec sed -i '' 's|StudBridge|AuthNode|g' {} \;
```

**Step 3: Verify**

```bash
grep -r "studbridge\|studBridge\|StudBridge" /Users/pavel/GolandProjects/studBridge --include="*.yaml" --include="*.json"
```
Expected: no output

---

### Task 5: Update README.md

**Files:**
- Modify: `README.md`

**Step 1: Replace display names**

```bash
sed -i '' 's|studBridge|authnode|g' /Users/pavel/GolandProjects/studBridge/README.md
sed -i '' 's|StudBridge|AuthNode|g' /Users/pavel/GolandProjects/studBridge/README.md
```

**Step 2: Verify**

```bash
grep -i "studbridge" /Users/pavel/GolandProjects/studBridge/README.md
```
Expected: no output (or only in historical context if intentionally kept)

---

### Task 6: Update go.sum files and workspace

**Files:**
- Modify: `authMicro/go.sum`, `userMicro/go.sum`
- Read: `go.work`

**Note:** The go.sum files contain checksums for the old published GitHub packages (`github.com/PavelShe11/studbridge/...`). These entries referenced published versions from GitHub. After renaming, inter-module references are resolved by the Go workspace (`go.work` uses `use` directives), so the published package entries become stale. Run `go work sync` to update.

**Step 1: Sync workspace**

```bash
cd /Users/pavel/GolandProjects/studBridge && go work sync
```

**Step 2: Tidy each module**

```bash
cd /Users/pavel/GolandProjects/studBridge/common && go mod tidy
cd /Users/pavel/GolandProjects/studBridge/authMicro/grpcApi && go mod tidy
cd /Users/pavel/GolandProjects/studBridge/authMicro && go mod tidy
cd /Users/pavel/GolandProjects/studBridge/userMicro && go mod tidy
```

**Step 3: Verify builds compile**

```bash
cd /Users/pavel/GolandProjects/studBridge/authMicro && go build ./...
cd /Users/pavel/GolandProjects/studBridge/userMicro && go build ./...
```
Expected: no errors

---

### Task 7: Commit

**Step 1: Stage and commit all changes**

```bash
cd /Users/pavel/GolandProjects/studBridge
git add -A
git commit -m "chore: rename project from studBridge to authnode"
```
