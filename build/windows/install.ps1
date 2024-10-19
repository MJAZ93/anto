# PowerShell script for Windows

# Output message
Write-Host "Running anto with 'init' parameter on Windows..."

# Check if the __MACOSX folder exists and remove it if it does
if (Test-Path "__MACOSX") {
    Write-Host "Removing __MACOSX..."
    Remove-Item -Recurse -Force "__MACOSX"
}

# Navigate to the extracted folder
Set-Location ".anto" -ErrorAction Stop

# Windows - Run the 'anto' binary with 'init' parameter
Write-Host "Running anto with 'init' parameter on Windows..."
# No need to set execute permission on Windows, just run the executable
Start-Process -NoNewWindow -Wait "./anto.exe" "init"

Write-Host "Installation complete!"
