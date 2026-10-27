Harf 0.2 — rewritten in **Go**, using native Windows controls.
No .NET, Electron, WebView, CGO or third-party runtime required.

- Harf-Setup.exe installs for the current Windows user; Harf.exe runs directly.
- Windows 10/11 x64. Exit the older Harf before upgrading.
- AltGr + O/G/S/C → ö/ğ/ş/ç; add Shift for uppercase.
- Ctrl+Alt+Y toggles automatic replacement (OFF at startup).
- Startup is opt-in. Closing hides to tray; Exit stops the app.
- Includes offline text conversion and clipboard copy.
- Uninstall: exit Harf and run `%LOCALAPPDATA%\Harf\Uninstall.cmd`.

CI checks Go conversion on Linux/Windows, Win32 memory layouts and native
window/input-hook registration. Manual typing/install/uninstall checks remain
necessary before a stable release. Preview binaries are unsigned; compare SHA-256.

Auto mode uses English US only and cannot identify all browser password fields
or unfinished URLs. Keep it OFF for sensitive input. The helper does not install
an OS keyboard-layout DLL.
