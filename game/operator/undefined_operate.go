package operator

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"lucksystem/game/engine"
	"lucksystem/game/runtime"
)

type LucaOperateUndefined struct {
}

// UndefinedOpcodeTracker tracks undefined opcodes encountered during script processing.
// It accumulates counts silently and provides a summary at the end.
var undefinedTracker = &opcodeTracker{
	counts: make(map[string]int),
}

type opcodeTracker struct {
	mu     sync.Mutex
	counts map[string]int
	total  int
}

// Track records an occurrence of an undefined opcode
func (t *opcodeTracker) Track(opcode string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.counts[opcode]++
	t.total++
}

// Reset clears all tracked data
func (t *opcodeTracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.counts = make(map[string]int)
	t.total = 0
}

// Summary returns a formatted summary string, or empty if no undefined opcodes
func (t *opcodeTracker) Summary() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.total == 0 {
		return ""
	}

	// Sort opcodes by count (descending)
	type entry struct {
		name  string
		count int
	}
	entries := make([]entry, 0, len(t.counts))
	for name, count := range t.counts {
		entries = append(entries, entry{name, count})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].count > entries[j].count
	})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[INFO] %d undefined opcodes skipped (%d unique types):\n", t.total, len(t.counts)))
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("  %-20s  x%d\n", e.name, e.count))
	}
	sb.WriteString("These are non-text opcodes (visual/audio/system) and can be safely ignored for translation work.")
	return sb.String()
}

// PrintUndefinedOpcodeSummary prints the accumulated summary of undefined opcodes
// and resets the tracker. Call this after RunScript() completes.
func PrintUndefinedOpcodeSummary() {
	summary := undefinedTracker.Summary()
	if summary != "" {
		fmt.Println(summary)
	}
	undefinedTracker.Reset()
}

func (g *LucaOperateUndefined) UNDEFINED(ctx *runtime.Runtime, opcode string) engine.HandlerFunc {
	// Track the undefined opcode silently (no per-opcode console output)
	code := ctx.Code()
	if len(opcode) == 0 {
		opcode = ToString("%X", code.Opcode)
	}
	undefinedTracker.Track(opcode)

	list, end := AllToUint16(code.ParamBytes)
	if end >= 0 {
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			list,
			code.ParamBytes[end],
		)
	} else {
		ctx.Script.SetOperateParams(ctx.CIndex, ctx.RunMode,
			list,
		)
	}
	return func() {
		// 下一步执行地址，为0则表示紧接着向下
		ctx.ChanEIP <- 0
	}
}
