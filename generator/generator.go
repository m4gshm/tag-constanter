package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"strings"
	"unicode"

	"github.com/m4gshm/fieldr/struc"
)

type Generator struct {
	Export     bool
	ExportVars bool
	ReturnRefs bool
	WrapType   bool
	Opts       *GenerateContentOptions
	buf        bytes.Buffer
	Name       string
}

type GenerateContentOptions struct {
	Fields               *bool
	Tags                 *bool
	TagsByFieldsMap      *bool
	TagValuesMap         *bool
	TagFieldsMap         *bool
	TagsValuesByFieldMap *bool

	GetFieldValue       *bool
	GetFieldValueByTag  *bool
	GetFieldValuesByTag *bool
	AsMap               *bool
	AsTagMap            *bool

	Strings *bool
}

func (g *Generator) printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) FormatSrc() ([]byte, error) {
	src := g.Src()
	fmtSrc, err := format.Source(src)
	if err != nil {
		return src, err
	}
	return fmtSrc, nil
}

func (g *Generator) Src() []byte {
	return g.buf.Bytes()
}

func (g *Generator) GenerateFile(str *struc.Struct) {
	g.Generate(str.PackageName, str.TypeName, str.TagNames, str.FieldNames, str.Fields)
}

const baseType = "string"

func (g *Generator) Generate(packageName string, typeName string, tagNames []struc.TagName, fieldNames []struc.FieldName,
	fields map[struc.FieldName]map[struc.TagName]struc.TagValue,
) {

	g.printf("// Code generated by '%s %s'; DO NOT EDIT.\n\n", g.Name, strings.Join(os.Args[1:], " "))
	g.printf("package %s\n", packageName)

	fieldType := getFieldType(typeName, g.Export)
	tagType := getTagType(typeName, g.Export)
	tagValueType := getTagValueType(typeName, g.Export)

	if g.WrapType {
		g.printf("type(\n")

		g.printf("%v %v\n", fieldType, baseType)
		g.printf("%v %v\n", arrayType(fieldType), "[]"+fieldType)
		g.printf("%v %v\n", tagType, baseType)
		g.printf("%v %v\n", arrayType(tagType), "[]"+tagType)
		g.printf("%v %v\n", tagValueType, baseType)
		g.printf("%v %v\n", arrayType(tagValueType), "[]"+tagValueType)

		g.printf(")\n")
	}

	g.printf("const(\n")

	g.generateFieldConstants(typeName, fieldNames)
	g.printf("\n")
	g.generateTagConstants(typeName, tagNames)
	g.printf("\n")
	g.generateTagFieldConstants(typeName, tagNames, fieldNames, fields)

	g.printf(")\n")

	g.printf("var(\n")

	opts := g.Opts
	if *opts.Fields {
		g.generateFieldsArrayVar(typeName, fieldNames)
	}

	if *opts.Tags {
		g.generateTagsArrayVar(typeName, tagNames)
	}

	if *opts.TagsByFieldsMap {
		g.generateFieldTagsMapVar(typeName, tagNames, fieldNames, fields)
	}

	if *opts.TagValuesMap {
		g.generateTagValuesMapVar(typeName, tagNames, fieldNames, fields)
	}

	if *opts.TagFieldsMap {
		g.generateTagFieldsMapVar(typeName, tagNames, fieldNames, fields)
	}

	if *opts.TagsValuesByFieldMap {
		g.generateFieldTagValueMapVar(fieldNames, tagNames, typeName, fields)
	}

	g.printf(")\n")

	if g.WrapType && *opts.Strings {
		g.generateArrayToStringsFunc(arrayType(fieldType), baseType)
		g.printf("\n")
		g.generateArrayToStringsFunc(arrayType(tagType), baseType)
		g.printf("\n")
		g.generateArrayToStringsFunc(arrayType(tagValueType), baseType)
		g.printf("\n")
	}

	returnRefs := g.ReturnRefs

	if *opts.GetFieldValue {
		g.generateGetFieldValueFunc(typeName, fieldNames, returnRefs)
		g.printf("\n")
	}
	if *opts.GetFieldValueByTag {
		g.generateGetFieldValueByTagFunc(typeName, fieldNames, tagNames, fields, returnRefs)
		g.printf("\n")
	}
	if *opts.GetFieldValuesByTag {
		g.generateGetFieldValuesByTagFunc(typeName, fieldNames, tagNames, fields, returnRefs)
		g.printf("\n")
	}
	if *opts.AsMap {
		g.generateAsMapFunc(typeName, fieldNames, returnRefs)
		g.printf("\n")
	}
	if *opts.AsTagMap {
		g.generateAsTagMapFunc(typeName, fieldNames, tagNames, fields, returnRefs)
		g.printf("\n")
	}

}

func arrayType(baseType string) string {
	return baseType + "s"
}

func getTagValueType(typeName string, export bool) string {
	return goName(typeName+"TagValue", export)
}

func getTagType(typeName string, export bool) string {
	return goName(typeName+"Tag", export)
}

func getFieldType(typeName string, export bool) string {
	return goName(typeName+"Field", export)
}

func goName(name string, export bool) string {
	first := rune(name[0])
	if export {
		first = unicode.ToUpper(first)
	} else {
		first = unicode.ToLower(first)
	}
	result := string(first) + name[1:]
	return result
}

func (g *Generator) generateFieldTagValueMapVar(fieldNames []struc.FieldName, tagNames []struc.TagName, typeName string, fields map[struc.FieldName]map[struc.TagName]struc.TagValue) {
	export := g.Export
	tagType := getTagType(typeName, export)
	fieldType := getFieldType(typeName, export)
	tagValueType := getTagValueType(typeName, export)

	var varValue string
	if g.WrapType {
		varValue = "map[" + fieldType + "]map[" + tagType + "]" + tagValueType + "{\n"
	} else {
		varValue = "map[" + baseType + "]map[" + baseType + "]" + baseType + "{\n"
	}
	for _, fieldName := range fieldNames {
		fieldConstName := getFieldConstName(typeName, fieldName, export)
		if g.WrapType {
			varValue += fieldConstName + ": map[" + tagType + "]" + tagValueType + "{"
		} else {
			varValue += fieldConstName + ": map[" + baseType + "]" + baseType + "{"
		}

		ti := 0
		for _, tagName := range tagNames {
			_, ok := fields[fieldName][tagName]
			if !ok {
				continue
			}
			if ti > 0 {
				varValue += ", "
			}

			tagConstName := getTagConstName(typeName, tagName, export)
			varValue += tagConstName + ": " + getTagValueConstName(typeName, tagName, fieldName, export)
			ti++
		}

		varValue += "},\n"
	}
	varValue += "}"

	varName := goName(typeName+"_FieldTagValue", g.ExportVars)

	g.printf("%v=%v\n\n", varName, varValue)
}

func (g *Generator) generateFieldTagsMapVar(typeName string, tagNames []struc.TagName, fieldNames []struc.FieldName, fields map[struc.FieldName]map[struc.TagName]struc.TagValue) {
	tagType := getTagType(typeName, g.Export)
	fieldType := getFieldType(typeName, g.Export)

	var varValue string
	if g.WrapType {
		varValue = "map[" + fieldType + "]" + arrayType(tagType) + "{\n"
	} else {
		varValue = "map[" + baseType + "][]" + baseType + "{\n"
	}
	for _, fieldName := range fieldNames {
		fieldConstName := getFieldConstName(typeName, fieldName, g.Export)

		if g.WrapType {
			varValue += fieldConstName + ": " + arrayType(tagType) + "{"
		} else {
			varValue += fieldConstName + ": []string{"
		}

		ti := 0
		for _, tagName := range tagNames {
			_, ok := fields[fieldName][tagName]
			if !ok {
				continue
			}

			if ti > 0 {
				varValue += ", "
			}
			tagConstName := getTagConstName(typeName, tagName, g.Export)
			varValue += tagConstName
			ti++
		}

		varValue += "},\n"
	}
	varValue += "}"

	varName := goName(typeName+"_FieldTags", g.ExportVars)

	g.printf("%v=%v\n\n", varName, varValue)
}

func (g *Generator) generateTagValuesMapVar(typeName string, tagNames []struc.TagName, fieldNames []struc.FieldName, fields map[struc.FieldName]map[struc.TagName]struc.TagValue) {
	tagType := getTagType(typeName, g.Export)
	tagValueType := getTagValueType(typeName, g.Export)
	var varValue string
	if g.WrapType {
		varValue = "map[" + tagType + "]" + arrayType(tagValueType) + "{\n"
	} else {
		varValue = "map[" + baseType + "][]" + baseType + "{\n"
	}
	for _, tagName := range tagNames {
		constName := getTagConstName(typeName, tagName, g.Export)

		if g.WrapType {
			varValue += constName + ": " + arrayType(tagValueType) + "{"
		} else {
			varValue += constName + ": []" + baseType + "{"
		}

		ti := 0
		for _, fieldName := range fieldNames {
			_, ok := fields[fieldName][tagName]
			if !ok {
				continue
			}

			if ti > 0 {
				varValue += ", "
			}
			tagConstName := getTagValueConstName(typeName, tagName, fieldName, g.Export)
			varValue += tagConstName
			ti++
		}

		varValue += "},\n"
	}
	varValue += "}"

	varName := goName(typeName+"_TagValues", g.ExportVars)

	g.printf("%v=%v\n\n", varName, varValue)
}

func (g *Generator) generateTagFieldsMapVar(typeName string, tagNames []struc.TagName, fieldNames []struc.FieldName, fields map[struc.FieldName]map[struc.TagName]struc.TagValue) {
	fieldType := getFieldType(typeName, g.Export)
	tagType := getTagType(typeName, g.Export)

	var varValue string
	if g.WrapType {
		varValue = "map[" + tagType + "]" + arrayType(fieldType) + "{\n"
	} else {
		varValue = "map[string][]string{\n"
	}
	for _, tagName := range tagNames {
		constName := getTagConstName(typeName, tagName, g.Export)

		if g.WrapType {
			varValue += constName + ": " + arrayType(fieldType) + "{"
		} else {
			varValue += constName + ": []string{"
		}

		ti := 0
		for _, field := range fieldNames {
			_, ok := fields[field][tagName]
			if !ok {
				continue
			}

			if ti > 0 {
				varValue += ", "
			}
			tagConstName := getFieldConstName(typeName, field, g.Export)
			varValue += tagConstName
			ti++
		}

		varValue += "},\n"
	}
	varValue += "}"

	varName := goName(typeName+"_TagFields", g.ExportVars)

	g.printf("%v=%v\n\n", varName, varValue)
}

func (g *Generator) generateTagFieldConstants(
	typeName string, tagNames []struc.TagName, fieldNames []struc.FieldName,
	fields map[struc.FieldName]map[struc.TagName]struc.TagValue,
) {
	for i, _tagName := range tagNames {
		if i > 0 {
			g.printf("\n")
		}
		for _, _fieldName := range fieldNames {
			_tagValue, ok := fields[_fieldName][_tagName]
			if ok {
				constName := getTagValueConstName(typeName, _tagName, _fieldName, g.Export)
				if g.WrapType {
					g.printf("%v=%v(\"%v\")\n", constName, getTagValueType(typeName, g.Export), _tagValue)
				} else {
					g.printf("%v=\"%v\"\n", constName, _tagValue)
				}
			}
		}
	}
}

func (g *Generator) generateFieldConstants(typeName string, fieldNames []struc.FieldName) {
	export := g.Export
	for _, name := range fieldNames {
		constName := getFieldConstName(typeName, name, export)
		if g.WrapType {
			g.printf("%v=%v(\"%v\")\n", constName, getFieldType(typeName, export), name)
		} else {
			g.printf("%v=\"%v\"\n", constName, name)
		}
	}
}

func (g *Generator) generateTagConstants(typeName string, tagNames []struc.TagName) {
	for _, name := range tagNames {
		constName := getTagConstName(typeName, name, g.Export)
		if g.WrapType {
			g.printf("%v=%v(\"%v\")\n", constName, getTagType(typeName, g.Export), name)
		} else {
			g.printf("%v=\"%v\"\n", constName, name)
		}
	}
}

func (g *Generator) generateFieldsArrayVar(typeName string, fieldNames []struc.FieldName) {
	fieldType := getFieldType(typeName, g.Export)
	var arrayVar string
	if g.WrapType {
		arrayVar = arrayType(fieldType) + "{"
	} else {
		arrayVar = "[]" + baseType + "{"
	}

	for i, fieldName := range fieldNames {
		if i > 0 {
			arrayVar += ", "
		}
		constName := getFieldConstName(typeName, fieldName, g.Export)
		arrayVar += constName
	}
	arrayVar += "}"
	varName := goName(typeName+"_Fields", g.ExportVars)
	g.printf("%v=%v\n\n", varName, arrayVar)
}

func (g *Generator) generateTagsArrayVar(typeName string, tagNames []struc.TagName) {
	var arrayVar string
	if g.WrapType {
		arrayVar = arrayType(getTagType(typeName, g.Export)) + "{"
	} else {
		arrayVar = "[]" + baseType + "{"
	}

	for i, tagName := range tagNames {
		if i > 0 {
			arrayVar += ", "
		}
		constName := getTagConstName(typeName, tagName, g.Export)
		arrayVar += constName
	}
	arrayVar += "}"
	varName := goName(typeName+"_Tags", g.ExportVars)
	g.printf("%v=%v\n\n", varName, arrayVar)
}

func (g *Generator) generateGetFieldValueFunc(typeName string, fieldNames []struc.FieldName, returnRefs bool) {

	var valType string
	if g.WrapType {
		valType = getFieldType(typeName, g.Export)
	} else {
		valType = "string"
	}

	valVar := "field"
	receiverVar := "v"
	receiverRef := asRefIfNeed(receiverVar, returnRefs)

	funcName := goName("GetFieldValue", g.Export)
	funcBody := "func (" + receiverVar + " *" + typeName + ") " + funcName + "(" + valVar + " " + valType + ") interface{} " +
		"{\n" + "switch " + valVar + " {\n"

	for _, fieldName := range fieldNames {
		fieldExpr := receiverRef + "." + string(fieldName)
		funcBody += "case " + getFieldConstName(typeName, fieldName, g.Export) + ":\n" +
			"return " + fieldExpr + "\n"
	}

	funcBody += "}\n" +
		"return nil" +
		"\n}\n"

	g.printf(funcBody)
}

func (g *Generator) generateGetFieldValueByTagFunc(typeName string, fieldNames []struc.FieldName, tagNames []struc.TagName, fields map[struc.FieldName]map[struc.TagName]struc.TagValue, returnRefs bool) {

	var valType string
	if g.WrapType {
		valType = getTagValueType(typeName, g.Export)
	} else {
		valType = "string"
	}

	valVar := "tag"
	receiverVar := "v"
	receiverRef := asRefIfNeed(receiverVar, returnRefs)

	funcName := goName("GetFieldValueByTagValue", g.Export)
	funcBody := "func (" + receiverVar + " *" + typeName + ") " + funcName + "(" + valVar + " " + valType + ") interface{} " +
		"{\n" + "switch " + valVar + " {\n"

	for _, fieldName := range fieldNames {
		fieldExpr := receiverRef + "." + string(fieldName)

		var caseExpr string
		for _, tagName := range tagNames {
			_, ok := fields[fieldName][tagName]
			if ok {
				if len(caseExpr) > 0 {
					caseExpr += ", "
				}
				caseExpr += getTagValueConstName(typeName, tagName, fieldName, g.Export)
			}
		}
		funcBody += "case " + caseExpr + ":\n" +
			"return " + fieldExpr + "\n"
	}

	funcBody += "}\n" +
		"return nil" +
		"\n}\n"

	g.printf(funcBody)
}

func (g *Generator) generateGetFieldValuesByTagFunc(typeName string, fieldNames []struc.FieldName, tagNames []struc.TagName, fields map[struc.FieldName]map[struc.TagName]struc.TagValue, returnRefs bool) {

	var valType string
	if g.WrapType {
		valType = getTagType(typeName, g.Export)
	} else {
		valType = baseType
	}

	valVar := "tag"
	receiverVar := "v"
	receiverRef := asRefIfNeed(receiverVar, returnRefs)

	resultType := "[]interface{}"

	funcName := goName("GetFieldValuesByTag", g.Export)
	funcBody := "func (" + receiverVar + " *" + typeName + ") " + funcName + "(" + valVar + " " + valType + ") " + resultType + " " +
		"{\n" + "switch " + valVar + " {\n"
	for _, tagName := range tagNames {

		caseExpr := getTagConstName(typeName, tagName, g.Export)
		fieldExpr := ""
		for _, fieldName := range fieldNames {
			_, ok := fields[fieldName][tagName]
			if ok {
				if len(fieldExpr) > 0 {
					fieldExpr += ", "
				}
				fieldExpr += receiverRef + "." + string(fieldName)
			}
		}
		if len(fieldExpr) > 0 {
			funcBody += "case " + caseExpr + ":\n" +
				"return " + resultType + "{" + fieldExpr + "}\n"
		}
	}

	funcBody += "}\n" +
		"return nil" +
		"\n}\n"

	g.printf(funcBody)
}

func asRefIfNeed(receiverVar string, returnRefs bool) string {
	receiverRef := receiverVar
	if returnRefs {
		receiverRef = "&" + receiverRef
	}
	return receiverRef
}

func (g *Generator) generateArrayToStringsFunc(arrayTypeName string, resultType string) {
	funcName := goName("Strings", g.Export)
	receiverVar := "v"
	g.printf("" +
		"func (" + receiverVar + " " + arrayTypeName + ") " + funcName + "() []" + resultType + " {\n" +
		"	strings := make([]" + resultType + ", 0, len(v))\n" +
		"	for i, val := range " + receiverVar + " {\n" +
		"		strings[i] = string(val)\n" +
		"		}\n" +
		"		return strings\n" +
		"	}\n")
}

func (g *Generator) generateAsMapFunc(typeName string, fieldNames []struc.FieldName, returnRefs bool) {
	export := g.Export

	receiverVar := "v"
	receiverRef := asRefIfNeed(receiverVar, returnRefs)

	keyType := baseType
	if g.WrapType {
		keyType = getFieldType(typeName, export)
	}

	funcName := goName("AsMap", export)
	funcBody := "" +
		"func (" + receiverVar + " *" + typeName + ") " + funcName + "() map[" + keyType + "]interface{} {\n" +
		"	return map[" + keyType + "]interface{}{\n"

	for _, fieldName := range fieldNames {
		funcBody += getFieldConstName(typeName, fieldName, export) + ": " + receiverRef + "." + string(fieldName) + ",\n"
	}
	funcBody += "" +
		"	}\n" +
		"}"

	g.printf(funcBody)
}

func (g *Generator) generateAsTagMapFunc(typeName string, fieldNames []struc.FieldName, tagNames []struc.TagName, fields map[struc.FieldName]map[struc.TagName]struc.TagValue, returnRefs bool) {
	receiverVar := "v"
	receiverRef := asRefIfNeed(receiverVar, returnRefs)

	keyType := baseType
	if g.WrapType {
		keyType = getTagValueType(typeName, g.Export)
	}

	valueType := "interface{}"

	varName := "tag"

	mapType := "map[" + keyType + "]" + valueType

	funcName := goName("AsTagMap", g.Export)
	funcBody := "" +
		"func (" + receiverVar + " *" + typeName + ") " + funcName + "(" + varName + " " + getTagType(typeName, g.Export) + ") " + mapType + " {\n" +
		"switch " + varName + " {\n" +
		""

	for _, tagName := range tagNames {
		funcBody += "case " + getTagConstName(typeName, tagName, g.Export) + ":\n" +
			"return " + mapType + "{\n"
		for _, fieldName := range fieldNames {
			_, ok := fields[fieldName][tagName]

			if ok {
				funcBody += getTagValueConstName(typeName, tagName, fieldName, g.Export) + ": " + receiverRef + "." + string(fieldName) + ",\n"
			}
		}

		funcBody += "}\n"
	}
	funcBody += "" +
		"	}\n" +
		"return nil" +
		"}"

	g.printf(funcBody)
}

func getTagConstName(typeName string, tag struc.TagName, export bool) string {
	return goName(typeName+"_"+string(tag), export)
}

func getTagValueConstName(typeName string, tag struc.TagName, field struc.FieldName, export bool) string {
	return goName(typeName+"_"+string(tag)+"_"+string(field), export)
}

func getFieldConstName(typeName string, fieldName struc.FieldName, export bool) string {
	return goName(typeName+"_"+string(fieldName), export)
}
