!include "MUI2.nsh"
!include "EnvVarUpdate.nsh"

Name "GoSteamRestarter"
OutFile "${OUTFILE}"
InstallDir "$PROGRAMFILES\GoSteamRestarter"
RequestExecutionLevel admin

!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_LANGUAGE "SimpChinese"

Section "Install"
  SetOutPath $INSTDIR
  File "gosteamrestarter-cli.exe"
  File "gosteamrestarter-desktop.exe"

  ; Add to PATH
  ${EnvVarUpdate} $0 "PATH" "A" "HKLM" "$INSTDIR"

  ; Create uninstaller
  WriteUninstaller "$INSTDIR\uninstall.exe"

  ; Start menu shortcuts
  CreateDirectory "$SMPROGRAMS\GoSteamRestarter"
  CreateShortcut "$SMPROGRAMS\GoSteamRestarter\GoSteamRestarter Desktop.lnk" "$INSTDIR\gosteamrestarter-desktop.exe"
  CreateShortcut "$SMPROGRAMS\GoSteamRestarter\Uninstall.lnk" "$INSTDIR\uninstall.exe"
SectionEnd

Section "Uninstall"
  Delete "$INSTDIR\gosteamrestarter-cli.exe"
  Delete "$INSTDIR\gosteamrestarter-desktop.exe"
  Delete "$INSTDIR\uninstall.exe"
  RMDir "$INSTDIR"

  ; Remove from PATH
  ${un.EnvVarUpdate} $0 "PATH" "R" "HKLM" "$INSTDIR"

  ; Remove shortcuts
  Delete "$SMPROGRAMS\GoSteamRestarter\GoSteamRestarter Desktop.lnk"
  Delete "$SMPROGRAMS\GoSteamRestarter\Uninstall.lnk"
  RMDir "$SMPROGRAMS\GoSteamRestarter"
SectionEnd
