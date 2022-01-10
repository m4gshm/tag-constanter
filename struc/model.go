package struc

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

const ReplaceableValueSeparator = "="
const KeyValueSeparator = ":"
const ListValuesSeparator = ","

var (
	excludeValues = map[TagName]map[TagValue]bool{}
)

type (
	TagName   = string
	TagValue  = string
	FieldName = string
	Package   struct{ Name, Path string }
	FieldType struct {
		Embedded, Ref  bool
		Name, FullName string
		Model          *Model
		Type           types.Type
	}

	//Model struct type model.
	Model struct {
		TypeName       string
		Package        Package
		FilePath       string
		FieldsTagValue map[FieldName]map[TagName]TagValue
		TagsFieldValue map[TagName]map[FieldName]TagValue
		FieldNames     []FieldName
		FieldsType     map[FieldName]FieldType
	}
)

// New - HierarchicalModel's default constructor.
func New(filePackages map[*ast.File]*packages.Package, files []*ast.File, fileSet *token.FileSet, typeName string) (*Model, error) {
	for _, file := range files {
		var (
			filePackage = filePackages[file]
			pkg         = filePackage.Types
		)
		lookup := pkg.Scope().Lookup(typeName)
		if lookup == nil {
			continue
		}
		var (
			fileInfo = fileSet.File(file.Pos())
			filePath = fileInfo.Name()
		)
		if builder, err := newBuilder(pkg, pkg, nil, typeName, filePath, handledStructs{}); err != nil {
			return nil, fmt.Errorf("new builder of %v: %w", typeName, err)
		} else if structModel, err := builder.newModel(lookup.Type()); err != nil {
			return nil, fmt.Errorf("new model of %v: %w", typeName, err)
		} else if structModel != nil {
			return structModel, nil
		}
	}
	return nil, nil
}

func newFieldTagValues(fieldTagNames []TagName, tagValues map[TagName]TagValue) map[TagName]TagValue {
	fieldTagValues := make(map[TagName]TagValue, len(fieldTagNames))
	for _, fieldTagName := range fieldTagNames {
		fieldTagValues[fieldTagName] = tagValues[fieldTagName]
	}
	return fieldTagValues
}

func parseTagValues(tags string) (map[TagName]TagValue, []TagName) {
	tagNames := make([]TagName, 0)
	tagValues := make(map[TagName]TagValue)

	var prevTagPos int
	tagValueLen := len(tags)
	for pos := 0; pos < tagValueLen; pos++ {
		character := rune(tags[pos])
		switch character {
		case '`', ' ':
			prevTagPos = pos + 1
		case ':':
			_tagName := TagName(tags[prevTagPos:pos])

			//parse TagValue
			pos++

			character = rune(tags[pos])
			tagValueBorder := '"'
			findEndBorder := false
			if character == tagValueBorder {
				pos++
				findEndBorder = true
			}
			tagDelim := ' '

			var endValuePos int
			for endValuePos = pos; endValuePos < tagValueLen; endValuePos++ {
				character = rune(tags[endValuePos])
				if findEndBorder && character == tagValueBorder {
					break
				} else if character == tagDelim {
					break
				}
			}

			tagContent := tags[pos:endValuePos]
			var excluded bool
			if excludedValues, ok := excludeValues[_tagName]; ok {
				excluded, ok = excludedValues[tagContent]
				excluded = excluded && ok
			}

			if !excluded {
				tagValues[_tagName] = tagContent
				tagNames = append(tagNames, _tagName)
			}

			prevTagPos = endValuePos
			pos = prevTagPos
		}
	}
	return tagValues, tagNames
}
