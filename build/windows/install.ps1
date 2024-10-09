# PowerShell script for Windows

# Output message
Write-Host "Running anto with 'init' parameter on Windows..."

# Make the 'anto' file executable (no need for chmod in Windows)
# Run the 'anto' executable with the 'init' parameter
if (Test-Path ".\anto.exe") {
    Write-Host "Executing anto.exe with 'init' parameter..."
    .\anto.exe init
} else {
    Write-Host "anto.exe not found, ensure it's in the current directory."
    exit 1
}

# Output completion message
Write-Host "Installation complete!"
