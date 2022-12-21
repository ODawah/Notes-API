package schemas

type Note struct {
	UUID  string
	Title string
	Text  string
}

type AllNotes struct {
	Notes []Note
}
