package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// Color escape codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

var (
	maxForRangeLevel = 3
	targetDirectory  = ""
	ignoreDirs       = make([]string, 0)
	exitMutex        sync.Mutex
)

func main() {

	args := os.Args[1:]

	if len(args) > 0 && args[0] != "" {
		maxForRangeLevel, _ = strconv.Atoi(args[0])
	}

	if len(args) > 1 {
		targetDirectory = args[1]
	}
	if len(args) > 2 {
		if err := json.Unmarshal([]byte(args[2]), &ignoreDirs); err != nil {
			fmt.Printf("json.Unmarshal err: %v, args[2]: %v", err, args[2])
		}
	}
	fmt.Printf("%s action params: %s\n", Yellow, Reset)
	fmt.Printf("%s \t maxForRangeLevel: %d, targetDirectory: %s %s\n", Yellow, maxForRangeLevel, targetDirectory, Reset)
	fmt.Printf("%s \t ignoreDirs: %v  %s\n", Yellow, ignoreDirs, Reset)
	if passMethods, failedMethods, err := ScanBizPath(); err != nil {
		fmt.Printf("%s err: %v %s\n", Yellow, err, Reset)
		fmt.Printf("%s passMethods: %s\n", Blue, Reset)
		for _, item := range passMethods {
			fmt.Printf("%s \t %s %s\n", Blue, item, Reset)
		}
		fmt.Printf("%s failedMethods: %s\n", Red, Reset)
		for _, item := range failedMethods {
			fmt.Printf("%s \t %s %s\n", Red, item, Reset)
		}
		defer func() {
			exitMutex.Lock()
			defer exitMutex.Unlock()

			// Exit with a non-zero code to indicate failure
			os.Exit(1)
		}()
	} else {
		fmt.Printf("%s passMethods: %s\n", Green, Reset)
		for _, item := range passMethods {
			fmt.Printf("%s \t %s %s\n", Green, item, Reset)
		}
		fmt.Printf("%s success... %s\n", Green, Reset)
	}
}

func ScanBizPath() ([]string, []string, error) {
	passMethodList := make([]string, 0)
	failedMethodList := make([]string, 0)
	// 遍历指定目录下的所有 Go 文件
	err := filepath.Walk(targetDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("%s Error walking directory: %s %s", Red, err, Reset)
			return err
		}
		if info.IsDir() && isPathInIgnoreDir(info.Name()) {
			// ignore dir
			fmt.Printf("%s ignore dir: %v, path: %s %s\n", Yellow, info.Name(), path, Reset)
			return nil
		}
		if strings.HasSuffix(info.Name(), ".go") {
			functionStats, err := checkFile(path)
			if err != nil {
				failedMethodList = append(failedMethodList, fmt.Sprintf("Error in file %s: %v", path, err))
			} else {
				for functionName, loopDepth := range functionStats {
					if loopDepth > maxForRangeLevel {
						failedMethodList = append(failedMethodList, fmt.Sprintf("Function %s in file %s, loop depth: %d", functionName, path, loopDepth))
					} else {
						passMethodList = append(passMethodList, fmt.Sprintf("Function %s in file %s, loop depth: %d", functionName, path, loopDepth))
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, nil, err
	}
	if len(failedMethodList) > 0 {
		return passMethodList, failedMethodList, errors.New("methods did not pass the test")
	} else {
		return passMethodList, failedMethodList, nil
	}

}

func checkFile(filePath string) (map[string]int, error) {
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, err
	}

	functionStats := make(map[string]int)

	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			// Check functions/methods for nested loops
			loopStats := countLoopDepth(fn.Body)
			functionStats[fn.Name.Name] = loopStats
		}
		return true
	})

	return functionStats, nil
}

func countLoopDepth(body *ast.BlockStmt) int {
	var maxDepth int
	if body == nil {
		return 0
	}
	for _, stmt := range body.List {
		ast.Inspect(stmt, func(n ast.Node) bool {

			var depth int
			switch n.(type) {
			case *ast.IfStmt:
				ifAst, _ := n.(*ast.IfStmt)
				_ = countLoopDepth(ifAst.Body)
			case *ast.BlockStmt:
				block, _ := n.(*ast.BlockStmt)
				_ = countLoopDepth(block)
			case *ast.ForStmt:
				st, _ := n.(*ast.ForStmt)
				depth = countLoopDepth(st.Body)
				if depth+1 > maxDepth {
					maxDepth = depth + 1
				}
			case *ast.RangeStmt:
				st, _ := n.(*ast.RangeStmt)
				depth = countLoopDepth(st.Body)
				if depth+1 > maxDepth {
					maxDepth = depth + 1
				}
			default:
			}
			return true
		})
	}

	return maxDepth
}

func isPathInIgnoreDir(path string) bool {
	ignoreMaps := make(map[string]bool, 0)
	for _, item := range ignoreDirs {
		ignoreMaps[item] = true
	}
	if _, existed := ignoreMaps[path]; existed {
		return true
	}
	return false
}
