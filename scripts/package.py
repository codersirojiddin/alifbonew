from pathlib import Path
from zipfile import ZipFile, ZIP_DEFLATED

root = Path(__file__).resolve().parents[1]
output = root / 'dist/downloads/Harf-Windows-source.zip'
output.parent.mkdir(parents=True, exist_ok=True)
with ZipFile(output, 'w', ZIP_DEFLATED) as archive:
    for name in ('Harf.cs', 'build.ps1', 'Build.cmd', 'README.md'):
        archive.write(root / 'desktop' / name, name)
print(output)
