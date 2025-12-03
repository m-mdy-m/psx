# Non-Functional Requirements
## PSX - Project Structure Checker

- **Version:** 1.0.0  
- **Author:** Genix (m-mdy-m)  
- **Contact:** bitsgenix@gmail.com  
- **Repository:** [m-mdy-m/psx](https://github.com/m-mdy-m/psx)
- **Date:** 2025-12-02
- **Status:** Draft  
- **Confidentiality:** Internal

---

## Overview

This document specifies the non-functional requirements (NFRs) for PSX. These requirements describe HOW the system should behave rather than WHAT it should do.

Categories:
1. Performance
2. Reliability
3. Usability
4. Security
5. Maintainability
6. Portability
7. Scalability

---

## 1. Performance Requirements

### NFR-1.1: Execution Speed
**Priority:** P0  
**Category:** Performance

**Requirements:**

**Small Projects (< 1,000 files):**
- Scan time: <1 second
- Config load: <50ms
- Report generation: <100ms
- Total execution: <1.5 seconds

**Medium Projects (1,000 - 10,000 files):**
- Scan time: <5 seconds
- Config load: <50ms
- Report generation: <200ms
- Total execution: <6 seconds

**Large Projects (10,000 - 100,000 files):**
- Scan time: <30 seconds
- Config load: <100ms
- Report generation: <500ms
- Total execution: <35 seconds

**Acceptance Criteria:**
- [ ] Meets timing targets on standard hardware
- [ ] Performance degrades linearly with project size
- [ ] Progress indicator shown for operations >2s
- [ ] Can abort long operations with Ctrl+C

**Test Method:**
```bash
# Benchmark script
time psx check --config bench.yml

# Should complete within targets
# Tested on: 2-core CPU, 4GB RAM
```

---

### NFR-1.2: Memory Usage
**Priority:** P0  
**Category:** Performance

**Requirements:**

**Small Projects:**
- Peak memory: <50MB
- Baseline: <20MB

**Medium Projects:**
- Peak memory: <100MB
- Baseline: <30MB

**Large Projects:**
- Peak memory: <200MB
- Baseline: <50MB

**Memory Efficiency:**
- No memory leaks
- Garbage collection efficient
- Streaming for large files (>10MB)
- Don't load entire project into memory

**Acceptance Criteria:**
- [ ] Memory usage within limits
- [ ] No memory leaks detected
- [ ] Handles files up to 100MB
- [ ] Works on low-memory systems (2GB RAM)

**Test Method:**
```bash
# Memory profiling
go test -memprofile mem.out
go tool pprof mem.out

# Valgrind on Linux
valgrind --leak-check=full ./psx check
```

---

### NFR-1.3: Binary Size
**Priority:** P1  
**Category:** Performance

**Requirements:**
- Binary size: <20MB (compressed)
- Binary size: <50MB (uncompressed)
- No external dependencies
- Includes all templates and resources

**Optimization:**
- Strip debug symbols in release builds
- Compress with UPX (optional)
- Dead code elimination
- Template compression

**Acceptance Criteria:**
- [ ] Linux binary <15MB
- [ ] macOS binary <16MB
- [ ] Windows binary <18MB
- [ ] Download time <10s on 10Mbps connection

---

### NFR-1.4: Startup Time
**Priority:** P1  
**Category:** Performance

**Requirements:**
- Cold start: <100ms
- Config load: <50ms
- First scan: <200ms additional
- Help display: <50ms

**Optimization:**
- Lazy loading of rules
- Config parsing optimized
- Minimal initialization

**Acceptance Criteria:**
- [ ] `psx --version` responds instantly (<10ms)
- [ ] `psx --help` shows in <50ms
- [ ] First scan starts within 200ms
- [ ] No noticeable delay on fast machines

---

### NFR-1.5: Parallelization
**Priority:** P1  
**Category:** Performance

**Requirements:**
- Rule execution parallelized
- File scanning concurrent
- Uses all available CPU cores
- No race conditions

**Concurrency Model:**
- Worker pool for rule execution
- Channel-based communication
- Context for cancellation
- Proper synchronization

**Acceptance Criteria:**
- [ ] Scales with CPU cores
- [ ] No race conditions (tested with -race flag)
- [ ] Graceful degradation on single-core systems
- [ ] Can limit parallelism with flag

**Test Method:**
```bash
# Race detection
go test -race ./...

# CPU usage should scale
# 1 core: ~100% usage
# 4 cores: ~400% usage
```

---

## 2. Reliability Requirements

### NFR-2.1: Error Handling
**Priority:** P0  
**Category:** Reliability

**Requirements:**

**Graceful Degradation:**
- Continue on non-critical errors
- Report errors clearly
- Suggest fixes when possible
- Never crash unexpectedly

**Error Types:**

**Recoverable Errors (warn and continue):**
- Permission denied on single file
- Malformed config (use defaults)
- Missing optional files
- Network timeout (if future features need it)

**Fatal Errors (exit cleanly):**
- Cannot read project root
- Config syntax completely invalid
- Out of memory
- Disk full on write operations

**Error Messages:**
- Clear description of what happened
- Location information (file, line)
- Suggestion for fixing
- Exit code indicates error type

**Acceptance Criteria:**
- [ ] No panics in production code
- [ ] All errors logged appropriately
- [ ] User-friendly error messages
- [ ] Errors include actionable suggestions

**Example Error Message:**
```
Error: Cannot read project configuration

File: /path/to/psx.yml
Line: 15
Problem: Invalid YAML syntax - unexpected character

Expected: key: value
Found: key value (missing colon)

Suggestion: Add a colon after 'key'
Documentation: https://psx.dev/docs/config

Exit code: 2
```

---

### NFR-2.2: Data Integrity
**Priority:** P0  
**Category:** Reliability

**Requirements:**

**File Operations:**
- Atomic writes (temp file + rename)
- Verify write success
- Rollback on failure
- Never corrupt existing files

**Safety Measures:**
- Check disk space before writing
- Validate content before writing
- Create backups (optional)
- Checksums for critical operations

**Acceptance Criteria:**
- [ ] No partial writes
- [ ] Files never corrupted
- [ ] Rollback works correctly
- [ ] Handles disk full gracefully

---

### NFR-2.3: Deterministic Behavior
**Priority:** P0  
**Category:** Reliability

**Requirements:**
- Same input = same output
- Rule execution order consistent
- No random behavior
- Timestamps excluded from comparison

**Testing:**
- Run twice on same project = identical results
- Parallel execution = same results
- Different OS = same validation results
- Different machines = same results

**Acceptance Criteria:**
- [ ] Deterministic test passes 100 times
- [ ] Results identical across platforms
- [ ] No timestamp dependencies in logic
- [ ] Randomness only for non-critical features (like example data)

---

### NFR-2.4: Fault Tolerance
**Priority:** P1  
**Category:** Reliability

**Requirements:**

**Handle Gracefully:**
- Missing files
- Incorrect permissions
- Disk full
- Network interruption (future)
- Corrupted config files
- Large files (>1GB)
- Deep directory structures (>100 levels)
- Symbolic link loops

**Recovery:**
- Retry transient failures (max 3 times)
- Fallback to defaults on config errors
- Skip problematic files with warning
- Continue processing when possible

**Acceptance Criteria:**
- [ ] Handles all listed scenarios
- [ ] Provides helpful error messages
- [ ] Doesn't hang or crash
- [ ] Logs issues for debugging

---

## 3. Usability Requirements

### NFR-3.1: Ease of Learning
**Priority:** P0  
**Category:** Usability

**Requirements:**
- New users productive in 5 minutes
- Help text comprehensive
- Examples provided
- Error messages educational

**Documentation:**
- Quick start guide (<5 minutes to read)
- Command examples for common tasks
- Troubleshooting section
- Video tutorial (optional, future)

**First-time Experience:**
```bash
# User runs PSX for first time
psx check

# Output includes:
# - What PSX found
# - What's wrong
# - How to fix it
# - Where to get help
```

**Acceptance Criteria:**
- [ ] Users can run basic check without reading docs
- [ ] Help text answers common questions
- [ ] Error messages include documentation links
- [ ] Examples cover 80% of use cases

---

### NFR-3.2: Ease of Use
**Priority:** P0  
**Category:** Usability

**Requirements:**

**Simplicity:**
- Default behavior is sensible
- No required configuration for basic usage
- Common tasks are simple commands
- Advanced features don't complicate basics

**Command Structure:**
```bash
# Simple (most common)
psx check
psx fix

# With options (common)
psx check --verbose
psx fix --interactive

# Advanced (rare)
psx check --config custom.yml --level error --output json
```

**Discoverability:**
- `--help` on every command
- Suggestions when command unknown
- Tab completion available
- Flag typos detected and suggested

**Acceptance Criteria:**
- [ ] 90% of users need only 3 commands
- [ ] No mandatory flags for basic usage
- [ ] Help text always available
- [ ] Typos suggested correctly

---

### NFR-3.3: Feedback Quality
**Priority:** P0  
**Category:** Usability

**Requirements:**

**Clear Output:**
- Severity visible (colors, symbols)
- Grouped by category
- Most important issues first
- Summary at the end

**Progress Indication:**
- Show what's happening
- Progress bar for long operations
- Percentage complete
- Estimated time remaining

**Actionable Messages:**
- What's wrong
- Why it matters
- How to fix it
- Command to run

**Example:**
```
✗ README_MISSING

What: No README.md file found in project root
Why: README is the first thing people see
How: Create a README describing your project

Quick fix: psx fix --rule README_MISSING
Manual fix: Create README.md and add:
  - Project name and description
  - Installation instructions
  - Usage examples
  - License information

Learn more: https://psx.dev/docs/rules/readme
```

**Acceptance Criteria:**
- [ ] Every error has clear message
- [ ] Suggestions are actionable
- [ ] Links to documentation work
- [ ] Messages are helpful, not condescending

---

### NFR-3.4: Accessibility
**Priority:** P1  
**Category:** Usability

**Requirements:**

**Terminal Compatibility:**
- Works in all major terminals
- Supports screen readers (through text output)
- Color-blind friendly (shapes + colors)
- Works in 80-column width

**Internationalization (Future):**
- English in v1.0
- i18n framework in place
- Error codes language-independent

**Acceptance Criteria:**
- [ ] Works in: bash, zsh, fish, PowerShell, cmd
- [ ] Color palette accessible (WCAG AA)
- [ ] Text-only mode available
- [ ] Works over SSH

---

## 4. Security Requirements

### NFR-4.1: Input Validation
**Priority:** P0  
**Category:** Security

**Requirements:**

**Config File:**
- YAML parsing safe (no code execution)
- Path validation (no directory traversal)
- Size limits (config <10MB)
- Regex validation (prevent ReDoS)

**User Input:**
- Command line args sanitized
- File paths validated
- No shell injection
- No command injection

**Validation Rules:**
```go
// Path validation
func validatePath(path string) error {
    // No absolute paths escaping project
    // No parent directory traversal
    // No symbolic link escaping
    // Maximum depth limit
}
```

**Acceptance Criteria:**
- [ ] All inputs validated
- [ ] No code execution from config
- [ ] Path traversal prevented
- [ ] ReDoS protection in place

---

### NFR-4.2: File System Safety
**Priority:** P0  
**Category:** Security

**Requirements:**

**Read Operations:**
- Respect file permissions
- Don't follow symlinks to sensitive areas
- Size limits (don't read >100MB files)
- Timeout on slow operations

**Write Operations:**
- Only in project directory
- Create with safe permissions (644 for files, 755 for dirs)
- Never overwrite without confirmation
- Atomic operations only

**Forbidden Operations:**
- No deleting files
- No executing files
- No network access
- No system calls (except filesystem)

**Acceptance Criteria:**
- [ ] Cannot escape project directory
- [ ] Cannot delete files
- [ ] Cannot execute code
- [ ] Cannot access sensitive files (/etc/passwd, etc.)

---

### NFR-4.3: No Telemetry
**Priority:** P0  
**Category:** Security/Privacy

**Requirements:**
- No data collection
- No network requests (in v1.0)
- No tracking
- No crash reporting to external services

**User Privacy:**
- Project structure never sent anywhere
- No analytics
- No update checks (in v1.0)
- Fully offline capable

**Acceptance Criteria:**
- [ ] No network code in binary
- [ ] Verified with network monitoring
- [ ] Privacy policy clear
- [ ] Open source (auditable)

---

### NFR-4.4: Dependency Security
**Priority:** P1  
**Category:** Security

**Requirements:**
- Minimal dependencies
- All dependencies audited
- Regular security updates
- No deprecated packages

**Current Dependencies (example):**
- Go standard library (trusted)
- gopkg.in/yaml.v3 (well-maintained)
- github.com/spf13/cobra (widely used)
- github.com/fatih/color (simple, safe)

**Acceptance Criteria:**
- [ ] All dependencies security-scanned
- [ ] Dependabot enabled
- [ ] No known vulnerabilities
- [ ] Dependencies kept up to date

---

## 5. Maintainability Requirements

### NFR-5.1: Code Quality
**Priority:** P0  
**Category:** Maintainability

**Requirements:**

**Code Standards:**
- Go idiomatic code
- `gofmt` formatted
- `golangci-lint` passes
- No code smells

**Metrics:**
- Cyclomatic complexity <15 per function
- Function length <100 lines
- File length <1000 lines
- Comment ratio >20%

**Documentation:**
- All public functions documented
- Complex logic explained
- Examples in godoc
- README for each package

**Acceptance Criteria:**
- [ ] Linter passes with no warnings
- [ ] Code review checklist satisfied
- [ ] All metrics within targets
- [ ] Documentation complete

---

### NFR-5.2: Test Coverage
**Priority:** P0  
**Category:** Maintainability

**Requirements:**

**Coverage Targets:**
- Unit tests: >80% coverage
- Integration tests: All major workflows
- E2E tests: All commands
- Performance tests: Key operations

**Test Quality:**
- Tests are deterministic
- Tests are fast (<5s for unit tests)
- Tests are isolated
- Tests have clear names

**Test Structure:**
```
test/
├── unit/              # Fast, isolated tests
├── integration/       # Test component integration
├── e2e/              # Full end-to-end scenarios
├── fixtures/         # Test data
└── benchmarks/       # Performance tests
```

**Acceptance Criteria:**
- [ ] Coverage >80% on all packages
- [ ] All critical paths tested
- [ ] CI runs all tests
- [ ] Tests pass on all platforms

---

### NFR-5.3: Code Organization
**Priority:** P1  
**Category:** Maintainability

**Requirements:**

**Clear Structure:**
- Logical package organization
- Clear separation of concerns
- Internal vs public API clear
- Dependencies flow downward

**Package Layout:**
```
cmd/          # Entry points
internal/     # Private implementation
  cli/        # Command line interface
  detector/   # Project detection
  rules/      # Rule engine
  fixer/      # Auto-fix logic
pkg/          # Public API
  schema/     # Data structures
  util/       # Utilities
```

**Acceptance Criteria:**
- [ ] No circular dependencies
- [ ] Clear package responsibilities
- [ ] Easy to find code
- [ ] New developers understand structure quickly

---

### NFR-5.4: Extensibility
**Priority:** P1  
**Category:** Maintainability

**Requirements:**

**Easy to Add:**
- New rules (plugin interface)
- New project types (detector interface)
- New output formats (formatter interface)
- New templates (template system)

**Design Patterns:**
- Strategy pattern for rules
- Factory pattern for detectors
- Observer pattern for events
- Template method for common logic

**Example - Adding New Rule:**
```go
// 1. Implement Rule interface
type MyNewRule struct{}

func (r *MyNewRule) Check(proj *Project) *Result {
    // Rule logic
}

// 2. Register rule
registry.Register("MY_NEW_RULE", &MyNewRule{})

// 3. Done! No other changes needed
```

**Acceptance Criteria:**
- [ ] Can add rule without modifying core
- [ ] Can add detector without breaking existing
- [ ] Can add formatter easily
- [ ] Examples documented

---

## 6. Portability Requirements

### NFR-6.1: Cross-Platform
**Priority:** P0  
**Category:** Portability

**Requirements:**

**Operating Systems:**
- Linux (Ubuntu 20.04+, Debian 10+, Fedora, Arch)
- macOS (11.0+, Intel and Apple Silicon)
- Windows (10+, Server 2019+)
- FreeBSD (best effort)

**Architectures:**
- amd64 (x86_64)
- arm64 (aarch64)
- 386 (x86) - if needed

**Platform Differences Handled:**
- Path separators (/ vs \)
- Line endings (LF vs CRLF)
- File permissions
- Case sensitivity
- Home directory location

**Acceptance Criteria:**
- [ ] Same behavior on all platforms
- [ ] CI tests on all platforms
- [ ] Binaries for all platforms
- [ ] Platform-specific code isolated

---

### NFR-6.2: Shell Compatibility
**Priority:** P1  
**Category:** Portability

**Requirements:**

**Unix Shells:**
- bash
- zsh
- fish
- sh (POSIX)

**Windows Shells:**
- PowerShell
- cmd.exe
- Git Bash

**Features:**
- Exit codes work correctly
- Output formatting correct
- Colors work (or disable gracefully)
- Completion scripts available

**Acceptance Criteria:**
- [ ] Tested in all major shells
- [ ] Completion works in bash/zsh/fish
- [ ] PowerShell experience good
- [ ] Works over SSH

---

### NFR-6.3: Container Support
**Priority:** P1  
**Category:** Portability

**Requirements:**

**Docker Images:**
- Official image on Docker Hub
- Minimal size (<50MB)
- Includes all features
- Regular updates

**Image Variants:**
```dockerfile
# Alpine-based (smallest)
FROM alpine:3.19
# Size: ~30MB

# Scratch-based (most minimal)
FROM scratch
# Size: ~20MB

# Distroless (secure)
FROM gcr.io/distroless/static
# Size: ~25MB
```

**Usage:**
```bash
docker run --rm -v $(pwd):/project psx/psx:latest check
```

**Acceptance Criteria:**
- [ ] Docker image <50MB
- [ ] Works in Docker and Podman
- [ ] Kubernetes compatible
- [ ] CI/CD examples provided

---

## 7. Scalability Requirements

### NFR-7.1: Project Size
**Priority:** P0  
**Category:** Scalability

**Requirements:**

**File Count:**
- Small: 1-1,000 files ✓
- Medium: 1,000-10,000 files ✓
- Large: 10,000-100,000 files ✓
- Huge: 100,000+ files (degraded performance OK)

**Directory Depth:**
- Handle up to 100 levels deep
- Detect circular symlinks
- Warn on excessive depth

**Individual File Size:**
- Small files (<1MB): Read fully
- Medium files (1-10MB): Stream
- Large files (10-100MB): Stream with limit
- Huge files (>100MB): Skip with warning

**Acceptance Criteria:**
- [ ] Tested with 100k file project
- [ ] Memory usage linear
- [ ] No stack overflow on deep trees
- [ ] Performance acceptable for realistic projects

---

### NFR-7.2: Rule Count
**Priority:** P1  
**Category:** Scalability

**Requirements:**
- Built-in rules: 20-30 in v1.0
- Support: 100+ rules without performance impact
- Custom rules: 50+ supported
- Rule execution parallel

**Performance:**
- 30 rules: <1s for small projects
- 100 rules: <3s for small projects
- Linear scaling with rule count

**Acceptance Criteria:**
- [ ] Can handle 100 rules
- [ ] Performance remains acceptable
- [ ] No memory issues with many rules
- [ ] Parallelization effective

---

### NFR-7.3: Concurrent Usage
**Priority:** P1  
**Category:** Scalability

**Requirements:**

**CI/CD:**
- Multiple projects checked simultaneously
- No resource contention
- Each run isolated
- No shared state

**Multi-user:**
- Each user has own config
- No cross-user interference
- Concurrent safe

**Acceptance Criteria:**
- [ ] Can run 10 instances simultaneously
- [ ] No race conditions
- [ ] No file locking issues
- [ ] No resource exhaustion

---

## Acceptance Criteria Summary

For v1.0 release, the following must be verified:

**Performance:**
- [ ] All timing targets met
- [ ] Memory usage within limits
- [ ] Binary size <20MB
- [ ] Startup time <100ms

**Reliability:**
- [ ] No crashes in testing
- [ ] Error handling complete
- [ ] Deterministic behavior verified
- [ ] Fault tolerance tested

**Usability:**
- [ ] User testing completed
- [ ] Help text comprehensive
- [ ] Error messages clear
- [ ] Accessible to all users

**Security:**
- [ ] Security audit passed
- [ ] No vulnerabilities
- [ ] Input validation complete
- [ ] Privacy requirements met

**Maintainability:**
- [ ] Code quality metrics met
- [ ] Test coverage >80%
- [ ] Documentation complete
- [ ] Extensibility verified

**Portability:**
- [ ] All platforms tested
- [ ] Docker images available
- [ ] Shell compatibility verified

**Scalability:**
- [ ] Large project tested
- [ ] Many rules tested
- [ ] Concurrent usage tested

---

## Testing Strategy

### Performance Testing
```bash
# Benchmarks
go test -bench=. -benchmem ./...

# Large project test
psx check --config stress-test.yml /path/to/huge/project

# Memory profiling
go test -memprofile mem.out
go tool pprof mem.out
```

### Reliability Testing
```bash
# Chaos testing
# - Kill process randomly
# - Fill disk during write
# - Remove permissions during scan
# - Corrupt config file

# Fuzzing
go test -fuzz=FuzzConfigParser
```

### Security Testing
```bash
# Static analysis
golangci-lint run
gosec ./...

# Dependency check
go list -m all | nancy sleuth

# Manual audit
# - Check all file operations
# - Review input validation
# - Test path traversal attempts
```

---

**Document Version:** 1.0.0  
**Status:** Draft  
**Next Review:** Before v1.0 release

*Non-functional requirements are as important as functional requirements. They ensure PSX is not just correct, but also fast, reliable, and pleasant to use.*