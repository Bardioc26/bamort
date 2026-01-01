# Frontend 
## Refaktor

* .github/instructions/vuejs3.instructions.md und .github/instructions/vue.instructions.md zusammen führen
* Vue von "Options API" auf "Composition API" umstellen wegen besserer Trennung von Visuals und Struktur

## Styling

* CSS soweit wie möglich aus den Komponenten und Views entfernen und für einheitliches Styling in base.css und main.css zusammen führen
* Styling der Komponenten vereinheitlichen
* CSS soweit wie möglich vereinfachen

# Backend

* PDF export fixen
* datensynchronisation DB Prod -> prepared_test_data
* strukturieren Datenexport/import verbessern
    * export import grupieren nach
        * Programmkonfigurationen
        * Daten
        * Regeln
        * Charakterdaten
        * Userdaten

## Refaktor

* Export Import Module neu grupieren

### refactor multi game system

* Datenstruktur gsm, character euipment,... nach Spielsystem gruppieren
* Middleware erstellen
* Regelmechanik 


# Deployment 

* depoyment verbessern, Konfiguration ausschließlich über .env file oder datenbank