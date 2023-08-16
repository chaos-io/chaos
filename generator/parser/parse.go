package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"

	"github.com/chaos-io/chaos/core/strcase"
	"github.com/chaos-io/chaos/logs"
)

// FileParser is the parser used by kit to parse go files.
type FileParser struct{}

// Parse will parse the go source.
func (fp *FileParser) Parse(src []byte) (*File, error) {
	f := NewFile()

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	pf, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	f.Package = pf.Name.Name

	for _, v := range pf.Decls {
		if dec, ok := v.(*ast.FuncDecl); ok {
			var st []NamedTypeValue
			var pr []NamedTypeValue
			var rs []NamedTypeValue

			if dec.Recv != nil {
				st = fp.parseFieldListAsNamedTypes(dec.Recv)
			}

			if dec.Type != nil {
				pr = fp.parseFieldListAsNamedTypes(dec.Type.Params)
				rs = fp.parseFieldListAsNamedTypes(dec.Type.Results)
			}

			bd := ""
			if dec.Body != nil {
				fst := token.NewFileSet()
				bt := bytes.NewBufferString("")
				err := format.Node(bt, fst, dec.Body)
				bd = bt.String()[1 : len(bt.String())-1]
				if err != nil {
					logs.Fatal(err)
				}
			}

			str := NamedTypeValue{}
			if len(st) > 0 {
				str = st[0]
			}
			fc := NewMethod(dec.Name.String(), str, bd, pr, rs)
			f.Methods = append(f.Methods, fc)
		}

		if dec, ok := v.(*ast.GenDecl); ok {
			switch dec.Tok {
			case token.IMPORT:
				f.Imports = fp.parseImports(dec.Specs)
			case token.CONST:
				f.Constants = append(f.Constants, fp.parseConstants(dec.Specs)...)
			case token.VAR:
				f.Vars = append(f.Vars, fp.parseVars(dec.Specs)...)
			case token.TYPE:
				fp.parseType(dec.Specs, &f)
			default:
				logs.Info("Skipping unknown Token Type")
			}
		}
	}

	return &f, nil
}
func (fp *FileParser) parseType(ds []ast.Spec, f *File) {
	for _, sp := range ds {
		tsp, ok := sp.(*ast.TypeSpec)
		if !ok {
			logs.Debug("Type spec is not TypeSpec type, odd, skipping")
			continue
		}
		switch tsp.Type.(type) {
		case *ast.InterfaceType:
			ift := tsp.Type.(*ast.InterfaceType)
			mth := fp.parseFieldListAsMethods(ift.Methods)
			intr := NewInterface(tsp.Name.Name, mth)
			intr.Methods = mth
			f.Interfaces = append(f.Interfaces, intr)
		case *ast.StructType:
			st := tsp.Type.(*ast.StructType)
			str := NewStruct(tsp.Name.Name, fp.parseFieldListAsNamedTypes(st.Fields))
			f.Structures = append(f.Structures, str)
		case *ast.FuncType:
			st := tsp.Type.(*ast.FuncType)
			f.FuncType = FuncType{
				Name:       tsp.Name.Name,
				Parameters: fp.parseFieldListAsNamedTypes(st.Params),
				Results:    fp.parseFieldListAsNamedTypes(st.Results),
			}
		default:
			logs.Info("Skipping unknown type")
		}
	}
}
func (fp *FileParser) parseImports(ds []ast.Spec) []NamedTypeValue {
	var imports []NamedTypeValue
	for _, sp := range ds {
		isp, ok := sp.(*ast.ImportSpec)
		if !ok {
			logs.Debug("Import spec is not ImportSpec type, odd, skipping")
			continue
		}
		ip := NewNameType("", "")
		if isp.Name != nil {
			ip.Name = isp.Name.Name
		}
		if isp.Path != nil {
			ip.Type = isp.Path.Value
		}
		imports = append(imports, ip)
	}
	return imports
}
func (fp *FileParser) parseVars(ds []ast.Spec) []NamedTypeValue {
	var vars []NamedTypeValue
	for _, sp := range ds {
		vsp, ok := sp.(*ast.ValueSpec)
		if !ok {
			logs.Debug("Var spec is not ValueSpec type, odd, skipping")
			continue
		}
		tp, ok := vsp.Type.(*ast.Ident)
		if len(vsp.Values) > 0 {
			fst := token.NewFileSet()
			bt := bytes.NewBufferString("")
			err := format.Node(bt, fst, vsp.Values[0])
			bd := bt.String()
			if err != nil {
				logs.Fatal(err)
			}
			if !ok {
				vars = append(vars, NewNameTypeValue(vsp.Names[0].Name, "", bd))
				continue
			}
			vars = append(vars, NewNameTypeValue(tp.Name, vsp.Names[0].Name, bd))
		} else {
			if !ok {
				vars = append(vars, NewNameType(vsp.Names[0].Name, ""))
				continue
			}
			vars = append(vars, NewNameType(tp.Name, vsp.Names[0].Name))
		}
	}
	return vars
}
func (fp *FileParser) parseConstants(ds []ast.Spec) []NamedTypeValue {
	var constants []NamedTypeValue
	for _, sp := range ds {
		vsp, ok := sp.(*ast.ValueSpec)
		if !ok {
			logs.Debug("Constant spec is not ValueSpec type, odd, skipping")
			continue
		}
		fst := token.NewFileSet()
		bt := bytes.NewBufferString("")
		err := format.Node(bt, fst, vsp.Values[0])
		bd := bt.String()
		if err != nil {
			logs.Fatal(err)
		}
		tp, ok := vsp.Type.(*ast.Ident)
		if !ok {
			constants = append(constants, NewNameTypeValue(vsp.Names[0].Name, "", bd))
			continue
		}
		constants = append(constants, NewNameTypeValue(tp.Name, vsp.Names[0].Name, bd))
	}
	return constants
}
func (fp *FileParser) parseFieldListAsNamedTypes(list *ast.FieldList) []NamedTypeValue {
	var ntv []NamedTypeValue
	if list != nil {
		for i, p := range list.List {
			typ := fp.getTypeFromExp(p.Type)
			logs.Debug(fmt.Sprintf("Type %s", typ))

			// Potentially N names
			var names []string
			for _, ident := range p.Names {
				names = append(names, ident.Name)
			}
			if len(names) == 0 {
				// Anonymous named type, give it a default name
				if strings.HasPrefix(typ, "[]") {
					names = append(names, strcase.ToLowerCamel(typ[2:3]+fmt.Sprintf("%d", i)))
				} else if strings.HasPrefix(typ, "*") {
					names = append(names, strcase.ToLowerCamel(typ[1:2]+fmt.Sprintf("%d", i)))
				} else {
					names = append(names, strcase.ToLowerCamel(typ[:1]+fmt.Sprintf("%d", i)))
				}
			}
			for _, name := range names {
				namedType := NewNameType(name, typ)
				logs.Debug(fmt.Sprintf("NamedType %+v", namedType))
				ntv = append(ntv, namedType)
			}
		}
	}
	return ntv
}
func (fp *FileParser) getTypeFromExp(e ast.Expr) string {
	tp := ""
	switch k := e.(type) {
	case *ast.Ident:
		tp = k.Name
	case *ast.SelectorExpr:
		logs.Debug("Type Selector, i.e. a third-party type")
		selectorIdent := fp.getTypeFromExp(k.X)
		tp = fmt.Sprintf(fmt.Sprintf("%s.%s", selectorIdent, k.Sel.Name))
	case *ast.StarExpr:
		starIndent := fp.getTypeFromExp(k.X)
		tp = "*" + starIndent
	case *ast.ArrayType:
		arrIndent := fp.getTypeFromExp(k.Elt)
		tp = "[]" + arrIndent
	case *ast.MapType:
		key := fp.getTypeFromExp(k.Key)
		value := fp.getTypeFromExp(k.Value)
		tp = "map[" + key + "]" + value
	case *ast.InterfaceType:
		tp = "interface{}"
	case *ast.Ellipsis:
		t := fp.getTypeFromExp(k.Elt)
		tp = "..." + t
	default:
		logs.Info("Type Expresion not supported")
		return ""
	}
	return tp
}
func (fp *FileParser) parseFieldListAsMethods(list *ast.FieldList) []Method {
	var mth []Method
	if list != nil {
		for _, p := range list.List {
			switch t := p.Type.(type) {
			case *ast.FuncType:
				m := Method{
					Name: p.Names[0].Name,
				}
				m.Parameters = fp.parseFieldListAsNamedTypes(t.Params)
				m.Results = fp.parseFieldListAsNamedTypes(t.Results)
				mth = append(mth, m)
			default:
				logs.Info("Skipping unknown type")
			}
		}
	}
	return mth
}

// NewFileParser returns a new parser.
func NewFileParser() *FileParser {
	return &FileParser{}
}
