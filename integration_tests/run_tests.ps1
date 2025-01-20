# Kill any existing process on port 8088
$existingProcess = Get-NetTCPConnection -LocalPort 8088 -ErrorAction SilentlyContinue
if ($existingProcess) {
    Stop-Process -Id $existingProcess.OwningProcess -Force
}

# Start the server in background
Push-Location examples\basic
$server = Start-Process -FilePath "go" -ArgumentList "run", "main.go" -PassThru -WindowStyle Hidden

# Function to cleanup
function Cleanup {
    Write-Host "Cleaning up..."
    if ($server) {
        Stop-Process -Id $server.Id -Force -ErrorAction SilentlyContinue
    }
    Pop-Location
}

# Setup cleanup on script exit
trap {
    Cleanup
    break
}

# Wait for server to be ready
Write-Host "Waiting for server to start..."
$ready = $false
for ($i = 1; $i -le 30; $i++) {
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:8088/std" -UseBasicParsing -ErrorAction SilentlyContinue
        if ($response.StatusCode -eq 200) {
            Write-Host "Server is ready!"
            $ready = $true
            break
        }
    }
    catch {
        Start-Sleep -Seconds 1
    }
}

if (-not $ready) {
    Write-Host "Server failed to start within 30 seconds"
    Cleanup
    exit 1
}

# Run the tests
Set-Location ..\..\integration_tests
hurl --test basic.hurl

# Cleanup
Cleanup 