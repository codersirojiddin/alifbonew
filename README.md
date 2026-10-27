# Harf

O‘zbek alifbosi uchun onlayn konvertor va **Go’da yozilgan Windows yordamchisi**.

- `dist/`: HTML/CSS/JavaScript sayt va eski ↔ yangi lotin konvertori.
- `desktop/`: Go moduli, native Win32 interfeys, AltGr, tray menyusi,
  ixtiyoriy avtoalmashtirish va joriy foydalanuvchi uchun installer.
- `.github/workflows/windows-preview.yml`: Linux/Windows testlari, Go build,
  native smoke test, SHA-256 va GitHub prerelease.

## Yuklab olish

[Releases](https://github.com/codersirojiddin/alifbonew/releases)dan **Go** deb
belgilangan eng yangi talqinni tanlang: `Harf-Setup.exe` — o‘rnatish uchun,
`Harf.exe` — o‘rnatmasdan ishlatish uchun. Oldingi Harf ochiq bo‘lsa, avval tray
menyusidan yoping. Go talqini .NET yoki boshqa qo‘shimcha runtime talab qilmaydi.

## Sayt

`dist/`ni statik hostingga joylash mumkin. Matn serverga yuborilmaydi. Konvertor
UTF-8 .txt/.md/.csv/.srt importi, harf registri, apostrof variantlari, havola/emailni
saqlash, nusxalash va .txt eksportini qo‘llaydi. Google Fonts yuklanmasa, tizim
shriftlari ishlatiladi.

```sh
node tests/converter.test.cjs
node --check dist/app.js
cd desktop
go test ./...
```

## Windows

Windows 10/11 x64. Go + native Win32; tashqi Go kutubxonasi, CGO va .NET yo‘q.
Batafsil: [desktop/README.md](desktop/README.md).

AltGr + O/G/S/C → ö/ğ/ş/ç. Shift → katta harf. Ctrl+Alt+Y → avtomatik rejim.
Avtomatik rejim va Windows bilan ishga tushish avval o‘chiq turadi. Bu fon yordamchisi,
Windows til menyusiga alohida klaviatura tartibini qo‘shmaydi. Hozircha preview;
kundalik yozish va o‘rnatish qo‘lda tekshirilishi kerak.

`python3 scripts/package.py` Go manba kodini yuklab olinadigan ZIPga yig‘adi.
Ilhom manbai: UzKlaviatura. Uning kodi yoki dizayn aktivlari ko‘chirilmagan.
