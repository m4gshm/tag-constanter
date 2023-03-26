// Code generated by 'fieldr'; DO NOT EDIT.

//go:build integration
// +build integration

package builder

import (
	"bytes"
	"example/sql_base"
	"time"
)

type EntityBuilder[ID any] struct {
	ID           ID
	Code         string
	ForeignID    ID
	Schema       string
	Version      int
	NoDB         *NoDBFieldsEntity
	Name         StringBasedType[string]
	Surname      string
	Values       []int32
	Ts           []*time.Time
	Versioned    sql_base.VersionedEntity
	Chan         chan map[time.Time]string
	SomeMap      map[StringBasedType[string]]bytes.Buffer
	Embedded     EmbeddedEntity
	OldForeignID *ForeignIDAwareEntity[ID]
}

func NewEntityBuilder[ID any]() *EntityBuilder[ID] {
	return &EntityBuilder[ID]{}
}

func (b *EntityBuilder[ID]) Build() *Entity[ID] {
	if b == nil {
		return &Entity[ID]{}
	}
	return &Entity[ID]{
		BaseEntity: &BaseEntity[ID]{
			ID: b.ID,
			RefCodeAwareEntity: &RefCodeAwareEntity{
				CodeAwareEntity: &CodeAwareEntity{
					Code: b.Code,
				},
			},
			ForeignIDAwareEntity: ForeignIDAwareEntity[ID]{
				ForeignID: b.ForeignID,
			},
		},
		Metadata: Metadata{
			Schema:  b.Schema,
			Version: b.Version,
		},
		NoDB:         b.NoDB,
		Name:         b.Name,
		Surname:      b.Surname,
		Values:       b.Values,
		Ts:           b.Ts,
		Versioned:    b.Versioned,
		Chan:         b.Chan,
		SomeMap:      b.SomeMap,
		Embedded:     b.Embedded,
		OldForeignID: b.OldForeignID,
	}
}

func (b *EntityBuilder[ID]) SetID(iD ID) *EntityBuilder[ID] {
	if b != nil {
		b.ID = iD
	}
	return b
}

func (b *EntityBuilder[ID]) SetCode(code string) *EntityBuilder[ID] {
	if b != nil {
		b.Code = code
	}
	return b
}

func (b *EntityBuilder[ID]) SetForeignID(foreignID ID) *EntityBuilder[ID] {
	if b != nil {
		b.ForeignID = foreignID
	}
	return b
}

func (b *EntityBuilder[ID]) SetSchema(schema string) *EntityBuilder[ID] {
	if b != nil {
		b.Schema = schema
	}
	return b
}

func (b *EntityBuilder[ID]) SetVersion(version int) *EntityBuilder[ID] {
	if b != nil {
		b.Version = version
	}
	return b
}

func (b *EntityBuilder[ID]) SetNoDB(noDB *NoDBFieldsEntity) *EntityBuilder[ID] {
	if b != nil {
		b.NoDB = noDB
	}
	return b
}

func (b *EntityBuilder[ID]) SetName(name StringBasedType[string]) *EntityBuilder[ID] {
	if b != nil {
		b.Name = name
	}
	return b
}

func (b *EntityBuilder[ID]) SetSurname(surname string) *EntityBuilder[ID] {
	if b != nil {
		b.Surname = surname
	}
	return b
}

func (b *EntityBuilder[ID]) SetValues(values []int32) *EntityBuilder[ID] {
	if b != nil {
		b.Values = values
	}
	return b
}

func (b *EntityBuilder[ID]) SetTs(ts []*time.Time) *EntityBuilder[ID] {
	if b != nil {
		b.Ts = ts
	}
	return b
}

func (b *EntityBuilder[ID]) SetVersioned(versioned sql_base.VersionedEntity) *EntityBuilder[ID] {
	if b != nil {
		b.Versioned = versioned
	}
	return b
}

func (b *EntityBuilder[ID]) SetChan(chan_ chan map[time.Time]string) *EntityBuilder[ID] {
	if b != nil {
		b.Chan = chan_
	}
	return b
}

func (b *EntityBuilder[ID]) SetSomeMap(someMap map[StringBasedType[string]]bytes.Buffer) *EntityBuilder[ID] {
	if b != nil {
		b.SomeMap = someMap
	}
	return b
}

func (b *EntityBuilder[ID]) SetEmbedded(embedded EmbeddedEntity) *EntityBuilder[ID] {
	if b != nil {
		b.Embedded = embedded
	}
	return b
}

func (b *EntityBuilder[ID]) SetOldForeignID(oldForeignID *ForeignIDAwareEntity[ID]) *EntityBuilder[ID] {
	if b != nil {
		b.OldForeignID = oldForeignID
	}
	return b
}

func (i *Entity[ID]) ToBuilder() *EntityBuilder[ID] {
	if i == nil {
		return &EntityBuilder[ID]{}
	}
	var (
		i_BaseEntity_ID                  ID
		r_CodeAwareEntity_Code           string
		r_ForeignIDAwareEntity_ForeignID ID
	)
	if r := i.BaseEntity; r != nil {
		i_BaseEntity_ID = r.ID
		if r := r.RefCodeAwareEntity; r != nil {
			if r := r.CodeAwareEntity; r != nil {
				r_CodeAwareEntity_Code = r.Code
			}
		}
		r_ForeignIDAwareEntity_ForeignID = r.ForeignIDAwareEntity.ForeignID
	}

	return &EntityBuilder[ID]{
		ID:           i_BaseEntity_ID,
		Code:         r_CodeAwareEntity_Code,
		ForeignID:    r_ForeignIDAwareEntity_ForeignID,
		Schema:       i.Metadata.Schema,
		Version:      i.Metadata.Version,
		NoDB:         i.NoDB,
		Name:         i.Name,
		Surname:      i.Surname,
		Values:       i.Values,
		Ts:           i.Ts,
		Versioned:    i.Versioned,
		Chan:         i.Chan,
		SomeMap:      i.SomeMap,
		Embedded:     i.Embedded,
		OldForeignID: i.OldForeignID,
	}
}
