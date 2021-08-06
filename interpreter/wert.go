package interpreter

type vollwert struct {
	links  lvalue
	rechts wert
}

func vrechts(v wert) *vollwert {
	return &vollwert{rechts: v}
}

type wert interface{ isWert() }

type IsWert struct{}

func (IsWert) isWert() {}

type logik struct {
	IsWert

	val bool
}

type ganz struct {
	IsWert

	val int64
}

type ptr struct {
	IsWert

	w **wert
}

type lvalue struct {
	IsWert

	w **wert
}

type struktur struct {
	IsWert

	fields map[string]*wert
}
