package sql

type col string

const (
	colID      col = "id"
	colName    col = "name"
	colSurname col = "surname"
	colValues  col = "values"
	colTs      col = "ts"
	colVersion col = "version"
)

func cols() []col {
	return []col{
		colID,
		colName,
		colSurname,
		colValues,
		colTs,
		colVersion,
	}
}

func (c col) field() string {
	switch c {
	case colID:
		return "BaseEntity.ID"
	case colName:
		return "Name"
	case colSurname:
		return "Surname"
	case colValues:
		return "Values"
	case colTs:
		return "Ts"
	case colVersion:
		return "Versioned.Version"
	}
	return ""
}

func (c col) val(s *Entity) interface{} {
	switch c {
	case colID:
		return s.BaseEntity.ID
	case colName:
		return s.Name
	case colSurname:
		return s.Surname
	case colValues:
		return s.Values
	case colTs:
		return s.Ts
	case colVersion:
		return s.Versioned.Version
	}
	return nil
}

func (c col) ref(s *Entity) interface{} {
	switch c {
	case colID:
		return &s.BaseEntity.ID
	case colName:
		return &s.Name
	case colSurname:
		return &s.Surname
	case colValues:
		return &s.Values
	case colTs:
		return &s.Ts
	case colVersion:
		return &s.Versioned.Version
	}
	return nil
}
