package asmap

type StructField string

const (
	BaseStructID StructField = "ID"
	BaseStructTS StructField = "TS"
	Name         StructField = "Name"
	Surname      StructField = "Surname"
	NoTag        StructField = "NoTag"
	Address      StructField = "Address"
	FlatCardNum  StructField = "CardNum"
	FlatBank     StructField = "Bank"
)

func (s *Struct[n]) AsMap() map[StructField]interface{} {
	if s == nil {
		return nil
	}
	m := map[StructField]interface{}{}
	if bs := s.BaseStruct; bs != nil {
		m[BaseStructID] = bs.ID
	}
	if bs := s.BaseStruct; bs != nil {
		if ts := bs.TS; ts != nil {
			m[BaseStructTS] = ts
		}
	}
	m[Name] = s.Name
	m[Surname] = s.Surname
	m[NoTag] = s.NoTag
	if a := s.Address; a != nil {
		m[Address] = a.AsMap()
	}
	m[FlatCardNum] = s.Flat.CardNum
	m[FlatBank] = s.Flat.Bank
	return m
}
