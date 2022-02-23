package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

type IAMPolicy struct {
	Version   string            `json:"Version"`
	Statement []PolicyStatement `json:"Statement"`
}
type PolicyStatement struct {
	Effect    string                       `json:"Effect"`
	Action    []string                     `json:"Action"`
	Resource  string                       `json:"Resource"`
	Condition map[string]map[string]string `json:"Condition"`
}

type IAMPolicyCondition map[string]IAMPolicyConditionKeyValue

type IAMPolicyConditionKeyValue map[string]interface{}

type AWSValue []string

func (v *AWSValue) UnmarshalJSON(input []byte) error {
	var raw interface{}
	json.Unmarshal(input, &raw)
	var elements []string
	switch item := raw.(type) {
	case string:
		elements = []string{item}
	case []interface{}:
		elements = make([]string, len(item))
		for i, it := range item {
			elements[i] = fmt.Sprintf("%s", it)
		}
	default:
		return fmt.Errorf("unsupported type %t in list", item)
	}
	*v = elements
	return nil
}

func main() {
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "./cloud_credential_tmpl.go", nil, 0)
	if err != nil {
		fmt.Println("Can't parse file", err)
	}
	jsFs, _ := ioutil.ReadFile("iam-policy.json")

	policy := IAMPolicy{}

	_ = json.Unmarshal([]byte(jsFs), &policy)

	exprs := make([]ast.Expr, 0, len(policy.Statement))
	for _, p := range policy.Statement {

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
			Value: buildStrings(p.Resource),
		}
		policyList[3] = &ast.KeyValueExpr{
			Key:   ast.NewIdent("Condition"),
			Value: buildKeyValueExpr(p.Condition),
		}

		compLit := ast.CompositeLit{Elts: policyList}
		exprs = append(exprs, &compLit)
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
								Key: ast.NewIdent("Version"),
								Value: &ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"somthing else\"",
								},
							},
							&ast.KeyValueExpr{
								Key: ast.NewIdent("Statement"),
								Value: &ast.CompositeLit{
									Type: &ast.ArrayType{
										Elt: ast.NewIdent("PolicyStatement"),
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

	fmt.Println("Modified AST:")
	printer.Fprint(os.Stdout, fs, file)
}

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
	default:
		panic("unsported type strings")
	}
}

func buildKeyValueExpr(input interface{}) ast.Expr {
	switch val := input.(type) {
	case map[string]map[string]string:
		exprs := make([]ast.Expr, 0, len(val))

		for k, v := range val {
			exprs = append(exprs, &ast.KeyValueExpr{
				Key:   buildStrings(k),
				Value: buildKeyValueExpr(v),
			})
		}

		return &ast.CompositeLit{
			Type: &ast.MapType{
				Key:   ast.NewIdent("string"),
				Value: ast.NewIdent("map[string]string"),
			},
			Elts: exprs,
		}
	case map[string]string:
		exprs := make([]ast.Expr, 0, len(val))
		for k, v := range val {
			exprs = append(exprs, &ast.KeyValueExpr{
				Key:   buildStrings(k),
				Value: buildStrings(v),
			})
		}
		return &ast.CompositeLit{
			Elts: exprs,
		}
	default:
		panic("unsported type keyvalue")
	}
}
