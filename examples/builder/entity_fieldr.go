package builder

import (
	"bytes"
	"example/sql_base"
	"time"
)

type EntityBuilder[ID any] struct {
	ID        ID
	NoDB      *NoDBFieldsEntity
	Name      StringBasedType[string]
	Surname   string
	Values    []int32
	Ts        []*time.Time
	Versioned sql_base.VersionedEntity
	Chan      chan map[time.Time]string
	SomeMap   map[StringBasedType[string]]bytes.Buffer
}

func (b EntityBuilder[ID]) Build() *Entity[ID] {
	return &Entity[ID]{
		BaseEntity: &BaseEntity[ID]{
			ID: b.ID,
		},
		NoDB:      b.NoDB,
		Name:      b.Name,
		Surname:   b.Surname,
		Values:    b.Values,
		Ts:        b.Ts,
		Versioned: b.Versioned,
		Chan:      b.Chan,
		SomeMap:   b.SomeMap,
	}
}
func (b *EntityBuilder[ID]) SetID(iD ID) *EntityBuilder[ID] {
	b.ID = iD
	return b
}

func (b *EntityBuilder[ID]) SetNoDB(noDB *NoDBFieldsEntity) *EntityBuilder[ID] {
	b.NoDB = noDB
	return b
}

func (b *EntityBuilder[ID]) SetName(name StringBasedType[string]) *EntityBuilder[ID] {
	b.Name = name
	return b
}

func (b *EntityBuilder[ID]) SetSurname(surname string) *EntityBuilder[ID] {
	b.Surname = surname
	return b
}

func (b *EntityBuilder[ID]) SetValues(values []int32) *EntityBuilder[ID] {
	b.Values = values
	return b
}

func (b *EntityBuilder[ID]) SetTs(ts []*time.Time) *EntityBuilder[ID] {
	b.Ts = ts
	return b
}

func (b *EntityBuilder[ID]) SetVersioned(versioned sql_base.VersionedEntity) *EntityBuilder[ID] {
	b.Versioned = versioned
	return b
}

func (b *EntityBuilder[ID]) SetChan(chan_ chan map[time.Time]string) *EntityBuilder[ID] {
	b.Chan = chan_
	return b
}

func (b *EntityBuilder[ID]) SetSomeMap(someMap map[StringBasedType[string]]bytes.Buffer) *EntityBuilder[ID] {
	b.SomeMap = someMap
	return b
}
