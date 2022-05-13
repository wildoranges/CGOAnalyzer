package tools

import (
	"go/parser"
	"go/token"
	"os"
)

//IsGo2CDir recursively check if there exists a file(ends with .go) under the given dir contains import "C"
//return true or false
func IsGo2CDir(pth string) bool {
	var flag bool
	flag = false
	fset := token.NewFileSet()
	msgReport(os.Stdout,showAll,"now checking %s\n", pth)
	f, _ := parser.ParseDir(fset, pth, nil, parser.ImportsOnly)
	/*if err != nil {
		msgReport(os.Stderr,true,"error occured while parsing %s\n", pth)
		msgReport(os.Stderr,true,err.Error())
	}*/
	for pkg, pkgast := range f {
		if flag {
			break
		}
		msgReport(os.Stdout,showAll,"now parsing:path:%s,pkg:%s\n", pth, pkg)
		for fn, srcfile := range pkgast.Files {
			if flag {
				break
			}
			for _, s := range srcfile.Imports {
				if flag {
					break
				}
				msgReport(os.Stdout,showAll,"pth:%s,pkg:%s,import:%s\n", pth, pkg, s.Path.Value)
				if s.Path.Value == "\"C\"" || s.Path.Value == "\"c\"" {
					flag = true
					msgReport(os.Stdout,showAll,"hit,go-c found in :%s,pkg:%s\n", fn, pkg)
					break
				}
			}
		}
	}
	return flag
}
