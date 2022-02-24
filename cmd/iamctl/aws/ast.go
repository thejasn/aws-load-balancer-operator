package aws

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

func buildStrings(input interface{}) ast.Expr {
	switch val := input.(type) {
	case string:
		return &ast.BasicLit{
			Kind:  token.STRING,
			Value: "\"" + val + "\"",
		}
	case []string:
		ret := make([]ast.Expr, 0, len(val))
		for _, s := range val {
			ret = append(ret, buildStrings(s))
		}
		return &ast.CompositeLit{
			Type: ast.NewIdent("[]string"),
			Elts: ret,
		}
	case AWSValue:
		ret := make([]ast.Expr, 0, len(val))
		for _, s := range val {
			ret = append(ret, buildStrings(s))
		}
		return &ast.CompositeLit{
			Type: ast.NewIdent("[]string"),
			Elts: ret,
		}
	default:
		panic("unsported type for string expr")
	}
}

func buildKeyValueExpr(input interface{}) ast.Expr {
	switch val := input.(type) {
	case *iamPolicyCondition:
		if val == nil {
			return ast.NewIdent("nil")
		}
		exprs := make([]ast.Expr, 0, 1)

		for k, v := range *val {
			exprs = append(exprs, &ast.KeyValueExpr{
				Key:   buildStrings(k),
				Value: buildKeyValueExpr(v),
			})
		}
		return &ast.CompositeLit{
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("cco"),
				Sel: ast.NewIdent("IAMPolicyCondition"),
			},
			Elts: exprs,
		}
	case iamPolicyConditionKeyValue:
		exprs := make([]ast.Expr, 0, 1)
		for k, v := range val {

			exprs = append(exprs, &ast.KeyValueExpr{
				Key:   buildStrings(k),
				Value: buildStrings(v),
			})
		}
		return &ast.CompositeLit{
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("cco"),
				Sel: ast.NewIdent("IAMPolicyConditionKeyValue"),
			},
			Elts: exprs,
		}
	default:
		panic("unsported type for key/val expr")
	}
}

func GenerateIAMPolicy(input, output, pkg string) {

	in := strings.ReplaceAll(filetemplate, "main", pkg)

	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "", in, 0)
	if err != nil {
		panic(fmt.Errorf("failed to parse template %v", err))
	}

	jsFs, err := ioutil.ReadFile(input)
	if err != nil {
		panic(fmt.Errorf("failed to read input file %v", err))
	}

	policy := iamPolicy{}

	err = json.Unmarshal([]byte(jsFs), &policy)
	if err != nil {
		panic(fmt.Errorf("failed to parse policy JSON %v", err))
	}

	exprs := make([]ast.Expr, 0, len(policy.Statement))
	for _, p := range policy.Statement {
		if len(p.Resource) > 1 {
			for _, r := range p.Resource {

				policyList := make([]ast.Expr, 4)

				policyList[0] = &ast.KeyValueExpr{
					Key:   ast.NewIdent("Effect"),
					Value: buildStrings(p.Effect),
				}
				policyList[1] = &ast.KeyValueExpr{
					Key:   ast.NewIdent("Action"),
					Value: buildStrings(p.Action),
				}

				policyList[2] = &ast.KeyValueExpr{
					Key:   ast.NewIdent("Resource"),
					Value: buildStrings(r),
				}

				policyList[3] = &ast.KeyValueExpr{
					Key:   ast.NewIdent("PolicyCondition"),
					Value: buildKeyValueExpr(p.Condition),
				}
				exprs = append(exprs, &ast.CompositeLit{Elts: policyList})
			}
			continue
		}

		policyList := make([]ast.Expr, 4)

		policyList[0] = &ast.KeyValueExpr{
			Key:   ast.NewIdent("Effect"),
			Value: buildStrings(p.Effect),
		}
		policyList[1] = &ast.KeyValueExpr{
			Key:   ast.NewIdent("Action"),
			Value: buildStrings(p.Action),
		}

		policyList[2] = &ast.KeyValueExpr{
			Key:   ast.NewIdent("Resource"),
			Value: buildStrings(p.Resource[0]),
		}

		policyList[3] = &ast.KeyValueExpr{
			Key:   ast.NewIdent("PolicyCondition"),
			Value: buildKeyValueExpr(p.Condition),
		}

		exprs = append(exprs, &ast.CompositeLit{Elts: policyList})
	}

	astutil.Apply(file, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch x := n.(type) {
		case *ast.ReturnStmt:
			c.Replace(&ast.ReturnStmt{
				Return: x.Pos(),
				Results: []ast.Expr{
					&ast.CompositeLit{
						Type: ast.NewIdent("IAMPolicy"),
						Elts: []ast.Expr{
							&ast.KeyValueExpr{
								Key:   ast.NewIdent("Version"),
								Value: buildStrings(policy.Version),
							},
							&ast.KeyValueExpr{
								Key: ast.NewIdent("Statement"),
								Value: &ast.CompositeLit{
									Type: &ast.ArrayType{
										Elt: &ast.SelectorExpr{
											X:   ast.NewIdent("cco"),
											Sel: ast.NewIdent("StatementEntry"),
										},
									},
									Elts: exprs,
								},
							},
						},
					},
				},
			})
		}

		return true
	})

	opFs, err := os.Create(output)
	if err != nil {
		panic(err)
	}

	cfg := &printer.Config{Mode: printer.UseSpaces, Tabwidth: 4, Indent: 0}
	cfg.Fprint(opFs, fs, file)
}
