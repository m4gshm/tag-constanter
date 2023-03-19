package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/m4gshm/gollections/mutable/omap"
	"github.com/m4gshm/gollections/slice"
	"golang.org/x/tools/go/packages"

	"github.com/m4gshm/fieldr/command"
	"github.com/m4gshm/fieldr/generator"
	"github.com/m4gshm/fieldr/logger"
	"github.com/m4gshm/fieldr/params"
	"github.com/m4gshm/fieldr/use"
)

func usage(commandLine *flag.FlagSet) func() {
	return func() {
		out := commandLine.Output()
		_, _ = fmt.Fprintf(out, params.Name+" is a tool for generating constants, variables, functions and methods"+
			" based on a structure model: name, fields, tags\n")
		_, _ = fmt.Fprintf(out, "Usage of "+params.Name+":\n")
		_, _ = fmt.Fprintf(out, "\t"+params.Name+" [flags] command1 [command-flags] command2 [command-flags]... command [command-flags]\n")
		_, _ = fmt.Fprintf(out, "Use \"command --help\" to get help of this one\n")
		_, _ = fmt.Fprintf(out, "Flags:\n")
		commandLine.PrintDefaults()
		_, _ = fmt.Fprintf(out, "  --help\n")
		_, _ = fmt.Fprintf(out, "\tshow this message\n")
		command.PrintUsage()
	}
}

func main() {
	if err := run(); err != nil {
		var uErr *use.Error
		if errors.As(err, &uErr) {
			fmt.Fprintf(os.Stderr, "err: "+uErr.Error()+"\n")
			flag.CommandLine.Usage()
		} else {
			log.Fatal(err.Error())
		}
	}
}

func run() error {
	appFile := os.Args[0]
	appArgs := os.Args[1:]

	configParser := newConfigFlagSet(appFile)
	flag.CommandLine = configParser

	debugFlag := configParser.Bool("debug", false, "enable debug logging")
	buildTags := params.MultiVal(configParser, "buildTag", []string{"fieldr"}, "include build tag")
	inputs := params.InFlag(configParser)
	packagePattern := configParser.String("package", ".", "used package")

	commonTypeConfig := params.NewTypeConfig(configParser)
	if err := configParser.Parse(appArgs); err != nil {
		return err
	}

	logger.Init(*debugFlag)
	logger.Debugf("common type config: type '%v', output '%v'", commonTypeConfig.Type, commonTypeConfig.Output)

	commands, args, err := parseCommands(configParser.Args())
	if err != nil {
		return err
	}
	if len(commands) == 0 {
		logger.Debugf("no command line generator commands")
	}
	if len(args) > 0 {
		logger.Debugf("unspent command line args %v\n", args)
	}

	workDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	fileSet := token.NewFileSet()

	wdSrcPkg, err := extractPackage(fileSet, *buildTags, workDir)
	if err != nil {
		return err
	}
	wdSrcFiles := wdSrcPkg.Syntax

	typeConfigs := omap.Empty[params.TypeConfig, []*command.Command]()

	typeConfig := *commonTypeConfig

	filesCmdArgs, err := newFilesCommentsConfig(wdSrcFiles, fileSet)
	if err != nil {
		return err
	}

	notCmdLineType := len(typeConfig.Type) == 0

	for _, f := range filesCmdArgs {
		for _, cmt := range f.commentArgs {
			name := strings.Join(cmt.args, " ")
			configParser := newConfigFlagSet(name)
			commentConfig := params.NewTypeConfig(configParser)
			if err := configParser.Parse(cmt.args); err != nil {
				return err
			}

			if notCmdLineType {
				if len(commentConfig.Type) != 0 {
					typeConfig.Type = commentConfig.Type
					if len(typeConfig.Output) == 0 {
						typeConfig.Output = commentConfig.Output
					}
					if len(typeConfig.OutPackage) == 0 {
						typeConfig.OutPackage = commentConfig.OutPackage
					}
					if len(typeConfig.OutBuildTags) == 0 {
						typeConfig.OutBuildTags = commentConfig.OutBuildTags
					}
					logger.Debugf("init first type %+v by comment type %+v", typeConfig, *commentConfig)
				}
				notCmdLineType = false
			}

			if commentConfig.Type == typeConfig.Type && commentConfig.Output == typeConfig.Output {
				logger.Debugf("skip comment config because its type and out are equal to prev: comment config %+v, prev %+v", commentConfig, typeConfig)
				//skip
			} else if len(commentConfig.Type) == 0 && commentConfig.Output == typeConfig.Output {
				//skip
				logger.Debugf("skip comment config because its out is equal to prev: comment config %+v, prev %+v", commentConfig, typeConfig)
			} else if len(commentConfig.Type) != 0 || len(commentConfig.Output) != 0 {
				if len(commentConfig.Type) == 0 {
					(*commentConfig).Type = typeConfig.Type
				}

				logger.Debugf("detect another type %+v\n", *commentConfig)

				if len(commands) == 0 {
					logger.Debugf("no commands for type %v", typeConfig)
					typeConfig = *commentConfig
				} else {
					typeConfigs.Set(typeConfig, commands)
					logger.Debugf("set type %+v, commands %d\n", typeConfig, len(commands))
					typeConfig = *commentConfig
					commands = []*command.Command{}
				}
			}

			unusedArgs := configParser.Args()
			cmtCommands, cmtArgs, err := parseCommands(unusedArgs)

			// cmtCommands, cmtArgs, err := parseCommands(cmt.args)
			if err != nil {
				var uErr *use.Error
				if errors.As(err, &uErr) {
					return use.FileCommentErr(uErr.Error(), f.astFile, f.tokenFile, cmt.comment)
				}
				return err
			} else if len(cmtCommands) == 0 {
				// logger.Debugf("no comment generator commands: file %s, line: %d args %v\n", f.file.Name, cmt.comment.Pos(), cmtArgs)
			} else if len(cmtArgs) > 0 {
				logger.Debugf("unspent comment line args: %v\n", cmtArgs)
			}
			commands = append(commands, cmtCommands...)
		}
	}

	typeConfigs.Set(typeConfig, commands)

	logger.Debugf("set type last %+v, commands: %s\n", typeConfig, strings.Join(slice.Convert(commands, (*command.Command).Name), ", "))

	srcPkg, err := extractPackage(fileSet, *buildTags, *packagePattern)
	if err != nil {
		return err
	}
	srcFiles := srcPkg.Syntax

	filePackages := make(map[*ast.File]*packages.Package)
	for _, file := range srcFiles {
		filePackages[file] = srcPkg
	}

	srcFiles, err = loadSrcFiles(*inputs, *buildTags, fileSet, srcFiles, filePackages)
	if err != nil {
		return err
	}

	return typeConfigs.Track(func(typeConfig params.TypeConfig, commands []*command.Command) error {
		logger.Debugf("using type config %+v\n", typeConfig)

		outputName := typeConfig.Output
		if outputName == "" {
			typeName := typeConfig.Type
			if len(typeName) == 0 {
				return use.Err("no type arg")
			}

			baseName := typeName + params.DefaultFileSuffix
			outputName = strings.ToLower(baseName)
		}

		if outputName, err = filepath.Abs(outputName); err != nil {
			return err
		}

		var outFile *ast.File
		var outFileInfo *token.File
		for _, file := range srcFiles {
			if info := fileSet.File(file.Pos()); info.Name() == outputName {
				outFileInfo = info
				outFile = file
				break
			}
		}

		var outPkg *packages.Package
		if outFile != nil {
			outPkg = filePackages[outFile]
		} else {
			var stat os.FileInfo
			stat, err = os.Stat(outputName)
			noExists := errors.Is(err, os.ErrNotExist)
			if noExists {
				dir := filepath.Dir(outputName)
				outPkg, err = dirPackage(dir, nil)
				if err != nil {
					return err
				} else if outPkg == nil {
					return fmt.Errorf("canot detenrime output package, path '%v'", dir)
				}
			} else if err != nil {
				return err
			} else {
				if stat.IsDir() {
					return fmt.Errorf("output file is directory")
				}
				outFileSet := token.NewFileSet()
				outFile, outPkg, err = loadFile(outputName, nil, outFileSet)
				if err != nil {
					return err
				}
				if outFile != nil {
					pos := outFile.Pos()
					outFileInfo = outFileSet.File(pos)
					if outFileInfo == nil {
						return fmt.Errorf("error of reading metadata of output file %v", outputName)
					}
				}
			}
		}

		g := generator.New(params.Name, typeConfig.OutBuildTags, outFile, outFileInfo, outPkg)

		ctx := &command.Context{TypeConfig: typeConfig, Generator: g, FilePackages: filePackages, Files: srcFiles, FileSet: fileSet}
		for _, c := range commands {
			if err := c.Run(ctx); err != nil {
				return err
			}
		}

		outPackageName := generator.OutPackageName(typeConfig.OutPackage, outPkg)
		if err := g.WriteBody(outPackageName); err != nil {
			return err
		}

		src, fmtErr := g.FormatSrc()

		const userWriteOtherRead = fs.FileMode(0644)
		if writeErr := ioutil.WriteFile(outputName, src, userWriteOtherRead); writeErr != nil {
			return fmt.Errorf("writing output: %s", writeErr)
		} else if fmtErr != nil {
			return fmt.Errorf("go src code formatting error: %s", fmtErr)
		}
		return nil
	})
}

func newConfigFlagSet(name string) *flag.FlagSet {
	configParser := flag.NewFlagSet(name, flag.ContinueOnError)
	configParser.Usage = usage(configParser)
	return configParser
}

func parseCommands(args []string) ([]*command.Command, []string, error) {
	commands := []*command.Command{}
	for len(args) > 0 {
		cmd := args[0]
		args = args[1:]

		if c := command.Get(cmd); c == nil {
			return nil, args, use.Err("unknowd command '" + cmd + "'")
		} else if a, err := c.Parse(args); err != nil {
			return nil, nil, err
		} else {
			args = a
			commands = append(commands, c)
		}
	}
	return commands, args, nil
}

type fileCmdArgs struct {
	astFile     *ast.File
	tokenFile   *token.File
	commentArgs []commentCmdArgs
}

func newFilesCommentsConfig(files []*ast.File, fileSet *token.FileSet) ([]fileCmdArgs, error) {
	result := []fileCmdArgs{}
	for _, file := range files {
		ft := fileSet.File(file.Pos())
		if args, err := getFileCommentCmdArgs(file, ft); err != nil {
			return nil, err
		} else if len(args) > 0 {

			result = append(result, fileCmdArgs{astFile: file, tokenFile: ft, commentArgs: args})
		}
	}
	return result, nil
}

type commentCmdArgs struct {
	comment *ast.Comment
	args    []string
}

func getFileCommentCmdArgs(file *ast.File, fInfo *token.File) ([]commentCmdArgs, error) {
	result := []commentCmdArgs{}
	for _, commentGroup := range file.Comments {
		for _, comment := range commentGroup.List {
			if args, err := getCommentCmdArgs(comment.Text); err != nil {
				return nil, err
			} else if len(args) > 0 {
				name := fInfo.Name()
				line := fInfo.Line(comment.Pos())
				logger.Debugf("extracted comment args: file %s, line %d, args %v", name, line, args)
				result = append(result, commentCmdArgs{comment: comment, args: args})
			}
		}
	}
	return result, nil
}

func getCommentCmdArgs(text string) ([]string, error) {
	prefix := "//" + params.CommentConfigPrefix
	if len(text) > 0 && strings.HasPrefix(text, prefix) {
		configComment := text[len(prefix)+1:]
		if len(configComment) > 0 {
			logger.Debugf("split comment args '%s'", configComment)
			if args, err := splitArgs(configComment); err != nil {
				return nil, fmt.Errorf("split cofig comment %v; %w", text, err)
			} else {
				logger.Debugf("comment args count %d, '%s'", len(args), strings.Join(args, ","))
				return args, nil
			}
		}
	}
	return nil, nil
}

func splitArgs(rawArgs string) ([]string, error) {
	var args []string
	for {
		rawArgs = strings.TrimLeft(rawArgs, " ")
		if len(rawArgs) == 0 {
			break
		}
		symbols := []rune(rawArgs)
		if symbols[0] == '"' {
			finished := false
			//start parsing quoted string
		quoted:
			for i := 1; i < len(symbols); i++ {
				c := symbols[i]
				switch c {
				case '\\':
					if i+1 == len(symbols) {
						return nil, errors.New("unexpected backslash at the end")
					}
					i++
				case '"':
					part := rawArgs[0 : i+1]
					arg, err := strconv.Unquote(part)
					if err != nil {
						return nil, fmt.Errorf("unquote string: %s: %w", part, err)
					}
					args = append(args, arg)
					rawArgs = string(symbols[i+1:])
					//finish parsing quoted string
					finished = true
					break quoted
				}
			}
			if !finished {
				return nil, errors.New("unclosed quoted string")
			}
		} else {
			i := strings.Index(rawArgs, " ")
			if i < 0 {
				i = len(rawArgs)
			}
			args = append(args, rawArgs[0:i])
			rawArgs = rawArgs[i:]
		}
	}
	return args, nil
}

func loadSrcFiles(inputs []string, buildTags []string, fileSet *token.FileSet, files []*ast.File, filePackages map[*ast.File]*packages.Package) ([]*ast.File, error) {
	for _, srcFile := range inputs {
		file, pkg, err := loadFile(srcFile, buildTags, fileSet)
		if err != nil {
			return nil, err
		}
		if _, ok := filePackages[file]; !ok {
			files = append(files, file)
			filePackages[file] = pkg
		}
	}
	return files, nil
}

func loadFile(srcFile string, buildTags []string, fileSet *token.FileSet) (*ast.File, *packages.Package, error) {
	isAbs := filepath.IsAbs(srcFile)
	if !isAbs {
		absFile, err := filepath.Abs(srcFile)
		if err != nil {
			return nil, nil, err
		}
		srcFile = absFile
	}
	file, err := parser.ParseFile(fileSet, srcFile, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}
	dir := filepath.Dir(srcFile)
	pkg, err := dirPackage(dir, buildTags)
	if err != nil {
		return nil, nil, err
	}
	return file, pkg, err
}

func dirPackage(dir string, buildTags []string) (*packages.Package, error) {
	pack, err := packages.Load(&packages.Config{Mode: packageMode, BuildFlags: buildTagsArg(buildTags)}, dir)
	if err != nil {
		return nil, err
	}
	for _, p := range pack {
		return p, nil
	}
	return nil, nil
}

func isDir(name string) (bool, error) {
	info, err := os.Stat(name)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

const packageMode = packages.NeedSyntax | packages.NeedModule | packages.NeedName | packages.NeedTypesInfo | packages.NeedTypes

func extractPackage(fileSet *token.FileSet, buildTags []string, patterns ...string) (*packages.Package, error) {
	_packages, err := packages.Load(&packages.Config{
		Fset: fileSet, Mode: packageMode, BuildFlags: buildTagsArg(buildTags),
	}, patterns...)
	if err != nil {
		return nil, err
	}
	if len(_packages) != 1 {
		return nil, fmt.Errorf("%d packages found", len(_packages))
	}
	pack := _packages[0]
	if errs := pack.Errors; len(errs) > 0 {
		logger.Debugf("package error; %v", errs[0])
	}
	return pack, nil
}

func buildTagsArg(buildTags []string) []string {
	return []string{fmt.Sprintf("-tags=%s", strings.Join(buildTags, " "))}
}
