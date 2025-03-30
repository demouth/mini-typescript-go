package printer

type EmitTextWriter interface {
	Clear()
	WriteStringLiteral(text string)
	WriteTrailingSemicolon(text string)
	WriteLine()
	String() string
}
