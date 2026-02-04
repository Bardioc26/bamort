#!/bin/bash
# Script to update, commit, and tag versions across the project

set -e

BACKEND_VERSION_FILE="backend/appsystem/version.go"
FRONTEND_VERSION_FILE="frontend/src/version.js"
FRONTEND_PACKAGE="frontend/package.json"
BACKEND_VERSION_MD="backend/VERSION.md"
FRONTEND_VERSION_MD="frontend/VERSION.md"

usage() {
    echo "Usage: $0 [-b backend_version] [-f frontend_version] [-c] [-t]"
    echo "  -b <version>   Update backend version"
    echo "  -f <version>   Update frontend version"
    echo "  -c             Commit using versions from files"
    echo "  -t             Tag using versions from files"
    echo "Examples:"
    echo "  $0 -b 0.1.31 -f 0.2.0" 
    echo "  $0 -b 0.1.31 -c -t"
    echo "  $0 -c -t"
    echo "So you can set the version at any time, commit later without worrying about commit messages and tag later when merged into main."
    exit 1
}

read_backend_version() {
    if [ ! -f "$BACKEND_VERSION_FILE" ]; then
        echo ""; return
    fi
    sed -n 's/.*const Version = "\(.*\)".*/\1/p' "$BACKEND_VERSION_FILE" | head -n1
}

read_frontend_version() {
    if [ ! -f "$FRONTEND_VERSION_FILE" ]; then
        echo ""; return
    fi
    sed -n "s/.*export const VERSION = '\(.*\)'.*/\1/p" "$FRONTEND_VERSION_FILE" | head -n1
}

BACKEND_VERSION_ARG=""
FRONTEND_VERSION_ARG=""
DO_COMMIT=false
DO_TAG=false

while [ $# -gt 0 ]; do
    case "$1" in
        -b)
            [ -z "$2" ] && usage
            BACKEND_VERSION_ARG="$2"
            shift 2
            ;;
        -f)
            [ -z "$2" ] && usage
            FRONTEND_VERSION_ARG="$2"
            shift 2
            ;;
        -c)
            DO_COMMIT=true
            shift
            ;;
        -t)
            DO_TAG=true
            shift
            ;;
        -h|--help)
            usage
            ;;
        *)
            usage
            ;;
    esac
done

if [ -z "$BACKEND_VERSION_ARG" ] && [ -z "$FRONTEND_VERSION_ARG" ] && [ "$DO_COMMIT" = false ] && [ "$DO_TAG" = false ]; then
    usage
fi

if [ -n "$BACKEND_VERSION_ARG" ]; then
    if [ -f "$BACKEND_VERSION_FILE" ]; then
        sed -i "s/const Version = \"[^\"]*\"/const Version = \"$BACKEND_VERSION_ARG\"/" "$BACKEND_VERSION_FILE"
        echo "✓ Updated $BACKEND_VERSION_FILE to $BACKEND_VERSION_ARG"
    else
        echo "⚠ Warning: $BACKEND_VERSION_FILE not found"
    fi

    if [ -f "$BACKEND_VERSION_MD" ]; then
        sed -i "s/## Current Version: .*/## Current Version: $BACKEND_VERSION_ARG/" "$BACKEND_VERSION_MD"
        echo "✓ Updated $BACKEND_VERSION_MD to $BACKEND_VERSION_ARG"
    fi
fi

if [ -n "$FRONTEND_VERSION_ARG" ]; then
    if [ -f "$FRONTEND_VERSION_FILE" ]; then
        sed -i "s/export const VERSION = '[^']*'/export const VERSION = '$FRONTEND_VERSION_ARG'/" "$FRONTEND_VERSION_FILE"
        echo "✓ Updated $FRONTEND_VERSION_FILE to $FRONTEND_VERSION_ARG"
    else
        echo "⚠ Warning: $FRONTEND_VERSION_FILE not found"
    fi

    if [ -f "$FRONTEND_PACKAGE" ]; then
        sed -i "s/\"version\": \"[^\"]*\"/\"version\": \"$FRONTEND_VERSION_ARG\"/" "$FRONTEND_PACKAGE"
        echo "✓ Updated $FRONTEND_PACKAGE to $FRONTEND_VERSION_ARG"
    else
        echo "⚠ Warning: $FRONTEND_PACKAGE not found"
    fi

    if [ -f "$FRONTEND_VERSION_MD" ]; then
        sed -i "s/## Current Version: .*/## Current Version: $FRONTEND_VERSION_ARG/" "$FRONTEND_VERSION_MD"
        echo "✓ Updated $FRONTEND_VERSION_MD to $FRONTEND_VERSION_ARG"
    fi
fi

BACKEND_VERSION_CURRENT=$(read_backend_version)
FRONTEND_VERSION_CURRENT=$(read_frontend_version)

if [ "$DO_COMMIT" = true ]; then
    if [ -z "$BACKEND_VERSION_CURRENT" ] && [ -z "$FRONTEND_VERSION_CURRENT" ]; then
        echo "❌ Cannot commit: version files missing" >&2
        exit 1
    fi

    FILES_TO_ADD=()
    [ -f "$BACKEND_VERSION_FILE" ] && FILES_TO_ADD+=("$BACKEND_VERSION_FILE")
    [ -f "$BACKEND_VERSION_MD" ] && FILES_TO_ADD+=("$BACKEND_VERSION_MD")
    [ -f "$FRONTEND_VERSION_FILE" ] && FILES_TO_ADD+=("$FRONTEND_VERSION_FILE")
    [ -f "$FRONTEND_PACKAGE" ] && FILES_TO_ADD+=("$FRONTEND_PACKAGE")
    [ -f "$FRONTEND_VERSION_MD" ] && FILES_TO_ADD+=("$FRONTEND_VERSION_MD")

    if [ ${#FILES_TO_ADD[@]} -eq 0 ]; then
        echo "❌ Cannot commit: no files to add" >&2
        exit 1
    fi

    git add "${FILES_TO_ADD[@]}"

    if [ -n "$BACKEND_VERSION_CURRENT" ] && [ -n "$FRONTEND_VERSION_CURRENT" ]; then
        if [ "$BACKEND_VERSION_CURRENT" = "$FRONTEND_VERSION_CURRENT" ]; then
            COMMIT_MSG="Bump version to $BACKEND_VERSION_CURRENT"
        else
            COMMIT_MSG="Bump backend to $BACKEND_VERSION_CURRENT, frontend to $FRONTEND_VERSION_CURRENT"
        fi
    elif [ -n "$BACKEND_VERSION_CURRENT" ]; then
        COMMIT_MSG="Bump backend to $BACKEND_VERSION_CURRENT"
    else
        COMMIT_MSG="Bump frontend to $FRONTEND_VERSION_CURRENT"
    fi

    git commit -m "$COMMIT_MSG"
    echo "✓ Committed: $COMMIT_MSG"
fi

if [ "$DO_TAG" = true ]; then
    if [ -z "$BACKEND_VERSION_CURRENT" ] && [ -z "$FRONTEND_VERSION_CURRENT" ]; then
        echo "❌ Cannot tag: version files missing" >&2
        exit 1
    fi

    if [ -n "$BACKEND_VERSION_CURRENT" ] && [ -n "$FRONTEND_VERSION_CURRENT" ] && [ "$BACKEND_VERSION_CURRENT" = "$FRONTEND_VERSION_CURRENT" ]; then
        git tag "v$BACKEND_VERSION_CURRENT"
        echo "✓ Tagged v$BACKEND_VERSION_CURRENT"
    else
        if [ -n "$BACKEND_VERSION_CURRENT" ]; then
            git tag "backend-v$BACKEND_VERSION_CURRENT"
            echo "✓ Tagged backend-v$BACKEND_VERSION_CURRENT"
        fi
        if [ -n "$FRONTEND_VERSION_CURRENT" ]; then
            git tag "frontend-v$FRONTEND_VERSION_CURRENT"
            echo "✓ Tagged frontend-v$FRONTEND_VERSION_CURRENT"
        fi
    fi
fi
