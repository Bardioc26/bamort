// Beispiel für Anfragen an die API für die Frontend-Entwicklung

/*
1. Anfrage zum Erlernen einer neuen Fertigkeit:
*/
fetch('/api/characters/123/skill-cost', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    name: 'Schwimmen',
    type: 'skill',
    action: 'learn'
  }),
})
.then(response => response.json())
.then(data => console.log(data))
.catch(error => console.error('Error:', error));

/*
2. Anfrage zum Verbessern einer vorhandenen Fertigkeit:
*/
fetch('/api/characters/123/skill-cost', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    name: 'Schwimmen',
    current_level: 12,
    type: 'skill',
    action: 'improve'
  }),
})
.then(response => response.json())
.then(data => console.log(data))
.catch(error => console.error('Error:', error));

/*
3. Anfrage zum Erlernen eines neuen Zaubers:
*/
fetch('/api/characters/123/skill-cost', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    name: 'Feuerball',
    type: 'spell',
    action: 'learn'
  }),
})
.then(response => response.json())
.then(data => console.log(data))
.catch(error => console.error('Error:', error));
