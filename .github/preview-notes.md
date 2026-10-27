Harf Windows input helper — preview build.

- **Harf-Setup.exe** installs for your Windows user. No administrator rights required.
- **Harf.exe** runs without installation.
- Windows 10/11 with .NET Framework 4.8 required.
- AltGr + O/G/S/C inserts ö/ğ/ş/ç. Add Shift for uppercase.
- Automatic replacement is opt-in. Toggle it in the app or with Ctrl+Alt+Y.
- Website-style forward/reverse conversion is available in the app window.
- Startup is opt-in; closing the window hides it to the tray. Use Exit to stop it.
- To uninstall, exit Harf and run `%LOCALAPPDATA%\Harf\Uninstall.cmd`.

This build is unsigned. CI compiles the application and checks its converter;
that does not replace manual testing of typing, installation and uninstallation
on a Windows PC. This helper does not install an OS keyboard-layout DLL.

Automatic replacement assumes a US Latin layout. It cannot reliably identify
browser password fields; keep automatic mode off while entering sensitive data.
Common editors and terminals are excluded from automatic replacement.

Compare the downloaded files against SHA256SUMS.txt. Report problems in this
repository with the app version, Windows version and steps to reproduce.
