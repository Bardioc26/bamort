# Frontend Version Management

## Current Version: 0.1.29

The frontend version is managed independently from the backend.

## Version Locations

1. **Primary source**: `/frontend/src/version.js`
   - Contains the VERSION constant
   - Exports version info functions

2. **Package metadata**: `/frontend/package.json`
   - Standard npm version field
   - Should match version.js

## Updating the Version

### Option 1: Using the update script (Recommended)
```bash
# Updates both backend and frontend
./scripts/update-version.sh 0.1.31
```

### Option 2: Manual update
Edit `/frontend/src/version.js`:
```javascript
export const VERSION = '0.1.31'  // Update this
```

And `/frontend/package.json`:
```json
{
  "version": "0.1.31"  // Update this
}
```

## Usage in Components

```javascript
import { getVersion, getVersionInfo } from '@/version'

// Get version string
const version = getVersion()  // "0.1.30"

// Get full info object
const info = getVersionInfo()  // { version: "0.1.30" }
```

## Landing Page Display

The landing page shows both:
- **Frontend Version**: From `/frontend/src/version.js`
- **Backend Version**: Fetched from `/api/public/version`

This allows users to see if frontend and backend are in sync.
