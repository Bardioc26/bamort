# Frontend Version Management

## Current Version: 0.1.22

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

## Git Commit Information

The git commit is injected via environment variable:
- Set `VITE_GIT_COMMIT` in `.env` or at build time
- Falls back to "unknown" if not set

Example `.env`:
```bash
VITE_GIT_COMMIT=d0c177b
```

## Usage in Components

```javascript
import { getVersion, getGitCommit, getVersionInfo } from '@/version'

// Get version string
const version = getVersion()  // "0.1.30"

// Get git commit
const commit = getGitCommit()  // "d0c177b" or "unknown"

// Get full info object
const info = getVersionInfo()  // { version: "0.1.30", gitCommit: "d0c177b" }
```

## Landing Page Display

The landing page shows both:
- **Frontend Version**: From `/frontend/src/version.js`
- **Backend Version**: Fetched from `/api/public/version`

This allows users to see if frontend and backend are in sync.

## Build-time Version Injection

To inject git commit at build time, update `vite.config.js`:

```javascript
import { defineConfig } from 'vite'
import { execSync } from 'child_process'

const gitCommit = execSync('git rev-parse --short HEAD').toString().trim()

export default defineConfig({
  define: {
    'import.meta.env.VITE_GIT_COMMIT': JSON.stringify(gitCommit)
  }
})
```
