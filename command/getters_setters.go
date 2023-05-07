package command

import (
	"flag"
	"fmt"

	"github.com/m4gshm/fieldr/generator"
	"github.com/m4gshm/fieldr/logger"
	"github.com/m4gshm/fieldr/params"
	"github.com/m4gshm/fieldr/struc"
	"github.com/m4gshm/gollections/op"
)

func NewGettersSetters() *Command {
	const (
		cmdName = "get-set"
	)
	var (
		flagSet         = flag.NewFlagSet(cmdName, flag.ContinueOnError)
		getPrefix       = flagSet.String("get-prefix", "", "getter methods prefix")
		setPrefix       = flagSet.String("set-prefix", "Set", "setter methods prefix")
		noExportMethods = flagSet.Bool("no-export", false, "no export generated methods")
		noRefReceiver   = flagSet.Bool("no-ref", false, "use value type (not pointer) for methods receiver")
		accessors       = flagSet.String("accessors", "get-set", "full access methods or getter or setter only (supported: get-set, get, set)")
		nolint          = params.Nolint(flagSet)
	)

	return New(
		cmdName, "generates getters, setters for a structure type",
		flagSet,
		func(context *Context) error {

			getters, setters := false, false
			switch *accessors {
			case "get-set":
				getters, setters = true, true
			case "get":
				getters = true
			case "set":
				setters = true
			default:
				return fmt.Errorf("usupported accessors '%s'", *accessors)
			}

			model, err := context.Model()
			if err != nil {
				return err
			}
			g := context.Generator
			pkgName, err := g.GetPackageName(model.Package.Name, model.Package.Path)
			if err != nil {
				return err
			}

			rec := generator.TypeReceiverVar(model.TypeName)
			fmn, fmb, err := generateGettersSetters(g, model, model, pkgName, rec, *getPrefix, *setPrefix, getters, setters, !(*noRefReceiver), !(*noExportMethods), *nolint, nil)
			if err != nil {
				return err
			}

			for i := range fmn {
				fieldMethodName := fmn[i]
				fieldMethodBody := fmb[i]
				if err := g.AddMethod(model.TypeName, fieldMethodName, fieldMethodBody); err != nil {
					return err
				}
			}
			return nil
		},
	)
}

func generateGettersSetters(
	g *generator.Generator, baseModel, fieldsModel *struc.Model, pkgName, receiverVar, getterPrefix, setterPrefix string,
	getters, setters, isReceiverReference, exportMethods, nolint bool, parentFieldInfo []generator.FieldInfo,
) ([]string, []string, error) {
	logger.Debugf("generate getters, setters: receiver %s, type %s, getterPrefix %s setterPrefix %s", receiverVar, baseModel.TypeName, getterPrefix, setterPrefix)
	fieldMethodBodies := []string{}
	fieldMethodNames := []string{}
	for _, fieldName := range fieldsModel.FieldNames {
		fieldType := fieldsModel.FieldsType[fieldName]
		if len(pkgName) > 0 {
			if !generator.IsExported(fieldName) {
				logger.Debugf("cannot generate getter, setter for private field %s for package %s", fieldName, pkgName)
				continue
			}

			if m := fieldType.Model; m != nil {
				if !generator.IsExported(m.TypeName) {
					logger.Debugf("cannot generate getter, setter for field %s with private type % for package %s", fieldName, m.TypeName, pkgName)
					continue
				}
			}
		}

		if fieldType.Embedded {
			ebmeddedFieldMethodNames, ebmeddedFieldMethodBodies, err := generateGettersSetters(
				g, baseModel, fieldType.Model, pkgName, receiverVar, getterPrefix, setterPrefix, getters, setters, isReceiverReference, exportMethods, nolint,
				append(parentFieldInfo, generator.FieldInfo{Name: fieldType.Name, Type: fieldType}))
			if err != nil {
				return nil, nil, err
			}
			fieldMethodBodies = append(fieldMethodBodies, ebmeddedFieldMethodBodies...)
			fieldMethodNames = append(fieldMethodNames, ebmeddedFieldMethodNames...)
		} else {
			fullFieldType, err := g.GetFullFieldTypeName(fieldType, false)
			if err != nil {
				return nil, nil, err
			}
			suffix := generator.LegalIdentName(generator.IdentName(fieldName, true))

			if len(getterPrefix) == 0 || getterPrefix == generator.Autoname {
				getterPrefix = op.IfElse(suffix == fieldName, "Get", "")
			}
			if getters {
				getterName := generator.IdentName(getterPrefix+suffix, exportMethods)
				logger.Debugf("getter %s", getterName)
				getterBody := generator.GenerateGetter(baseModel, pkgName, receiverVar, getterName, fieldName, fullFieldType, g.OutPkgPath, nolint, isReceiverReference, parentFieldInfo)
				fieldMethodBodies = append(fieldMethodBodies, getterBody)
				fieldMethodNames = append(fieldMethodNames, getterName)
			}
			if setters {
				setterName := generator.IdentName(setterPrefix+suffix, exportMethods)
				logger.Debugf("setter %s", setterName)
				setterBody := generator.GenerateSetter(baseModel, pkgName, receiverVar, setterName, fieldName, fullFieldType, g.OutPkgPath, nolint, isReceiverReference, parentFieldInfo)
				fieldMethodBodies = append(fieldMethodBodies, setterBody)
				fieldMethodNames = append(fieldMethodNames, setterName)
			}
		}
	}
	return fieldMethodNames, fieldMethodBodies, nil
}
