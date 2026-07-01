!include "MUI2.nsh"

Name "GoSteamRestarter"
OutFile "${OUTFILE}"
InstallDir ""
RequestExecutionLevel highest

Var InstallMode

!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_LANGUAGE "SimpChinese"

Section /o "为所有用户安装（需要管理员权限）" SecAllUsers
  StrCpy $InstallMode "allusers"
SectionEnd

Section "仅为当前用户安装" SecCurrentUser
  StrCpy $InstallMode "currentuser"
SectionEnd

Function .onInit
  ; Default to current user
  StrCpy $InstallMode "currentuser"
  StrCpy $INSTDIR "$LOCALAPPDATA\GoSteamRestarter"
FunctionEnd

Function .onSelChange
  ; Make sections mutually exclusive
  SectionGetFlags ${SecAllUsers} $0
  SectionGetFlags ${SecCurrentUser} $1

  IntOp $0 $0 & 1
  IntCmp $0 1 allUsers
  ; Current user selected
  StrCpy $INSTDIR "$LOCALAPPDATA\GoSteamRestarter"
  Goto done
  allUsers:
  ; All users selected - uncheck current user
  SectionGetFlags ${SecCurrentUser} $1
  IntOp $1 $1 & 0xFFFFFFFE
  SectionSetFlags ${SecCurrentUser} $1
  StrCpy $INSTDIR "$PROGRAMFILES\GoSteamRestarter"
  done:
FunctionEnd

Section "安装文件" SecInstall
  SectionIn RO
  SetOutPath $INSTDIR
  File "${BINDIR}\gosteamrestarter-cli.exe"
  File "${BINDIR}\gosteamrestarter-desktop.exe"

  ; Add to PATH based on install mode
  StrCmp $InstallMode "allusers" 0 +3
    EnVar::SetHKLM
    Goto addpath
    EnVar::SetHKCU
  addpath:
  EnVar::AddValue "PATH" "$INSTDIR"

  ; Create uninstaller
  WriteUninstaller "$INSTDIR\uninstall.exe"

  ; Start menu shortcuts
  StrCmp $InstallMode "allusers" 0 +3
    SetShellVarContext all
    Goto shortcuts
    SetShellVarContext current
  shortcuts:
  CreateDirectory "$SMPROGRAMS\GoSteamRestarter"
  CreateShortcut "$SMPROGRAMS\GoSteamRestarter\GoSteamRestarter Desktop.lnk" "$INSTDIR\gosteamrestarter-desktop.exe"
  CreateShortcut "$SMPROGRAMS\GoSteamRestarter\Uninstall.lnk" "$INSTDIR\uninstall.exe"

  ; Save install mode for uninstaller
  FileOpen $0 "$INSTDIR\.installmode" w
  FileWrite $0 $InstallMode
  FileClose $0
SectionEnd

Section "Uninstall"
  ; Read install mode
  FileOpen $0 "$INSTDIR\.installmode" r
  FileRead $0 $1
  FileClose $0

  Delete "$INSTDIR\gosteamrestarter-cli.exe"
  Delete "$INSTDIR\gosteamrestarter-desktop.exe"
  Delete "$INSTDIR\.installmode"
  Delete "$INSTDIR\uninstall.exe"
  RMDir "$INSTDIR"

  ; Remove from PATH
  StrCmp $1 "allusers" 0 +3
    EnVar::SetHKLM
    Goto rmpath
    EnVar::SetHKCU
  rmpath:
  EnVar::DeleteValue "PATH" "$INSTDIR"

  ; Remove shortcuts
  StrCmp $1 "allusers" 0 +3
    SetShellVarContext all
    Goto rmshortcuts
    SetShellVarContext current
  rmshortcuts:
  Delete "$SMPROGRAMS\GoSteamRestarter\GoSteamRestarter Desktop.lnk"
  Delete "$SMPROGRAMS\GoSteamRestarter\Uninstall.lnk"
  RMDir "$SMPROGRAMS\GoSteamRestarter"
SectionEnd
