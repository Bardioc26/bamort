# Plan a new Feature
The Application backend is written in GoLang the frontend in Vue3.
The packages "transfero" and "importero" are import/export formats you will NOT use exept for the definition of the importero.CharacterImport data struchture. 

Make a plan for implementing a new feature. Take the following things into account:
- the feature is to inport an export data.
- the data is all about character data
    - descriptiv data
        - male, female
        - character class
        - size, age, social class, ...
        - ...
    - stats
        - St, Gw, In,...
        - ...
    - skills
    - spell
    - equipment
        - weapons
        - containers
        - transportation equipment
        - ...    
    - character history
        - learned skill and spells (date, costs, level)
        - gained expirience points or gold(date ,amount)
        - ...
- the import data can come from different systems and contain da data in different formats
- the import source may require tailored conversion into the common internal format (we refer to it as BMRT-Format)
- if all information is transformed into the BMRT-Format we want to check if the Skills, Spells, Equipment and other things are available in the Masterdata of the game system set by the import data.
    - link items with exisiting master data
    - add new items to masterdata 
    - OR mark item as personal Item (players house rules)
- later I want to add new formats 
    - new formats should be plugable. 
        - At best by droping a file into the right folder
        - Or by adding the new package and registering it on one location
- the export of the data must create ecactly the same output as we had in the input file

# Chat GPT thoughts

Got it — I’ll treat the document as the *baseline spec* and make the feature **complete, extensible, and implementation-ready**, while staying system-agnostic.

Below is a **comprehensive proposal** that:

* expands the feature scope,
* suggests concrete implementation techniques,
* and evaluates best practices for maintainability, correctness, and extensibility.

All assumptions are derived from your uploaded plan .

---

## 1. Clarify & Complete the Feature Scope

### 1.1 Core Goal (Refined)

Build a **lossless, bidirectional import/export system** for character data that:

* accepts multiple external formats,
* converts them into a **canonical internal representation (BMRT-Format)**,
* validates and reconciles data against master data,
* and exports data **logical equivalent** to the original input format.

This implies:

* **Round-trip integrity** is a hard requirement.
* Import and export are *two sides of the same adapter*.

---

## 2. Canonical Model: BMRT-Format (Foundation)

### 2.1 BMRT-Format Design Principles

Best practices:

* **Explicit schema** (JSON Schema / Protobuf / Avro)
* **Versioned** (e.g. `bmrtVersion: 1.0`)
* **Loss-tolerant but not lossy**

Key techniques:

* Store *raw source metadata* alongside normalized data
* Preserve unknown fields instead of discarding them

### 2.2 Suggested BMRT Structure (High Level)

```text
BMRTCharacter
├── Meta
│   ├── sourceSystem
│   ├── sourceFormat
│   ├── importTimestamp
│   ├── bmrtVersion
│   └── rawSourceSnapshot
│
├── Identity
│   ├── name
│   ├── gender
│   ├── species
│   ├── class
│   ├── age
│   ├── size
│   └── socialClass
│
├── Stats
│   ├── primaryStats (ST, GW, IN, ...)
│   ├── derivedStats
│   └── modifiers
│
├── Skills
├── Spells
├── Equipment
│   ├── weapons
│   ├── containers
│   └── transportation
│
├── History
│   ├── skillLearningEvents
│   ├── spellLearningEvents
│   ├── xpEvents
│   └── currencyEvents
│
└── Extensions
    └── formatSpecificData
```

**Why this matters**
This structure:

* isolates system-specific quirks,
* supports partial imports,
* and enables strict validation without blocking creativity.

---

## 3. Import Pipeline Architecture

### 3.1 Multi-Stage Import Pipeline (Recommended)

```text
Input File
  ↓
Format Detector
  ↓
Format Adapter (Plugin)
  ↓
Intermediate DTO
  ↓
BMRT Transformer
  ↓
Validation & Reconciliation
  ↓
BMRT Character
```

### 3.2 Format Detection

Techniques:

* File extension + MIME type
* Magic bytes
* Optional user override

Best practice:

* Detection should be **pure and side-effect free**

---

## 4. Pluggable Format System (Critical Requirement)

### 4.1 Plugin Contract (Strongly Recommended)

Each format plugin should implement:

```text
ICharacterFormatPlugin
├── canHandle(input): boolean
├── parse(input): ParsedCharacter
├── toBMRT(parsed): BMRTCharacter
├── fromBMRT(bmrt): ParsedCharacter
├── serialize(parsed): OutputFile
└── getMetadata()
```

### 4.2 Plugin Discovery Strategies

## 4.2 Plugin Discovery Strategies 

### Core Constraint Shift (Important)

Unlike JVM or Node ecosystems, **Go does not support safe, dynamic runtime loading of arbitrary plugins across platforms** in a production-friendly way.

Therefore:

* `.so` Go plugins are **not recommended** for extensibility
* Process isolation and API contracts are preferred
* Microservices become a *best practice*, not an over-architecture

---

## 4.2.1 Recommended Architecture: Adapter Microservices

### High-Level Idea

Each import/export format is implemented as an **independent adapter service** that communicates with the central application via a **stable API contract**.

```text
Vue3 Frontend
     ↓
Central Go Backend
     ↓
Format Adapter Service(s)
     ↓
BMRT-Format
```

---

## 4.2.2 Adapter Service Contract (Canonical)

Each adapter service must expose a minimal, strict API:

```text
POST /detect
POST /import
POST /export
GET  /metadata
```

### Endpoint Responsibilities

#### `GET /metadata`

```json
{
  "id": "foundry-vtt",
  "name": "Foundry VTT Character Format",
  "supportedExtensions": [".json"],
  "supportedVersions": ["10.x", "11.x"],
  "bmrtVersion": "1.0",
  "capabilities": ["import", "export"]
}
```

#### `POST /detect`

* Input: raw file or metadata
* Output: confidence score (`0.0 – 1.0`)

#### `POST /import`

* Input: raw source file
* Output: `BMRTCharacter + SourceMapping`

#### `POST /export`

* Input: `BMRTCharacter + SourceMapping`
* Output: serialized file **logical equivalent** to original format

---

## 4.2.3 Why Microservices Are the Best Fit for Go

### Advantages

✔ No dynamic linking issues
✔ Language-agnostic (adapters can be Go, Node, Python, etc.)
✔ Crash isolation (bad plugin ≠ downed backend)
✔ Clear versioning and lifecycle
✔ Easy to add/remove adapters without redeploying core app

This aligns well with Go’s philosophy:

> “Communicate by sharing memory? No. Share memory by communicating.”

---

## 4.2.4 Discovery Strategies

### Strategy A — Service Registry (Recommended)

Adapter services register themselves at startup.

```text
Central App
└── Adapter Registry
    ├── foundry-vtt
    ├── roll20
    └── custom-house-rules
```

**Implementation options:**

* Static config (YAML / ENV)
* Service discovery (Consul, etcd)
* Kubernetes labels + discovery

**Best practice**:

* Central app polls `/metadata`
* Registry is cached and versioned

---

### Strategy B — Drop-in via Container Deployment (Closest to “Drop a File”)

Instead of dropping a file:

* Drop a **container** (Docker image)
* Add it to a watched directory or compose file

```text
/adapters
 ├── adapter-foundry/
 │   └── docker-compose.yaml
 └── adapter-roll20/
```

This preserves your original intent:

> *“At best by dropping a file into the right folder”*

Just modernized and production-safe.

---

### Strategy C — Embedded Adapters (Fallback / Core Formats)

For **core or official formats only**:

* Compile adapters directly into the Go backend
* Implement behind the same interface as remote adapters

```go
type FormatAdapter interface {
    Detect(input []byte) float64
    Import(input []byte) (*BMRTCharacter, error)
    Export(bmrt *BMRTCharacter) ([]byte, error)
    Metadata() AdapterMetadata
}
```

**Important rule**:

> Embedded adapters must behave exactly like remote ones.

---

## 4.2.5 Adapter Selection Flow

```text
Upload File
   ↓
Ask all adapters: /detect
   ↓
Select highest confidence
   ↓
Import via selected adapter
   ↓
Store adapter ID + mapping
```

This ensures:

* Deterministic exports
* Correct adapter reuse later
* Multiple adapters can support the same format version

---

## 4.2.6 Frontend (Vue3) Considerations

The Vue3 frontend:

* Never talks to adapters directly
* Talks only to the central backend
* Receives adapter metadata for UI decisions

Example:

* “Detected format: Foundry VTT (confidence 92%)”
* “3 unresolved items need mapping”

This keeps:
✔ frontend simple
✔ security tight
✔ adapters replaceable

---

## 4.2.7 Security & Stability Best Practices

* Adapters run in isolated processes or containers
* Timeouts on all adapter calls
* Strict payload size limits
* Schema validation on responses
* Adapter version pinning per import

---

## 4.2.8 Summary of the Strategy

| Concern                 | Chosen Solution            |
| ----------------------- | -------------------------- |
| Runtime extensibility   | Microservice adapters      |
| Go compatibility        | API-based contracts        |
| Plugin safety           | Process isolation          |
| “Drop-in” extensibility | Container-based deployment |
| Round-trip fidelity     | Adapter-bound export       |
| Future formats          | Zero core redeploy         |

---

### Net Result

You get:

* **True pluggability**
* **Production-safe extensibility**
* **Language freedom**
* **Perfect import/export symmetry**
* **A clean Go + Vue3 architecture**

---

## 5. Master Data Reconciliation

### 5.1 Reconciliation Strategy

When importing Skills, Spells, Equipment:

1. **Exact match** → link to master data
2. **Fuzzy match** → suggest mapping
3. **No match**:

   * add to master data **OR**
   * mark as *Personal / House Rule*

### 5.2 Suggested Techniques

* Name + category + system matching
* Optional synonym dictionaries

### 5.3 Conflict Resolution UX (Even if not UI yet)

* Deterministic defaults
* Optional user intervention hooks
* Audit trail for decisions

---

## 6. Validation & Integrity

### 6.1 Validation Layers

* Schema validation (BMRT)
* Business rules (e.g. stat caps)
* Referential integrity (master data links)

### 6.2 Best Practices

* Validation should be **configurable**
* Warnings ≠ errors
* Never block import unless data is unrecoverable

---

## 7. Export: Guaranteed Round-Trip Fidelity

### 7.1 Export Requirements

* Same format
* Same structure
* Same values
* Same ordering (if applicable)
* Same optional fields

### 7.2 Key Technique: Source Snapshot

Store:

* original raw file
* mapping metadata
* field-level transformations

This allows:

* perfect regeneration
* future re-export even if BMRT evolves

---

## 8. Versioning & Evolution

### 8.1 Format Versioning

* Each plugin declares supported versions
* BMRT migrators handle internal evolution

### 8.2 Backward Compatibility

* Never delete fields
* Deprecate instead
* Provide migration utilities

---

## 9. Testing Strategy (Non-Optional)

### 9.1 Required Test Types

* Golden file tests (import → export equality)
* Cross-format consistency tests
* Plugin isolation tests
* Master data reconciliation tests

### 9.2 Best Practice

Every plugin ships with:

* example files
* expected BMRT output
* expected re-export output

---

## 10. Non-Functional Requirements

### Performance

* Streaming parsers for large files
* Lazy loading for history

### Security

* Sandbox plugins
* Validate inputs aggressively
* No arbitrary code execution

### Observability

* Structured logs per import step
* Import reports
* Debug mode with transformation diffs

---

## 11. What This Adds Beyond the Original Plan

✔ Explicit round-trip guarantee
✔ Plugin contracts and lifecycle
✔ Validation & reconciliation strategy
✔ Master data conflict handling
✔ Versioning & future-proofing
✔ Testing and operability best practices

---

