package apicheck

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// sdkModuleDir returns the filesystem path of the pinned firestore module.
func sdkModuleDir() (string, error) {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "cloud.google.com/go/firestore")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("go list firestore module: %w", err)
	}
	dir := strings.TrimSpace(string(out))
	if dir == "" {
		return "", fmt.Errorf("empty firestore module dir")
	}
	return dir, nil
}

// discoverSDKTypesWithMethods returns exported type names in the firestore package
// that have at least one exported method (value or pointer receiver).
func discoverSDKTypesWithMethods(sdkDir string) (map[string][]string, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, sdkDir, func(info os.FileInfo) bool {
		name := info.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			return false
		}
		// Skip subdirectories handled separately; ParseDir is non-recursive for files in sdkDir only.
		return true
	}, 0)
	if err != nil {
		return nil, err
	}
	pkg, ok := pkgs["firestore"]
	if !ok {
		return nil, fmt.Errorf("firestore package not found in %s", sdkDir)
	}

	methods := map[string][]string{}
	seen := map[string]map[string]bool{}

	for _, f := range pkg.Files {
		// Only files directly in sdkDir (ignore nested if any slipped in).
		if filepath.Dir(fset.File(f.Pos()).Name()) != sdkDir {
			continue
		}
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Recv == nil || fd.Name == nil || !fd.Name.IsExported() {
				continue
			}
			if len(fd.Recv.List) == 0 {
				continue
			}
			typeName := recvTypeName(fd.Recv.List[0].Type)
			if typeName == "" || !ast.IsExported(typeName) {
				continue
			}
			if seen[typeName] == nil {
				seen[typeName] = map[string]bool{}
			}
			if seen[typeName][fd.Name.Name] {
				continue
			}
			seen[typeName][fd.Name.Name] = true
			methods[typeName] = append(methods[typeName], fd.Name.Name)
		}
	}
	return methods, nil
}

func recvTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return recvTypeName(t.X)
	default:
		return ""
	}
}
