package foundation

type contextKey int

const (
	UserIDKey contextKey = iota
	ClaimsKey
	NoteKey
)
