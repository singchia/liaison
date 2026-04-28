; NSIS install/uninstall hooks for Liaison Desktop.
;
; Tauri's NSIS template invokes the macros below at well-defined
; points in the install / uninstall flow. We use NSIS_HOOK_PREUNINSTALL
; to call the bundled binary with --cleanup-credentials, which removes
; the per-user PAT entry from Windows Credential Manager before the
; uninstaller deletes program files. This way a clean uninstall
; doesn't leave a stale credential entry behind.

!macro NSIS_HOOK_PREUNINSTALL
  ; Best-effort. /TIMEOUT lets us not wait forever if the binary
  ; hangs. The cleanup itself is fast (one keyring delete call).
  ExecWait '"$INSTDIR\Liaison.exe" --cleanup-credentials'
!macroend
