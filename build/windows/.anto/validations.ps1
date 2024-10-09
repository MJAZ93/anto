# PowerShell script for Windows

# Function to exit script on error
function exit_on_error {
    param([string]$message)
    Write-Host $message
    exit 1
}

# Check if the anto.exe file exists
$antoExe = "anto.exe"
if (-Not (Test-Path $antoExe)) {
    exit_on_error "anto.exe not found"
}

# Call the anto.exe binary with the 'validate' parameter
Write-Host "Running anto.exe with 'validate' parameter..."
.\anto.exe validate || exit_on_error "Validation failed"

# Add any extra validation or testing script here
Write-Host "Any additional validation or testing can be added here."

Write-Host "Validation complete!"
