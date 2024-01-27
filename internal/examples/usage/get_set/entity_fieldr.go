// Code generated by 'fieldr'; DO NOT EDIT.

package get_set

import "time"

func (e *Entity[ID]) Id() ID {
	if e != nil {
		if be := e.BaseEntity; be != nil {
			return be.id
		}
	}

	var no ID
	return no
}

func (e *Entity[ID]) SetId(id ID) {
	if e != nil {
		if be := e.BaseEntity; be != nil {
			be.id = id
		}
	}
}

func (e *Entity[ID]) Name() string {
	if e != nil {
		return e.name
	}

	var no string
	return no
}

func (e *Entity[ID]) SetName(name string) {
	if e != nil {
		e.name = name
	}
}

func (e *Entity[ID]) Surname() string {
	if e != nil {
		return e.surname
	}

	var no string
	return no
}

func (e *Entity[ID]) SetSurname(surname string) {
	if e != nil {
		e.surname = surname
	}
}

func (e *Entity[ID]) Ts() time.Time {
	if e != nil {
		return e.ts
	}

	var no time.Time
	return no
}

func (e *Entity[ID]) SetTs(ts time.Time) {
	if e != nil {
		e.ts = ts
	}
}