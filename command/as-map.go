package command

import (
	"flag"

	"github.com/m4gshm/fieldr/coderewriter"
	"github.com/m4gshm/fieldr/generator"
	"github.com/m4gshm/fieldr/params"
	"github.com/m4gshm/fieldr/struc"
)

func NewAsMapMethod() *Command {
	const (
		cmdName    = "as-map"
		genContent = "method/function"
	)

	const transformerTriggers = "<no condition (empty)>, " + string(generator.RewriteTriggerType) + ", " + string(generator.RewriteTriggerField)

	var transformFieldValueFormat = "trigger" + struc.KeyValueSeparator + "trigger_value" + struc.KeyValueSeparator + "engine" +
		struc.ReplaceableValueSeparator + "engine_format" + "; supported triggers '" + transformerTriggers +
		"', engine '" + string(generator.RewriteEngineFmt) + "'"

	var (
		flagSet             = flag.NewFlagSet(cmdName, flag.ContinueOnError)
		name                = flagSet.String("name", "", "function/method name")
		export              = params.Export(flagSet)
		snake               = params.Snake(flagSet)
		keyType             = flagSet.String("key-type", "", "generated constants type, use "+generator.Autoname+" for autoname")
		ref                 = flagSet.Bool("ref", false, "use struct field references in generated method")
		fun                 = flagSet.Bool("func", false, "generate function in place of struct method")
		all                 = flagSet.Bool("all", false, "use exported and private fields in generated "+genContent)
		nolint              = params.Nolint(flagSet)
		hardcode            = flagSet.Bool("hardcode", false, "hardcode field name in generated "+genContent+" (don't generate constants based on field name)")
		fieldValueRewriters = params.MultiVal(flagSet, "rewrite", []string{}, "field value rewriting applied to generated "+genContent+"; "+
			"format - "+transformFieldValueFormat)
		flats = params.MultiVal(flagSet, "flat", []string{}, "apply generator to fields of nested structs")
	)

	c := New(
		cmdName, "generates method or functon that converts the struct type to a map",
		flagSet,
		func(context *Context) error {
			g := context.Generator
			model, err := context.Model()
			if err != nil {
				return err
			}
			// excludedFields := make(map[struc.FieldName]struct{})

			kType := *keyType
			if kType == generator.Autoname {
				kType = generator.GetFieldType(model.TypeName, *export, *snake)
				if err := g.AddType(kType, generator.BaseConstType); err != nil {
					return err
				}
			} else if len(kType) == 0 {
				kType = generator.BaseConstType
			}

			flatsSet := toSet(*flats)
			if constants, err := g.GenerateFieldConstants(model, kType, *export, *snake, *all, flatsSet); err != nil {
				return err
			} else if rewriter, err := coderewriter.New(*fieldValueRewriters); err != nil {
				return err
			} else if _, funcName, funcBody, err := g.GenerateAsMapFunc(
				model, *name, kType, constants /*excludedFields, */, flatsSet, rewriter, *export, *snake, *ref, *fun, *nolint, *hardcode,
			); err != nil {
				return err
			} else if err := g.AddFuncOrMethod(funcName, funcBody); err != nil {
				return err
			}
			return nil

		},
	)
	return c
}
