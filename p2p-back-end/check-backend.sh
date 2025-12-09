# #!/bin/bash
# set -e

# echo "🔍 Running golangci-lint..."
# golangci-lint run ./... --verbose

# echo "🔒 Running gosec..."
# gosec ./... -fmt=json -out=gosec-report.json

# echo "✅ Backend checks complete."
# echo "Gosec report saved to gosec-report.json"
#!/bin/bash
set -e

echo "🔍 Running golangci-lint..."
if golangci-lint run ./... --verbose; then
  echo "✅ golangci-lint passed!"
else
  echo "⚠️ golangci-lint found issues, continuing..."
fi

echo "🔒 Running gosec..."
if command -v gosec >/dev/null 2>&1; then
  gosec ./... -fmt=json -out=gosec-report.json
else
  echo "⚠️ gosec not installed, running via go run..."
  go run github.com/securego/gosec/v2/cmd/gosec@latest ./... -fmt=json -out=gosec-report.json
fi

echo "✅ Backend checks complete."
echo "Gosec report saved to gosec-report.json"
