package models

type Character struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	// ... (Alle deine Felder aus dem Pflichtenheft)
	Fertigkeiten []Fertigkeit `gorm:"foreignKey:CharacterID" json:"fertigkeiten"`
	Zauber       []Zauber     `gorm:"foreignKey:CharacterID" json:"zauber"`
	// ...
	// Für Zeitstempel (falls du möchtest)
	// CreatedAt time.Time
	// UpdatedAt time.Time
}
