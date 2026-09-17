package apicheck

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

type wrapperMethod struct {
	Recv  string // e.g. queryWrapper
	Name  string
	File  string
	Start int
	End   int
	Key   string // Recv.Name
}

func TestWrapperMethodCoverage(t *testing.T) {
	profile := os.Getenv("FSMOCK_COVERPROFILE")
	if profile == "" {
		t.Skip("FSMOCK_COVERPROFILE not set; accuracy gate sets this")
	}

	root, err := moduleRoot()
	if err != nil {
		t.Fatal(err)
	}

	methods, err := listWrapperMethods(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(methods) == 0 {
		t.Fatal("no *xxxWrapper methods found — listWrapperMethods broken?")
	}

	covered, err := parseCoverProfile(profile, root)
	if err != nil {
		t.Fatalf("parse cover profile %s: %v\n  Fix: run accuracy gate so cover.out is produced", profile, err)
	}

	usedWaivers := map[string]bool{}
	var missing []string

	for _, m := range methods {
		ok := methodCovered(covered, m)
		if ok {
			continue
		}
		if w, ok := waivers[m.Key]; ok {
			if w.Reason == "" || w.Issue == "" {
				t.Errorf("waiver %q missing Reason or Issue", m.Key)
			}
			usedWaivers[m.Key] = true
			t.Logf("waive %s: %s (%s)", m.Key, w.Reason, w.Issue)
			continue
		}
		missing = append(missing, fmt.Sprintf("%s (%s:%d)", m.Key, filepath.Base(m.File), m.Start))
	}

	if len(missing) > 0 {
		t.Errorf("%d wrapper method(s) lack coverage and have no waiver:\n  %s\n  Fix: add an integration test in fstest/ that exercises the method, or add waivers[%q] with Reason+Issue in internal/apicheck/waivers.go",
			len(missing), strings.Join(missing, "\n  "), "recv.Method")
	}

	for key, w := range waivers {
		if usedWaivers[key] {
			continue
		}
		// Waiver unused: either method is now covered, or method was removed.
		found := false
		var m wrapperMethod
		for i := range methods {
			if methods[i].Key == key {
				found = true
				m = methods[i]
				break
			}
		}
		if !found {
			t.Errorf("stale waiver %q — no such wrapper method; delete it from waivers.go", key)
			continue
		}
		if methodCovered(covered, m) {
			t.Errorf("stale waiver %q — method is now covered; delete it from waivers.go (%s)", key, w.Reason)
		}
	}
}

func moduleRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", wd)
		}
		dir = parent
	}
}

func listWrapperMethods(root string) ([]wrapperMethod, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	var out []wrapperMethod
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		path := filepath.Join(root, name)
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil, err
		}
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Recv == nil || fd.Name == nil || !fd.Name.IsExported() {
				continue
			}
			recv := receiverIdent(fd.Recv)
			if recv == "" || !strings.HasSuffix(recv, "Wrapper") {
				continue
			}
			start := fset.Position(fd.Pos()).Line
			end := fset.Position(fd.End()).Line
			out = append(out, wrapperMethod{
				Recv:  recv,
				Name:  fd.Name.Name,
				File:  path,
				Start: start,
				End:   end,
				Key:   recv + "." + fd.Name.Name,
			})
		}
	}
	return out, nil
}

func receiverIdent(recv *ast.FieldList) string {
	if recv == nil || len(recv.List) == 0 {
		return ""
	}
	expr := recv.List[0].Type
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}
	id, ok := expr.(*ast.Ident)
	if !ok {
		return ""
	}
	return id.Name
}

type coverBlock struct {
	File  string
	Start int
	End   int
	Count int
}

func parseCoverProfile(path, root string) ([]coverBlock, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var blocks []coverBlock
	sc := bufio.NewScanner(f)
	// mode: set
	if !sc.Scan() {
		return nil, fmt.Errorf("empty cover profile")
	}
	for sc.Scan() {
		line := sc.Text()
		// format: path:startLine.startCol,endLine.endCol numStmt count
		colon := strings.LastIndex(line, ":")
		if colon < 0 {
			continue
		}
		filePart := line[:colon]
		rest := line[colon+1:]
		fields := strings.Fields(rest)
		if len(fields) != 3 {
			continue
		}
		rangePart := fields[0]
		count, err := strconv.Atoi(fields[2])
		if err != nil {
			continue
		}
		comma := strings.Index(rangePart, ",")
		if comma < 0 {
			continue
		}
		startS := strings.Split(rangePart[:comma], ".")[0]
		endS := strings.Split(rangePart[comma+1:], ".")[0]
		start, _ := strconv.Atoi(startS)
		end, _ := strconv.Atoi(endS)

		filePath := filePart
		if !filepath.IsAbs(filePath) {
			// cover profiles use module path prefix
			const mod = "github.com/akmalsyrf/go-firestore-mock/v2/"
			if strings.HasPrefix(filePath, mod) {
				filePath = filepath.Join(root, strings.TrimPrefix(filePath, mod))
			} else if strings.HasPrefix(filePath, root) {
				// already absolute-ish
			} else {
				filePath = filepath.Join(root, filepath.Base(filePath))
			}
		}
		blocks = append(blocks, coverBlock{File: filePath, Start: start, End: end, Count: count})
	}
	return blocks, sc.Err()
}

func methodCovered(blocks []coverBlock, m wrapperMethod) bool {
	base := filepath.Base(m.File)
	for _, b := range blocks {
		if b.Count <= 0 {
			continue
		}
		if filepath.Base(b.File) != base && b.File != m.File {
			continue
		}
		// Overlap with method body
		if b.End < m.Start || b.Start > m.End {
			continue
		}
		return true
	}
	return false
}
