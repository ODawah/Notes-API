package schemas

type Note struct {
	UUID  string `json:"uuid"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type AllNotes struct {
	Notes []Note
}

type GetNote struct {
	Title string `json:"title"`
}
