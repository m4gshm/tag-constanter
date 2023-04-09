// Code generated by 'fieldr'; DO NOT EDIT.

package internal

import (
	"example/get_set"
	"time"
)

func GetID[ID any](e *get_set.Entity[ID]) ID {
	if e != nil {
		if be := e.BaseEntity; be != nil {
			return be.ID
		}
	}

	var no ID
	return no
}

func SetID[ID any](e *get_set.Entity[ID], iD ID) {
	if e != nil {
		if be := e.BaseEntity; be != nil {
			be.ID = iD
		}
	}
}

func GetCode[ID any](e *get_set.Entity[ID]) string {
	if e != nil {
		if be := e.BaseEntity; be != nil {
			if rcae := be.RefCodeAwareEntity; rcae != nil {
				if cae := rcae.CodeAwareEntity; cae != nil {
					return cae.Code
				}
			}
		}
	}

	var no string
	return no
}

func SetCode[ID any](e *get_set.Entity[ID], code string) {
	if e != nil {
		if be := e.BaseEntity; be != nil {
			if rcae := be.RefCodeAwareEntity; rcae != nil {
				if cae := rcae.CodeAwareEntity; cae != nil {
					cae.Code = code
				}
			}
		}
	}
}

func GetNoDB[ID any](e *get_set.Entity[ID]) *get_set.NoDBFieldsEntity {
	if e != nil {
		return e.NoDB
	}

	var no *get_set.NoDBFieldsEntity
	return no
}

func SetNoDB[ID any](e *get_set.Entity[ID], noDB *get_set.NoDBFieldsEntity) {
	if e != nil {
		e.NoDB = noDB
	}
}

func GetValues[ID any](e *get_set.Entity[ID]) []int32 {
	if e != nil {
		return e.Values
	}

	var no []int32
	return no
}

func SetValues[ID any](e *get_set.Entity[ID], values []int32) {
	if e != nil {
		e.Values = values
	}
}

func GetTs[ID any](e *get_set.Entity[ID]) []*time.Time {
	if e != nil {
		return e.Ts
	}

	var no []*time.Time
	return no
}

func SetTs[ID any](e *get_set.Entity[ID], ts []*time.Time) {
	if e != nil {
		e.Ts = ts
	}
}

func GetEmbedded[ID any](e *get_set.Entity[ID]) get_set.EmbeddedEntity {
	if e != nil {
		return e.Embedded
	}

	var no get_set.EmbeddedEntity
	return no
}

func SetEmbedded[ID any](e *get_set.Entity[ID], embedded get_set.EmbeddedEntity) {
	if e != nil {
		e.Embedded = embedded
	}
}