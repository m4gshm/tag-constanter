# Fieldr

Generator of various enumerated constants, types, functions based on structure fields and their tags.

## Supported commands

* enum-const - generate constants based on template applied to struct fields.
* as-map - generates method or functon that converts the struct type to a map.

## Enum constant example

create a file `entity.go` with the below content:

```go
package usage

//go:generate fieldr -type Entity enum-const -val .json -list jsons
type Entity struct {
    Id   int    `json:"id"`
    Name string `json:"name"`
}
```

then running this command in the same directory:

```bash script
go generate .
```

will be generated `entity_fieldr.go` file with the next content:

```go
// Code generated by 'fieldr'; DO NOT EDIT.

package usage

const (
    entityJsonId   = "id"
    entityJsonName = "name"
)

func jsons() []string { return []string{entityJsonId, entityJsonName} }
```

this consist of constants based on the `json` tag contents and the method `jsons` that enumerates these constants.

To get extended help of the command, use the following:

```bash
fields enum-const --help
```

 See more examples [here](./examples/)