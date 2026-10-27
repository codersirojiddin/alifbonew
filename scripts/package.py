from pathlib import Path
from zipfile import ZipFile, ZIP_DEFLATED

root = Path(__file__).resolve().parents[1]
output = root / 'dist/downloads/Harf-Windows-source.zip'
output.parent.mkdir(parents=True, exist_ok=True)
with ZipFile(output, 'w', ZIP_DEFLATED) as archive:
    for path in sorted((root / 'desktop').rglob('*')):
        if path.is_file() and 'out' not in path.relative_to(root / 'desktop').parts:
            archive.write(path, path.relative_to(root / 'desktop'))
print(output)
