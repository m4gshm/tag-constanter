// Code generated by 'fieldr'; DO NOT EDIT.

package enrich_enum

func (e Enum) Name() string {
	switch e {
	case AA:
		return "AA"
	case BB:
		return "BB"
	case CC:
		return "CC"
	case DD:
		return "DD"
	default:
		return ""
	}
}

func EnumAll() []Enum {
	return []Enum{
		AA,
		BB,
		CC,
		DD,
	}
}

func EnumByName(name string) (v Enum, ok bool) {
	ok = true
	switch name {
	case "AA":
		v = AA
	case "BB":
		v = BB
	case "CC":
		v = CC
	case "DD":
		v = DD
	default:
		ok = false
	}
	return
}

func EnumByValue(value int) (e Enum, ok bool) {
	ok = true
	switch value {
	case 1:
		e = AA
	case 2:
		e = BB
	case 3:
		e = CC
	case 4:
		e = DD
	default:
		ok = false
	}
	return
}
