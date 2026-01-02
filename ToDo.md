# Frontend 

* verbessere Versionserstellung

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
* API Dokumentation
* ./testdata  neu erstellen und aktuell halten
* in jedem Package eine README.md erstellen in der kurz erklärt wird wozu das package dient, welche Abhängigkeiten bestehen, wie es zu benutzen ist und wie die tests funktionieren.
* verbessere Versionserstellung
* Waffenfertigkeiten haben keine Bonuseigenschaft.
* waffenfertigkeiten müssen in andere Katagorien eingeteilt werden. Nahkampf, Schusswaffen, Verteidigungswaffen etc.

## Refaktor

* Export Import Module neu grupieren
* export_temp sollte nicht im backend liegen
* routing verbessern
* templates für PDF export sollten nicht direkt im backend liegen

### refactor multi game system

* Datenstruktur gsm, character euipment,... nach Spielsystem gruppieren
* Middleware erstellen
* Regelmechanik 


# Deployment 

* depoyment verbessern, Konfiguration ausschließlich über .env file oder datenbank