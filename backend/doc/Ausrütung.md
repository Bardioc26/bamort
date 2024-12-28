# API Examples
## Create Ausruestung

POST /ausruestung

{
  "character_id": 1,
  "name": "Magic Sword",
  "anzahl": 1,
  "gewicht": 2.5,
  "wert": 150.0,
  "beinhaltet_in": null,
  "beschreibung": "A glowing sword imbued with magical energy.",
  "bonus": 5,
  "ist_magisch": true,
  "abw": 10,
  "ausgebrannt": false
}

## Get All Ausruestung for a Character

GET /ausruestung/1 Response:
[
    {
      "ausruestung_id": 1,
      "character_id": 1,
      "name": "Magic Sword",
      "anzahl": 1,
      "gewicht": 2.5,
      "wert": 150.0,
      "beinhaltet_in": null,
      "beschreibung": "A glowing sword imbued with magical energy.",
      "bonus": 5,
      "ist_magisch": true,
      "abw": 10,
      "ausgebrannt": false
    }
  ]
  
## Update Ausruestung

PUT /ausruestung/1

{
    "anzahl": 2,
    "gewicht": 3.0
  }
  
## Delete Ausruestung

DELETE /ausruestung/1

## Summary

These endpoints provide complete CRUD functionality for managing Ausruestung items. They integrate with the database and offer flexibility to interact with equipment data.