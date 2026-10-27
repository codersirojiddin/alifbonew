# Harf — Go Windows helper

The Windows application and self-installer are written in Go. Native Win32 controls,
Go standard library only: no .NET, Electron, WebView, CGO or third-party libraries.
Target: Windows 10/11 x64. Running the built app does not require Go to be installed.

## Build

Install Go and run `Build.cmd`, or `powershell -File build.ps1` in this directory.
The script runs Go tests and vet, then creates `out/Harf.exe` and `out/Harf-Setup.exe`.
They contain the same Go binary; the filename selects installation mode.

Cross-compile from Linux/macOS:

```sh
cd desktop
go test ./...
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags='-H=windowsgui -s -w' -o Harf.exe ./cmd/harf
```

CI checks the portable converter on Linux, conversion and Win32 memory layouts on
Windows, and native window/hook registration with a bounded smoke test. Preview
installers and SHA-256 hashes are published under GitHub Releases. Manual typing,
installation and uninstallation tests are still needed before a stable release.

## Use

- Harf-Setup.exe installs into `%LOCALAPPDATA%\Harf` for the current user.
  Exit an older Harf instance before upgrading. No administrator access required.
- Harf.exe runs directly without installation.
- AltGr (Ctrl+Alt) + O/G/S/C inserts ö/ğ/ş/ç. Add Shift for uppercase.
- Automatic replacement starts OFF. Toggle it in the window/tray or Ctrl+Alt+Y.
- Automatic mode recognizes sh, ch, o', g' and backtick variants, and only runs
  with the English US keyboard layout to avoid incorrect physical-key mappings.
- Closing the window hides it to the tray. Use Exit to stop it.
- Windows startup is opt-in; only the current user's Harf Run registry entry is used.
- The window includes an offline forward/reverse converter and clipboard copy.

This background helper does not install an OS keyboard-layout DLL. It keeps only
the preceding candidate rune in memory and does not save or transmit keystrokes.
Common editors/terminals and detectable native password controls are excluded.
Browser password fields cannot be identified reliably: keep auto mode OFF while
entering passwords. Elevated apps can refuse input injection. Exit Harf to release
all registered shortcuts. Automatic replacement cannot protect unfinished URLs,
emails or foreign words in advance; use the full-text converter for those.

## Uninstall

Exit Harf and run `%LOCALAPPDATA%\Harf\Uninstall.cmd`. It removes the application,
Start Menu shortcut and current user's Harf startup entry.

## Manual Windows acceptance checks

Test installation, Start Menu launch, startup, upgrade and uninstall on a clean
standard-user account. Test all AltGr/Shift shortcuts, rapid typing, focus changes,
mouse clicks, Caps Lock, Backspace, auto-mode toggling, clipboard, excluded/elevated
apps, tray exit and duplicate-instance handling. Check the documented browser
password-field limitation. CI smoke testing does not cover every interaction.
