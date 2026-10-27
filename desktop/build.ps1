$ErrorActionPreference = 'Stop'
Push-Location $PSScriptRoot
try {
  $env:CGO_ENABLED = '0'
  $env:GOOS = 'windows'
  $env:GOARCH = 'amd64'
  go test ./...
  if ($LASTEXITCODE -ne 0) { throw 'Go tests failed.' }
  go vet -unsafeptr=false ./...
  if ($LASTEXITCODE -ne 0) { throw 'Go vet failed.' }
  New-Item out -ItemType Directory -Force | Out-Null
  go build -trimpath -ldflags '-H=windowsgui -s -w' -o out/Harf.exe ./cmd/harf
  if ($LASTEXITCODE -ne 0) { throw 'Windows Go build failed.' }
  Copy-Item out/Harf.exe out/Harf-Setup.exe -Force
  Get-FileHash out/Harf-Setup.exe -Algorithm SHA256
  Write-Host 'Go build complete: out/Harf.exe and out/Harf-Setup.exe'
} finally {
  Pop-Location
}
