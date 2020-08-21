package service

import "github.com/twcrone/inventoryservice/database"

type Note struct {
	Id   int    `json: "id"`
	Note string `json: "message"`
}

func GetAllNotes() []Note {
	results, err := database.DbConn.Query(`select id, note from notes`)
	if err != nil {
		return nil
	}
	defer results.Close()
	notes := make([]Note, 0)
	for results.Next() {
		var note Note
		_ = results.Scan(&note.Id, &note.Note)
		notes = append(notes, note)
	}
	return notes
}
