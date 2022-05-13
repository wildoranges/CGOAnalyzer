package tools

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

type FuncItem struct{
	label string
}

const GoCTypeTrans int = 0
const StdCLibFunc int = 1
const FuncUnknown int = 2
const StdCType int = 3

var FuncKind = map[int]string{
	GoCTypeTrans:"GoCTypeTrans",
	StdCLibFunc:"StdCLibFunc",
	FuncUnknown:"Unknown",
	StdCType:"StdCType",
}

func (i FuncItem)KindText(name string) string{
	if _, ok := cStdType[name];ok{
		return FuncKind[StdCType]
	}else if _, ok := TypeTrans[name];ok{
		return FuncKind[GoCTypeTrans]
	}
	return FuncKind[FuncUnknown]
}

func (i FuncItem)Label()string {
	return i.label
}

type void struct{

}

var member void

var TypeTrans = map[string]void{
	"GoString":member,
	"CString":member,
	"GoBytes":member,
}

var cStdType = map[string]void{//only part of the std c ctypes and macros
	"schar":         member,
	"uchar":         member,
	"ushort":        member,
	"uint":          member,
	"ulong":         member,
	"longlong":      member,
	"ulonglong":     member,
	"complexfloat":  member,
	"complexdouble": member,
	"float":		 member,
	"double":	     member,
	"short":		 member,
	"int":			 member,
	"char":			 member,
	"size_t":		 member,
	"intptr_t":		 member,
	"uintptr_t":     member,
	"u_char":        member,
	"u_short":		 member,
	"u_int":		 member,
	"u_long":		 member,
	"quad_t":		 member,
	"u_quad_t":		 member,
	"uint_t":		 member,
	"int8_t":		 member,
	"int16_t":		 member,
	"int32_t":		 member,
	"int64_t":		 member,
	"uint8_t":		 member,
	"uint16_t":		 member,
	"uint32_t":		 member,
	"uint64_t":      member,
	"bool":			 member,
}


//GoCFuncVisitor contains the number of each function while parsing a file/repo
type GoCFuncVisitor struct {
	curFile map[string]int
	fileName string
}

func isCStdType(name string) bool {
	_, ok := cStdType[name]
	return ok
}

//Visit counts the number of each function while walking traverse an AST
func (v GoCFuncVisitor) Visit(n ast.Node) ast.Visitor {
	if call, ok := n.(*ast.CallExpr); ok {
		if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
			funcName := fun.Sel.Name
			if pkgname, ok := fun.X.(*ast.Ident); ok {
				msgReport(os.Stdout,showAll,"now checking func:%s\n", funcName)
				if pkgname.Name == "C" {
					if !isCStdType(funcName) {
						msgReport(os.Stdout,showAll,"Go2CFunc:%s\n", funcName)
						if value, ok := v.curFile[funcName]; ok {
							v.curFile[funcName] = value + 1
						} else {
							v.curFile[funcName] = 1
						}
					}
				}
			}
		}
	}
	return v
}

//CntGo2CFunc counts the go-c function in each repo,return two maps.
//the first map contains the overall result.its key is C function name and its value is the number of the function
//the second map contains the result in each repo,the key is repo's name and the value is a map containing
//the repo's local result
func CntGo2CFunc(repos []string) (ItemInfo) {
	var Info ItemInfo
	Info.i = FuncItem{"function"}
	Info.total = make(map[string]int)
	Info.local = make(map[string]map[string]int)
	Info.pathInfo = make(map[string]StringSet)
	for _, repo := range repos {
		msgReport(os.Stdout,showAll,"now parsing repo:%s\n", repo)
		repoM := make(map[string]int)
		alldirs, err := GetAllDirs(repo)
		if err != nil {
			msgReport(os.Stderr,true,err.Error()+"\n")
			continue
		}
		for _, repoDir := range alldirs {
			set := token.NewFileSet()
			f, err := parser.ParseDir(set, repoDir, nil, 0)
			if err != nil {
				msgReport(os.Stderr,true,err.Error()+"\n")
			}
			for pkg, pkgast := range f {
				msgReport(os.Stdout,showAll,"now parsing pkg:%s\n", pkg)
				for filename, srcfile := range pkgast.Files {
					curfile := make(map[string]int)
					m := GoCFuncVisitor{curFile: curfile,fileName: filename}
					msgReport(os.Stdout,showAll,"now parsing file:%s\n", filename)
					ast.Walk(m, srcfile)
					for k, v := range m.curFile{
						if Info.pathInfo[k] == nil {
							Info.pathInfo[k] = make(StringSet)
						}
						pathRecord(dumpAll,Info,k,filename)
						if value, ok := repoM[k]; ok{
							repoM[k] = value + v
						}else{
							repoM[k] = v
						}
					} 
				}
			}
		}
		Info.local[repo] = repoM
		for k, v := range repoM {
			if value, ok := Info.total[k]; ok {
				Info.total[k] = value + v
			} else {
				Info.total[k] = v
			}
		}
	}
	return Info
}
