# Derived Values API Dokumentation

Das neue System für abgeleitete Werte teilt die Berechnungen in zwei Kategorien auf:

## 1. Statische Felder (ohne Würfelwürfe)
**Endpunkt:** `POST /api/characters/calculate-static-fields`

Berechnet alle Felder, die keine Würfelwürfe benötigen.

### Request:
```json
{
  "st": 70,
  "gs": 60,
  "gw": 65,
  "ko": 75,
  "in": 50,
  "zt": 30,
  "au": 55,
  "rasse": "Menschen",
  "typ": "Krieger"
}
```

### Response:
```json
{
  "ausdauer_bonus": 10,
  "schadens_bonus": 2,
  "angriffs_bonus": 0,
  "abwehr_bonus": 0,
  "zauber_bonus": -1,
  "resistenz_bonus_koerper": 1,
  "resistenz_bonus_geist": 0,
  "resistenz_koerper": 12,
  "resistenz_geist": 11,
  "abwehr": 11,
  "zaubern": 10,
  "raufen": 6
}
```

## 2. Würfelfelder (mit Würfelwürfen)
**Endpunkt:** `POST /api/characters/calculate-rolled-field`

Berechnet einzelne Felder mit Würfelwürfen, die vom Frontend bereitgestellt werden.

### PA (Persönliche Ausstrahlung):
```json
{
  "st": 70, "gs": 60, "gw": 65, "ko": 75, "in": 50, "zt": 30, "au": 55,
  "rasse": "Menschen",
  "typ": "Krieger",
  "field": "pa",
  "roll": 55
}
```

Response:
```json
{
  "field": "pa",
  "value": 55,
  "formula": "1d100 + 4×In/10 - 20",
  "details": {
    "roll": 55,
    "in_bonus": 20,
    "base_modifier": -20,
    "modifier": 0
  }
}
```

### WK (Willenskraft):
```json
{
  "st": 70, "gs": 60, "gw": 65, "ko": 75, "in": 50, "zt": 30, "au": 55,
  "rasse": "Menschen",
  "typ": "Krieger",
  "field": "wk",
  "roll": 45
}
```

### LP Max (Lebenspunkte Maximum):
```json
{
  "st": 70, "gs": 60, "gw": 65, "ko": 75, "in": 50, "zt": 30, "au": 55,
  "rasse": "Menschen",
  "typ": "Krieger",
  "field": "lp_max",
  "roll": 2
}
```

### AP Max (Abenteuerpunkte Maximum):
```json
{
  "st": 70, "gs": 60, "gw": 65, "ko": 75, "in": 50, "zt": 30, "au": 55,
  "rasse": "Menschen",
  "typ": "Krieger",
  "field": "ap_max",
  "roll": 3
}
```

### B Max (Bewegungsweite):
```json
{
  "st": 70, "gs": 60, "gw": 65, "ko": 75, "in": 50, "zt": 30, "au": 55,
  "rasse": "Menschen",
  "typ": "Krieger",
  "field": "b_max",
  "roll": [2, 1, 3, 2]
}
```

## Verwendung im Frontend

### 1. Alle statischen Felder auf einmal berechnen:
```javascript
const calculateStaticFields = async (attributes, race, type) => {
  const response = await fetch('/api/characters/calculate-static-fields', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      st: attributes.st,
      gs: attributes.gs,
      gw: attributes.gw,
      ko: attributes.ko,
      in: attributes.in,
      zt: attributes.zt,
      au: attributes.au,
      rasse: race,
      typ: type
    })
  });
  return await response.json();
};
```

### 2. Einzelne Würfelfelder berechnen:
```javascript
const calculateRolledField = async (attributes, race, type, field, diceRoll) => {
  const response = await fetch('/api/characters/calculate-rolled-field', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      ...attributes,
      rasse: race,
      typ: type,
      field: field,
      roll: diceRoll
    })
  });
  return await response.json();
};

// Beispiele:
// PA: roll = einzelner Würfelwert (1-100)
// WK: roll = einzelner Würfelwert (1-100)  
// LP Max: roll = einzelner Würfelwert (1-3)
// AP Max: roll = einzelner Würfelwert (1-3)
// B Max: roll = Array von Würfelwerten (je nach Rasse 2-4 Werte von 1-3)
```

## Würfelwürfe im Frontend generieren

```javascript
const rollD100 = () => Math.floor(Math.random() * 100) + 1;
const rollD3 = () => Math.floor(Math.random() * 3) + 1;

const getMovementDiceCount = (race) => {
  switch(race) {
    case 'Gnome':
    case 'Halblinge':
      return 2;
    case 'Zwerge':
      return 3;
    default:
      return 4;
  }
};

const rollMovement = (race) => {
  const diceCount = getMovementDiceCount(race);
  return Array.from({length: diceCount}, () => rollD3());
};
```
