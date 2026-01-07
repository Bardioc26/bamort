# Third-Party Notices

This repository depends on third-party components. Those components remain licensed under their own terms.

## Backend (Go)

You can generate a license inventory locally, for example:

- `cd backend && go install github.com/google/go-licenses@latest`
- `cd backend && go-licenses report ./...`

Notable licenses included:
- MIT
- BSD-2-Clause / BSD-3-Clause
- Apache-2.0 (e.g. `github.com/pdfcpu/pdfcpu`, `gopkg.in/yaml.v2`)
- MPL-2.0 (e.g. `github.com/go-sql-driver/mysql`)

## Frontend (Node/Vue)

You can generate a license inventory locally, for example:

- `cd frontend && npm install`
- `cd frontend && npx license-checker --production --summary`

Notable licenses included:
- MIT
- BSD-2-Clause / BSD-3-Clause
- ISC

## Icons

- `frontend/src/components/icons/IconTooling.vue` contains an icon from MaterialDesignIcons (Templarian), licensed under Apache-2.0 (see the file header comment for the upstream link).

## Artwork / images

The repository contains image assets (e.g. under `frontend/public/` and template assets under `backend/templates/`).

Ensure you have redistribution rights for all included artwork (and add explicit attribution/license information if any files originate from third parties).
