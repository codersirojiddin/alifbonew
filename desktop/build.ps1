$ErrorActionPreference = 'Stop'
$root = $PSScriptRoot
$candidates = @(
  "$env:WINDIR\Microsoft.NET\Framework64\v4.0.30319\csc.exe",
  "$env:WINDIR\Microsoft.NET\Framework\v4.0.30319\csc.exe"
)
$compiler = $candidates | Where-Object { Test-Path $_ } | Select-Object -First 1
if (!$compiler) { throw '.NET Framework 4.8 C# compiler was not found. Install the .NET Framework 4.8 Developer Pack.' }
$output = Join-Path $root 'out'
New-Item $output -ItemType Directory -Force | Out-Null
$app = Join-Path $output 'Harf.exe'
& $compiler /nologo /target:winexe /platform:anycpu /optimize+ "/out:$app" /reference:System.Windows.Forms.dll /reference:System.Drawing.dll (Join-Path $root 'Harf.cs')
if ($LASTEXITCODE -ne 0) { throw 'C# compilation failed.' }
$test = Start-Process -FilePath $app -ArgumentList '--self-test' -Wait -PassThru
if ($test.ExitCode -ne 0) { throw 'Desktop converter tests failed.' }
Copy-Item $app (Join-Path $output 'Harf-Setup.exe') -Force
Get-FileHash (Join-Path $output 'Harf-Setup.exe') -Algorithm SHA256
Write-Host 'Build complete. Test Harf.exe in Windows before distributing Harf-Setup.exe.'
