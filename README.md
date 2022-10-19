# Enable colors/escape sequences in cmd.exe
```
REG ADD HKCU\CONSOLE /f /v VirtualTerminalLevel /t REG_DWORD /d 1
```
