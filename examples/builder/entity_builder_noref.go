// Code generated by 'fieldr'; DO NOT EDIT.

package builder

import (
	"bytes"
	"example/sql_base"
	"time"
)

type EntityBuilderVal[ID any] struct {
	ID        ID
	Code      string
	ForeignID ID
	NoDB      *NoDBFieldsEntity
	Name      StringBasedType[string]
	Surname   string
	Values    []int32
	Ts        []*time.Time
	Versioned sql_base.VersionedEntity
	Chan      chan map[time.Time]string
	SomeMap   map[StringBasedType[string]]bytes.Buffer
	Embedded  EmbeddedEntity
}

func NewEntityBuilderVal[ID any]() *EntityBuilderVal[ID] {
	return &EntityBuilderVal[ID]{}
}

func (b EntityBuilderVal[ID]) Build() Entity[ID] {
	return Entity[ID]{
		BaseEntity: &BaseEntity[ID]{
			ID: b.ID,
			CodeAwareEntity: &CodeAwareEntity{
				Code: b.Code,
			},
			ForeignIDAwareEntity: ForeignIDAwareEntity[ID]{
				ForeignID: b.ForeignID,
			},
		},
		NoDB:      b.NoDB,
		Name:      b.Name,
		Surname:   b.Surname,
		Values:    b.Values,
		Ts:        b.Ts,
		Versioned: b.Versioned,
		Chan:      b.Chan,
		SomeMap:   b.SomeMap,
		Embedded:  b.Embedded,
	}
}

func (b EntityBuilderVal[ID]) SetID(iD ID) EntityBuilderVal[ID] {
	b.ID = iD
	return b
}

func (b EntityBuilderVal[ID]) SetCode(code string) EntityBuilderVal[ID] {
	b.Code = code
	return b
}

func (b EntityBuilderVal[ID]) SetForeignID(foreignID ID) EntityBuilderVal[ID] {
	b.ForeignID = foreignID
	return b
}

func (b EntityBuilderVal[ID]) SetNoDB(noDB *NoDBFieldsEntity) EntityBuilderVal[ID] {
	b.NoDB = noDB
	return b
}

func (b EntityBuilderVal[ID]) SetName(name StringBasedType[string]) EntityBuilderVal[ID] {
	b.Name = name
	return b
}

func (b EntityBuilderVal[ID]) SetSurname(surname string) EntityBuilderVal[ID] {
	b.Surname = surname
	return b
}

func (b EntityBuilderVal[ID]) SetValues(values []int32) EntityBuilderVal[ID] {
	b.Values = values
	return b
}

func (b EntityBuilderVal[ID]) SetTs(ts []*time.Time) EntityBuilderVal[ID] {
	b.Ts = ts
	return b
}

func (b EntityBuilderVal[ID]) SetVersioned(versioned sql_base.VersionedEntity) EntityBuilderVal[ID] {
	b.Versioned = versioned
	return b
}

func (b EntityBuilderVal[ID]) SetChan(chan_ chan map[time.Time]string) EntityBuilderVal[ID] {
	b.Chan = chan_
	return b
}

func (b EntityBuilderVal[ID]) SetSomeMap(someMap map[StringBasedType[string]]bytes.Buffer) EntityBuilderVal[ID] {
	b.SomeMap = someMap
	return b
}

func (b EntityBuilderVal[ID]) SetEmbedded(embedded EmbeddedEntity) EntityBuilderVal[ID] {
	b.Embedded = embedded
	return b
}
