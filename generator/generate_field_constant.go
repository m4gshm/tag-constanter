package generator

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"github.com/m4gshm/fieldr/logger"
	"github.com/m4gshm/fieldr/struc"
	"github.com/pkg/errors"
)

type stringer struct {
	val      string
	callback func()
}

func (c *stringer) String() string {
	c.callback()
	return c.val
}

var _ fmt.Stringer = (*stringer)(nil)

func (g *Generator) GenerateFieldConstants(model *struc.Model, fieldType string, fieldNames []struc.FieldName, export, snake, wrapType bool) error {
	typeName := model.TypeName
	g.addConstDelim()
	for _, fieldName := range fieldNames {
		constName := GetFieldConstName(typeName, fieldName, export, snake)
		constVal := g.GetConstValue(fieldType, fieldName, wrapType)
		if err := g.addConst(constName, constVal); err != nil {
			return err
		}
	}
	return nil
}

type constResult struct{ name, field, value string }

func (g *Generator) GenerateFieldConstant(
	model *struc.Model, value, name, typ string, export, snake, nolint, compact, usePrivate, refAccessor, valAccessor bool,
) error {
	valueTmpl := wrapTemplate(value)
	nameTmpl := wrapTemplate(name)

	wrapType := len(typ) > 0
	if !wrapType {
		typ = BaseConstType
	} else if err := g.AddType(typ, BaseConstType); err != nil {
		return err
	}

	usedTags := []string{}
	usedTagsSet := map[string]struct{}{}

	constants := make([]constResult, 0)
	for _, f := range model.FieldNames {
		var (
			inExecute bool
			fieldName = f
			tags      = map[string]interface{}{}
		)

		if isFieldExcluded(fieldName, usePrivate) {
			continue
		}

		if tagVals := model.FieldsTagValue[fieldName]; tagVals != nil {
			for k, v := range tagVals {
				tag := k
				tags[tag] = &stringer{val: v, callback: func() {
					if !inExecute {
						return
					}
					if _, ok := usedTagsSet[tag]; !ok {
						usedTagsSet[tag] = struct{}{}
						usedTags = append(usedTags, tag)
						logger.Debugf("use tag '%s'", tag)
					}
				}}
			}
		}

		parse := func(name string, data interface{}, funcs template.FuncMap, tmplVal string) (string, error) {
			logger.Debugf("parse template for \"%s\" %s\n", name, tmplVal)
			tmpl, err := template.New(value).Option("missingkey=zero").Funcs(funcs).Parse(tmplVal)
			if err != nil {
				return "", fmt.Errorf("parse: of '%s', template %s: %w", name, tmplVal, err)
			}

			buf := bytes.Buffer{}
			logger.Debugf("template context %+v\n", tags)
			inExecute = true
			if err = tmpl.Execute(&buf, data); err != nil {
				inExecute = false
				return "", fmt.Errorf("compile: of '%s': field '%s', template %s: %w", name, fieldName, tmplVal, err)
			}
			inExecute = false
			cmpVal := buf.String()
			logger.Debugf("parse result: of '%s'; %s\n", name, cmpVal)
			return cmpVal, nil
		}

		funcs := addCommonFuncs(template.FuncMap{
			"struct": func() map[string]interface{} { return map[string]interface{}{"name": model.TypeName} },
			"name":   func() string { return fieldName },
			"field":  func() map[string]interface{} { return map[string]interface{}{"name": fieldName} },
			"tag":    func() map[string]interface{} { return tags },
		})

		if val, err := parse(fieldName+" const val", tags, funcs, valueTmpl); err != nil {
			return err
		} else if len(nameTmpl) > 0 {
			constName, err := parse(fieldName+" const name", tags, funcs, nameTmpl)
			if err != nil {
				return err
			}
			constants = append(constants, constResult{field: fieldName, name: strings.ReplaceAll(constName, ".", ""), value: val})
		} else {
			constants = append(constants, constResult{field: fieldName, value: val})
		}
	}

	for _, c := range constants {
		constName := c.name
		if len(constName) == 0 {
			constName = g.getTagTemplateConstName(model.TypeName, c.field, usedTags, export, snake)
			logger.Debugf("apply auto constant name '%s'", constName)
		} else {
			logger.Debugf("template generated constant name '%s'", constName)
		}
		if len(c.value) != 0 {
			if err := g.addConst(constName, g.GetConstValue(typ, c.value, wrapType)); err != nil {
				return err
			}
		} else {
			logger.Infof("constant without value: '%s'", constName)
		}
	}
	g.addConstDelim()

	if wrapType {
		exportFunc := export
		if funcBody, funcName, err := generateAggregateFunc(typ, constants, exportFunc, compact, nolint); err != nil {
			return err
		} else if err := g.AddFunc(funcName, funcBody); err != nil {
			return err
		}
		g.addFunсDelim()

		if funcBody, funcName, err := g.aenerateConstFieldFunc(typ, constants, exportFunc, nolint); err != nil {
			return err
		} else if err := g.AddFunc(funcName, funcBody); err != nil {
			return err
		}
		g.addFunсDelim()

		if refAccessor || valAccessor {
			structPackage, err := g.StructPackage(model)
			if err != nil {
				return err
			}
			if valAccessor {
				if funcBody, funcName, err := g.GenerateConstValueFunc(model, structPackage, typ, constants, exportFunc, nolint, false); err != nil {
					return err
				} else if err := g.AddFunc(funcName, funcBody); err != nil {
					return err
				}
			}
			if refAccessor {
				if funcBody, funcName, err := g.GenerateConstValueFunc(model, structPackage, typ, constants, exportFunc, nolint, true); err != nil {
					return err
				} else if err := g.AddFunc(funcName, funcBody); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func addCommonFuncs(funcs template.FuncMap) template.FuncMap {
	toString := func(val interface{}) string {
		str := ""
		switch vt := val.(type) {
		case string:
			str = vt
		case fmt.Stringer:
			str = vt.String()
		case fmt.GoStringer:
			str = vt.GoString()
		default:
			str = fmt.Sprint(val)
			logger.Debugf("toString: val '%v', result '%s'", val, str)
		}
		return str
	}
	toStrings := func(vals []interface{}) []string {
		results := make([]string, len(vals))
		for i, val := range vals {
			results[i] = toString(val)
		}
		return results
	}
	rexp := func(expr interface{}, val interface{}) (string, error) {
		sexpr := toString(expr)
		str := toString(val)
		if len(sexpr) == 0 {
			return "", errors.New("empty regexp: val '" + str + "'")
		}
		r, err := regexp.Compile(sexpr)
		if err != nil {
			return "", err
		}
		submatches := r.FindStringSubmatch(str)
		names := r.SubexpNames()
		if len(names) <= len(submatches) {
			for i, groupName := range names {
				if groupName == "v" {
					submatch := submatches[i]
					if len(submatch) == 0 {
						return submatch, nil
					}
				}
			}
		}
		if len(submatches) > 0 {
			s := submatches[len(submatches)-1]
			return s, nil
		}
		return "", nil
	}

	snakeFunc := func(val interface{}) string {
		sval := toString(val)
		if len(sval) == 0 {
			return ""
		}
		last := len(sval) - 1
		symbols := []rune(sval)
		result := make([]rune, 0)
		for i := 0; i < len(symbols); i++ {
			cur := symbols[i]
			result = append(result, cur)
			if i < last {
				next := symbols[i+1]
				if unicode.IsLower(cur) && unicode.IsUpper(next) {
					result = append(result, '_', unicode.ToLower(next))
					i++
				}
			}
		}
		return string(result)
	}

	toUpper := func(val interface{}) string {
		return strings.ToUpper(toString(val))
	}

	toLower := func(val interface{}) string {
		return strings.ToLower(toString(val))
	}

	join := func(vals ...interface{}) string {
		result := strings.Join(toStrings(vals), "")
		return result
	}

	strOr := func(vals ...string) string {
		if len(vals) == 0 {
			return ""
		}

		for _, val := range vals {
			if len(val) > 0 {
				return val
			}
		}
		return vals[len(vals)-1]
	}
	f := template.FuncMap{
		"OR":          strOr,
		"conc":        join,
		"concatenate": join,
		"join":        join,
		"regexp":      rexp,
		"rexp":        rexp,
		"snake":       snakeFunc,
		"toUpper":     toUpper,
		"toLower":     toLower,
		"up":          toUpper,
		"low":         toLower,
	}
	for k, v := range f {
		funcs[k] = v
	}
	return funcs
}

func wrapTemplate(text string) string {
	if len(text) == 0 {
		return text
	}
	if !strings.Contains(text, "{{") {
		text = "{{" + strings.ReplaceAll(text, "\\", "\\\\") + "}}"
		logger.Debugf("constant template transformed to '%s'", text)
	}
	return text
}

func generateAggregateFunc(typ string, constants []constResult, export, compact, nolint bool) (string, string, error) {
	var (
		funcName  = goName(typ+"s", export)
		arrayType = "[]" + typ
	)

	arrayBody := arrayType + "{"

	compact = compact || len(constants) <= oneLineSize
	if !compact {
		arrayBody += "\n"
	}

	i := 0
	for _, constant := range constants {
		if compact && i > 0 {
			arrayBody += ", "
		}
		arrayBody += constant.name
		if !compact {
			arrayBody += ",\n"
		}
		i++
	}
	arrayBody += "}"

	return "func " + funcName + "() " + arrayType + " { return " + arrayBody + "}", funcName, nil
}

func (g *Generator) aenerateConstFieldFunc(typ string, constants []constResult, export, nolint bool) (string, string, error) {
	var (
		funcName     = goName("Field", export)
		receiverVar  = "c"
		returnType   = BaseConstType
		returnNoCase = "\"\""
	)

	funcBody := "func (" + receiverVar + " " + typ + ") " + funcName + "() " + returnType
	funcBody += " {" + g.noLint(nolint) + "\n" +
		"switch " + receiverVar + " {\n" +
		""

	for _, constant := range constants {
		funcBody += "case " + constant.name + ":\n" +
			"return \"" + constant.field + "\"\n"
	}

	funcBody += "}\n"
	funcBody += "" +
		"return " + returnNoCase +
		"}\n"

	return funcBody, MethodName(typ, funcName), nil
}

func MethodName(typ, fun string) string { return typ + "." + fun }

func (g *Generator) GenerateConstValueFunc(model *struc.Model, pkg, typ string, constants []constResult, export, nolint, ref bool) (string, string, error) {
	var (
		funcName     = goName("Val", export)
		receiverVar  = "c"
		argName      = "s"
		argType      = "*" + getTypeName(model.TypeName, pkg)
		returnTypes  = "interface{}"
		returnNoCase = "nil"
		pref         = ""
	)

	if ref {
		pref = "&"
		funcName = goName("Ref", export)
	}

	funcBody := "func (" + receiverVar + " " + typ + ") " + funcName + "(" + argName + " " + argType + ") " + returnTypes
	funcBody += " {" + g.noLint(nolint) + "\n" +
		"switch " + receiverVar + " {\n" +
		""

	for _, constant := range constants {
		funcBody += "case " + constant.name + ":\n" +
			"return " + pref + argName + "." + constant.field + "\n"
	}

	funcBody += "}\n"
	funcBody += "" +
		"return " + returnNoCase +
		"}\n"

	return funcBody, MethodName(typ, funcName), nil
}
