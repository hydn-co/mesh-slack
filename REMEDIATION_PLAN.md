# mesh-slack PR #3 Remediation Plan

**PR**: Develop ([hydn-co/mesh-slack/pull/3](https://github.com/hydn-co/mesh-slack/pull/3))  
**Issues Found**: 5 open code review comments  
**Date**: 2026-04-28

---

## Summary

This PR adds fully functional Slack collectors and actions backed by a shared Slack API client, replacing stubbed implementations. The code review identified **5 critical issues** spanning:
- Unicode handling and message validation
- Resource management (HTTP connection leaks)
- Performance (unnecessary API calls, memory inefficiency)
- Code correctness (unnecessary allocations, error handling)

---

## Phase 1: Resource Management & Connection Leaks (Critical)

**Priority**: HIGH  
**Impact**: Production stability; prevents HTTP connection exhaustion  
**Estimated Effort**: 1–2 hours

### Issue 1.1: HTTP Response Body Leak in `slack_api/client.go`

**File**: `internal/slack_api/client.go`  
**Function**: `readResponseBody`  
**Problem**: Early return on context cancellation without closing `resp.Body` causes HTTP connection leaks.

**Current Code Pattern**:
```go
func readResponseBody(ctx context.Context, resp *http.Response) ([]byte, error) {
    if ctx.Err() != nil {
        return nil, ctx.Err()  // ❌ resp.Body not closed
    }
    // ...
    return body, nil
}
```

**Solution**:
- Add defensive `resp` nil check before body close
- Ensure `resp.Body.Close()` is called on all exit paths (including context cancellation)
- Use `defer` to guarantee cleanup

**Expected Fix**:
```go
func readResponseBody(ctx context.Context, resp *http.Response) ([]byte, error) {
    if resp == nil {
        return nil, errors.New("response is nil")
    }
    defer resp.Body.Close()
    
    if ctx.Err() != nil {
        return nil, ctx.Err()
    }
    // ... rest of function
}
```

**Validation**:
- [ ] Run `go vet ./internal/slack_api/...` – should pass
- [ ] Run `go test ./internal/slack_api/...` – connection leak tests pass
- [ ] Manual review: trace all return paths, confirm defer covers all

---

## Phase 2: Performance Optimization (High)

**Priority**: HIGH  
**Impact**: API quota efficiency, memory footprint, latency  
**Estimated Effort**: 2–3 hours total

### Issue 2.1: Unnecessary Slack API Call in `users/resolve_user_id.go`

**File**: `internal/users/resolve_user_id.go`  
**Function**: `ResolveUserIDsByEmails`  
**Problem**: Makes a ListUsers API call even when `emails` is empty; loop break is checked after first page fetch.

**Current Code Pattern**:
```go
func ResolveUserIDsByEmails(ctx context.Context, client *slack.Client, emails []string) (map[string]string, error) {
    result := make(map[string]string, len(emails))
    cursor := ""
    for {
        // ... fetch page
        // ... process results
        if cursor == "" {
            break  // ❌ break checked AFTER fetch, not before loop entry
        }
    }
    return result, nil
}
```

**Solution**:
- Add early return for `len(emails) == 0` before any API calls
- Clarifies function behavior: "if no emails provided, return empty map immediately"

**Expected Fix**:
```go
func ResolveUserIDsByEmails(ctx context.Context, client *slack.Client, emails []string) (map[string]string, error) {
    if len(emails) == 0 {
        return make(map[string]string), nil
    }
    result := make(map[string]string, len(emails))
    cursor := ""
    for {
        // ... (rest unchanged)
    }
}
```

**Validation**:
- [ ] Add unit test: `TestResolveUserIDsByEmails_EmptyInput` verifying no API calls made
- [ ] Run `go test ./internal/users/...` – test passes
- [ ] Verify in mock client that `ListUsers` is never called with empty input

---

### Issue 2.2: Unnecessary Map Clone in `slack_api/client.go`

**File**: `internal/slack_api/client.go`  
**Function**: `describeSlackError`  
**Problem**: Clones `slackErrorDescriptions` on every call but never mutates it; wastes allocation on error path.

**Current Code Pattern**:
```go
var slackErrorDescriptions = map[string]string{
    "invalid_auth":    "Authentication failed: invalid token or permissions",
    "team_access_not_granted": "...",
    // ...
}

func describeSlackError(code string) string {
    clone := maps.Clone(slackErrorDescriptions)  // ❌ unnecessary allocation
    if desc, ok := clone[code]; ok {
        return desc
    }
    return "Unknown Slack API error"
}
```

**Solution**:
- Direct map lookup without clone; map is read-only
- Eliminates allocation on every error

**Expected Fix**:
```go
func describeSlackError(code string) string {
    if desc, ok := slackErrorDescriptions[code]; ok {
        return desc
    }
    return "Unknown Slack API error"
}
```

**Validation**:
- [ ] Run `go vet ./internal/slack_api/...` – should pass
- [ ] Run benchmarks if error path is hot: `go test -bench=. ./internal/slack_api/...`
- [ ] Verify no race conditions (slackErrorDescriptions is const/read-only)

---

### Issue 2.3: Memory-Inefficient Channel Emission in `slack_channels_collector.go`

**File**: `internal/collectors/slack_channels_collector.go`  
**Function**: `SlackChannelsCollector.Start`  
**Problem**: Collects all channels into memory before emitting; expensive for large workspaces (100k+ channels).

**Current Code Pattern**:
```go
func (c *SlackChannelsCollector) Start(ctx context.Context) error {
    var channels []*Channels  // ❌ allocates for ALL channels
    err := c.enumerators.ForEach(func(ch *Channel) error {
        channels = append(channels, ch)  // collect all in slice
        return nil
    })
    // after loop, emit all at once
    for _, ch := range channels {
        c.emit(ch)
    }
}
```

**Solution**:
- Emit channels inside the `ForEach` callback for streaming behavior
- Avoids buffering entire workspace; processes incrementally

**Expected Fix**:
```go
func (c *SlackChannelsCollector) Start(ctx context.Context) error {
    return c.enumerators.ForEach(func(ch *Channel) error {
        // Emit directly, don't buffer
        c.emit(ch)
        return nil
    })
}
```

**Validation**:
- [ ] Add memory profile test: compare peak memory before/after fix
- [ ] Run `go test ./internal/collectors/...` – all pass
- [ ] Manual review: verify emit callback is idempotent and safe to call repeatedly
- [ ] Verify error handling: if emit fails, does ForEach propagate error correctly?

---

## Phase 3: Correctness & Character Encoding (High)

**Priority**: HIGH  
**Impact**: Message validation correctness; prevents incorrect rejections of valid Unicode  
**Estimated Effort**: 2–3 hours

### Issue 3.1: Byte Length vs. Rune Length in `slack_channel_message_post_action.go`

**File**: `internal/actions/slack_channel_message_post_action.go`  
**Function**: `verifyMessage`  
**Problem**: Uses `len(trimmed)` which counts bytes, not Unicode characters. Slack's 4000-character limit applies to runes, not bytes.

**Example**:
- Emoji (4-byte UTF-8): `len("😀") == 4` but should count as `1` character
- Accented characters: `len("café") == 5` but should count as `4` characters

**Current Code Pattern**:
```go
func verifyMessage(msg string) error {
    trimmed := strings.TrimSpace(msg)
    if len(trimmed) > 4000 {  // ❌ counts bytes, not runes
        return fmt.Errorf("message exceeds 4000 character limit: %d bytes", len(trimmed))
    }
    return nil
}
```

**Solution**:
- Use `utf8.RuneCountInString` to count actual characters
- Update error message to report character count, not byte count
- Add tests for multi-byte sequences (emoji, accented chars, CJK)

**Expected Fix**:
```go
import "unicode/utf8"

func verifyMessage(msg string) error {
    trimmed := strings.TrimSpace(msg)
    charCount := utf8.RuneCountInString(trimmed)
    if charCount > 4000 {
        return fmt.Errorf("message exceeds 4000 character limit: %d characters", charCount)
    }
    return nil
}
```

**Validation**:
- [ ] Add test cases:
  - Pure ASCII: `strings.Repeat("a", 4001)` should fail
  - Emoji: `strings.Repeat("😀", 4001)` should fail (emit exact char count)
  - Mixed: ASCII + emoji + accents, verify count is correct
- [ ] Run `go test ./internal/actions/... -v` – all pass
- [ ] Verify error messages are user-friendly and accurate

---

## Phase 4: Integration Testing & Validation (Medium)

**Priority**: MEDIUM  
**Impact**: Confidence in fixes; prevents regressions  
**Estimated Effort**: 3–4 hours

### Test Suite Enhancements

**New Integration Tests**:
1. **Connection Leak Prevention**:
   - Mock HTTP server that cancels context mid-response
   - Verify connection is closed, not left in TIME_WAIT

2. **API Call Minimization**:
   - Mock Slack client tracking call count
   - Verify `ResolveUserIDsByEmails([])` makes 0 API calls
   - Verify `ResolveUserIDsByEmails(["a@b.com"])` makes 1 API call

3. **Streaming Emission**:
   - Mock enumerator with 10k channels
   - Profile memory usage: verify no full list buffered in memory
   - Verify callback is invoked for each channel incrementally

4. **Unicode Validation**:
   - Test matrix: ASCII, emoji, CJK, accents, combining diacritics
   - Verify all pass/fail boundaries at exactly 4000 characters

---

## Phase 5: Go Best Practices Audit (Low)

**Priority**: LOW  
**Impact**: Code maintainability, consistency with mesh standards  
**Estimated Effort**: 1–2 hours

### Suggested Enhancements (Optional but Recommended)

1. **Error Wrapping**: Use `fmt.Errorf("%w", err)` for error context
2. **Nil Checks**: Defensive checks for nil pointers before dereferencing
3. **Test Coverage**: Ensure all error paths have tests
4. **Comments**: Document non-obvious logic (e.g., why certain limits exist)
5. **Linting**: Run `golangci-lint run ./...` to catch any missed issues

---

## Implementation Roadmap

### Timeline: ~10–12 hours total

| Phase | Duration | Order | Status |
|-------|----------|-------|--------|
| Phase 1: Resource Mgmt | 1–2h | **First** | Not started |
| Phase 2: Performance | 2–3h | **Second** | Not started |
| Phase 3: Correctness | 2–3h | **Third** | Not started |
| Phase 4: Testing | 3–4h | **Fourth** | Not started |
| Phase 5: Audit | 1–2h | **Optional** | Not started |

---

## Verification Checklist

- [ ] All 5 issues understood and triaged
- [ ] Phase 1 complete: no connection leaks
- [ ] Phase 2 complete: no unnecessary API calls or allocations
- [ ] Phase 3 complete: Unicode handling correct
- [ ] Phase 4 complete: new tests pass, no regressions
- [ ] `go vet ./...` passes
- [ ] `go test ./...` passes with 100% relevant coverage
- [ ] `golangci-lint run ./...` clean (if configured)
- [ ] PR review comments resolved and conversations marked done

---

## Notes

- All fixes are backward-compatible; no API changes required
- Fixes improve stability (phase 1), efficiency (phase 2), and correctness (phase 3)
- Tests provide confidence and prevent regressions
- Aligns with Go best practices and mesh project standards

