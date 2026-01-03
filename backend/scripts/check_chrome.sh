#!/bin/bash
# check_chrome.sh - Verify Chrome/Chromium is available for PDF export

echo "=== Chrome/Chromium Availability Check ==="
echo ""

# Check for Chrome in common locations
CHROME_PATHS=(
    "/usr/bin/google-chrome"
    "/usr/bin/google-chrome-stable"
    "/usr/bin/chromium"
    "/usr/bin/chromium-browser"
    "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
    "/snap/bin/chromium"
)

FOUND=false

for path in "${CHROME_PATHS[@]}"; do
    if [ -x "$path" ]; then
        echo "✓ Found Chrome at: $path"
        VERSION=$("$path" --version 2>/dev/null || echo "Unknown")
        echo "  Version: $VERSION"
        FOUND=true
        
        # Suggest setting CHROME_BIN if not already set
        if [ -z "$CHROME_BIN" ]; then
            echo "  Suggestion: export CHROME_BIN=\"$path\""
        fi
        echo ""
    fi
done

# Check PATH
echo "Checking PATH for chrome/chromium..."
if command -v google-chrome &> /dev/null; then
    echo "✓ 'google-chrome' found in PATH"
    google-chrome --version
    FOUND=true
elif command -v chromium-browser &> /dev/null; then
    echo "✓ 'chromium-browser' found in PATH"
    chromium-browser --version
    FOUND=true
elif command -v chromium &> /dev/null; then
    echo "✓ 'chromium' found in PATH"
    chromium --version
    FOUND=true
else
    echo "✗ No chrome/chromium found in PATH"
fi

echo ""
echo "Current CHROME_BIN: ${CHROME_BIN:-not set}"
echo ""

if [ "$FOUND" = true ]; then
    echo "✓ Chrome/Chromium is available - PDF export should work"
    exit 0
else
    echo "✗ Chrome/Chromium NOT found - PDF export will FAIL"
    echo ""
    echo "Please install Chrome or Chromium:"
    echo "  Debian/Ubuntu: sudo apt-get install chromium-browser"
    echo "  Alpine: apk add chromium"
    echo "  macOS: brew install --cask google-chrome"
    echo ""
    echo "Or set CHROME_BIN to point to your Chrome installation"
    exit 1
fi
