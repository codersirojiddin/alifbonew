# Harf Windows helper — preview source

This is original source code, not a compiled or Windows-tested release.
No executable is included in this source ZIP. The website must not advertise
an available .exe until the Windows validation below is complete.

## Build on Windows 10 / 11

1. Extract this entire ZIP into one directory.
2. Run `Build.cmd`, or run `powershell -NoProfile -File .\build.ps1`.
3. The script compiles with the Windows .NET Framework compiler, runs converter
   self-tests and creates `out/Harf.exe` and `out/Harf-Setup.exe`.
4. Test `out/Harf.exe` first. Build success alone does not verify keyboard input.
5. `Harf-Setup.exe` copies the application into `%LOCALAPPDATA%\Harf`, adds a
   Start Menu shortcut and starts it. No administrator access is requested.
   Close an existing Harf instance before installing an update.

The .NET Framework 4.8 compiler is required; if absent, install its Developer Pack.
`Build.cmd` uses a process-scoped PowerShell execution policy and does not change
your machine's persistent execution-policy setting.

## Controls

- AltGr (or Ctrl+Alt) + O / G / S / C inserts ö / ğ / ş / ç.
- Add Shift for uppercase. Caps Lock is not used by these shortcuts.
- Ctrl+Alt+Y toggles automatic replacement, which starts OFF each session.
- Automatic mode recognizes ASCII `sh`, `ch`, `o'`, `g'`, and backtick variants
  on a US Latin physical layout. This is an input helper, NOT a registered
  Windows keyboard layout. OS layout DLL installation is not implemented.
- Tray menu: show window, toggle automatic replacement, exit.
- Closing the window hides it to the tray; use Exit to stop it entirely.
- Windows startup is opt-in through the checkbox.
- The window includes forward/reverse offline text conversion and copy.

The helper does not log, save or transmit keystrokes. The automatic mode keeps
only the preceding candidate character in memory. Common code editors and
terminals are excluded. Native password controls are excluded where detectable;
browser password fields cannot be reliably identified, so keep automatic mode
OFF while entering passwords. AltGr shortcuts remain active while Harf runs.
Keyboard hooks and shortcut collisions can affect other applications; exit Harf
to release all shortcuts. Input into elevated apps may be blocked by Windows.

Automatic replacement is intentionally a basic typing aid: unlike the full-text
converter, it cannot identify URLs/emails/foreign-language words before you finish
typing. Use the website or desktop text box when contextual preservation matters.

## Uninstall

Exit Harf from the tray menu. Run `%LOCALAPPDATA%\Harf\Uninstall.cmd`.
It removes the startup entry, Start Menu shortcut and application directory.
The helper does not modify system keyboard layouts.

## Required Windows acceptance checks before a release

- Run conversion self-tests through `build.ps1`.
- Notepad: AltGr letter insertion, Shift uppercase, automatic digraphs,
  mixed case, rapid typing, Backspace, mouse clicks, navigation and focus changes.
- Confirm Ctrl+Alt+Y, Caps Lock, held keys, tray exit and duplicate-instance handling.
- Confirm shortcuts do not type into a newly focused window after activation.
- Verify automatic mode is off at startup and disabled in excluded applications.
- Check password-control exclusions and document the browser limitation.
- Install on a clean standard-user account; test Start Menu launch, opt-in startup,
  update after exit, uninstall, 32-bit/64-bit behavior and input restrictions.
- Test with a screen reader and a non-US current keyboard layout.
- Review hook latency and antivirus results; sign release binaries when possible.

## Website download integration

After Windows validation, publish `Harf-Setup.exe` and its SHA-256 hash in your own
GitHub Release. Replace the source-package link in the website with the actual
release download URL and update the Windows section/FAQ to match tested behavior.
Never link the sample project's installer as though it were a Harf release.
