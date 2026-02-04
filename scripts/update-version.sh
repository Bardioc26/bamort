#!/bin/bash
# Script to update version number across the project

set -e

if [ -z "$1" ]; then
    echo "Usage: $0 <backend_version> [frontend_version] [auto]"
    echo "Example: $0 0.1.31              # Sets both to 0.1.31"
    echo "Example: $0 0.1.31 0.2.0        # Sets different versions"
    echo "Example: $0 0.1.31 auto         # Sets both to 0.1.31 and auto-commits"
    echo "Example: $0 0.1.31 0.2.0 auto   # Sets different versions and auto-commits"
    exit 1
fi

# Parse arguments
BACKEND_VERSION="$1"
AUTO_COMMIT=false

# Check if second argument is "auto"
if [ "$2" = "auto" ]; then
    FRONTEND_VERSION="$1"
    AUTO_COMMIT=true
elif [ -z "$2" ]; then
    FRONTEND_VERSION="$1"
elif [ "$3" = "auto" ]; then
    FRONTEND_VERSION="$2"
    AUTO_COMMIT=true
else
    FRONTEND_VERSION="$2"
fi


if [ "$BACKEND_VERSION" = "$FRONTEND_VERSION" ]; then
    echo "Updating both backend and frontend to version $BACKEND_VERSION..."
else
    echo "Updating backend to $BACKEND_VERSION and frontend to $FRONTEND_VERSION..."
fi

# Update backend version
BACKEND_VERSION_FILE="backend/appsystem/version.go"
if [ -f "$BACKEND_VERSION_FILE" ]; then
    sed -i "s/const Version = \"[^\"]*\"/const Version = \"$BACKEND_VERSION\"/" "$BACKEND_VERSION_FILE"
    echo "✓ Updated $BACKEND_VERSION_FILE to $BACKEND_VERSION"
else
    echo "⚠ Warning: $BACKEND_VERSION_FILE not found"
fi

# Update frontend version
FRONTEND_VERSION_FILE="frontend/src/version.js"
if [ -f "$FRONTEND_VERSION_FILE" ]; then
    sed -i "s/export const VERSION = '[^']*'/export const VERSION = '$FRONTEND_VERSION'/" "$FRONTEND_VERSION_FILE"
    echo "✓ Updated $FRONTEND_VERSION_FILE to $FRONTEND_VERSION"
else
    echo "⚠ Warning: $FRONTEND_VERSION_FILE not found"
fi

# Update frontend package.json version
FRONTEND_PACKAGE="frontend/package.json"
if [ -f "$FRONTEND_PACKAGE" ]; then
    sed -i "s/\"version\": \"[^\"]*\"/\"version\": \"$FRONTEND_VERSION\"/" "$FRONTEND_PACKAGE"
    echo "✓ Updated $FRONTEND_PACKAGE to $FRONTEND_VERSION"
else
    echo "⚠ Warning: $FRONTEND_PACKAGE not found"
fi

# Update VERSION.md files
BACKEND_VERSION_MD="backend/VERSION.md"
if [ -f "$BACKEND_VERSION_MD" ]; then
    sed -i "s/## Current Version: .*/## Current Version: $BACKEND_VERSION/" "$BACKEND_VERSION_MD"
    echo "✓ Updated $BACKEND_VERSION_MD to $BACKEND_VERSION"
fi

FRONTEND_VERSION_MD="frontend/VERSION.md"
if [ -f "$FRONTEND_VERSION_MD" ]; then
    sed -i "s/## Current Version: .*/## Current Version: $FRONTEND_VERSION/" "$FRONTEND_VERSION_MD"
    echo "✓ Updated $FRONTEND_VERSION_MD to $FRONTEND_VERSION"
fi

echo ""
echo "✅ Version update complete!"
echo "   Backend:  $BACKEND_VERSION"
echo "   Frontend: $FRONTEND_VERSION"
echo ""

# Auto-commit if requested
if [ "$AUTO_COMMIT" = true ]; then
    echo "Auto-committing changes..."
    if [ "$BACKEND_VERSION" = "$FRONTEND_VERSION" ]; then
        git add backend/appsystem/version.go frontend/src/version.js frontend/package.json frontend/VERSION.md
        git commit -m "Bump version to $BACKEND_VERSION"
        git tag v$BACKEND_VERSION
        echo "✓ Committed and tagged as v$BACKEND_VERSION"
    else
        git add backend/config/version.go frontend/src/version.js frontend/package.json frontend/VERSION.md
        git commit -m "Bump backend to $BACKEND_VERSION, frontend to $FRONTEND_VERSION"
        git tag backend-v$BACKEND_VERSION 
        git tag frontend-v$FRONTEND_VERSION
        git tag v$BACKEND_VERSION -m "Backend version $BACKEND_VERSION, Frontend version $FRONTEND_VERSION"
        echo "✓ Committed and tagged as backend-v$BACKEND_VERSION, frontend-v$FRONTEND_VERSION, v$BACKEND_VERSION"
    fi
    echo ""
    echo "Next step:"
    echo "  Push: git push && git push --tags"
else
    echo "Next steps:"
    echo "1. Review changes: git diff"
    if [ "$BACKEND_VERSION" = "$FRONTEND_VERSION" ]; then
        echo "2. Commit changes: git commit -am 'Bump version to $BACKEND_VERSION'"
        echo "3. Tag release: git tag v$BACKEND_VERSION"
        git tag v$BACKEND_VERSION
    else
        echo "2. Commit changes: git commit -am 'Bump backend to $BACKEND_VERSION, frontend to $FRONTEND_VERSION'"
        echo "3. Tag releases: git tag backend-v$BACKEND_VERSION && git tag frontend-v$FRONTEND_VERSION"
        git tag backend-v$BACKEND_VERSION 
        git tag frontend-v$FRONTEND_VERSION
        git tag v$BACKEND_VERSION -m "Backend version $BACKEND_VERSION, Frontend version $FRONTEND_VERSION"
    fi
    echo "4. Push: git push && git push --tags"
fi
