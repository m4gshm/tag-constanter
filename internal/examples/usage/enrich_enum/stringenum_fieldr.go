// Code generated by 'fieldr'; DO NOT EDIT.

package enrich_enum

func (s StringEnum) Name() string {
	switch s {
	case FIRST:
		return "FIRST"
	case SECOND:
		return "SECOND"
	case THIRD:
		return "THIRD"
	default:
		return ""
	}
}

func StringEnumAll() []StringEnum {
	return []StringEnum{
		FIRST,
		SECOND,
		THIRD,
	}
}

func StringEnumByName(name string) (e StringEnum, ok bool) {
	ok = true
	switch name {
	case "FIRST":
		e = FIRST
	case "SECOND":
		e = SECOND
	case "THIRD":
		e = THIRD
	default:
		ok = false
	}
	return
}

func StringEnumByValue(value string) (e StringEnum, ok bool) {
	ok = true
	switch value {
	case "first one":
		e = FIRST
	case "one more":
		e = SECOND
	case "any third":
		e = THIRD
	default:
		ok = false
	}
	return
}
