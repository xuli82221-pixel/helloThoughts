param(
    [Parameter(Mandatory=$true)][ValidateSet('codex','claude-code')][string]$Client,
    [string]$Executable = (Join-Path $PSScriptRoot 'memo.exe'),
    [Parameter(Mandatory=$true)][ValidatePattern('^[a-fA-F0-9]{64}$')][string]$ExpectedSha256,
    [string]$ThoughtsHome = 'auto',
    [switch]$Apply
)
$ErrorActionPreference = 'Stop'
if ([Environment]::OSVersion.Platform -ne [PlatformID]::Win32NT -or -not [Environment]::Is64BitOperatingSystem) { throw 'Windows x64 is required' }
$taskExe = (Resolve-Path -LiteralPath $Executable).Path
if (-not $ThoughtsHome -or $ThoughtsHome -eq 'auto') {
    $taskProfile = [Environment]::GetFolderPath([Environment+SpecialFolder]::UserProfile)
    if (-not $taskProfile) { throw 'Cannot locate the Windows account profile; no registration was changed' }
    $ThoughtsHome = Join-Path $taskProfile 'Thoughts'
}
if ($ThoughtsHome -notmatch '^[a-zA-Z]:[\\/]' -or $ThoughtsHome.StartsWith('\\') -or $ThoughtsHome.StartsWith('//')) { throw 'Select an absolute local Thoughts home directory' }
$taskHash = (Get-FileHash -LiteralPath $taskExe -Algorithm SHA256).Hash.ToLowerInvariant()
if ($taskHash -ne $ExpectedSha256.ToLowerInvariant()) { throw 'SHA256 mismatch. Nothing was installed.' }
$taskHome = [IO.Path]::GetFullPath($ThoughtsHome)
$taskCli = if ($Client -eq 'codex') { 'codex' } else { 'claude' }
Write-Output ('Verified executable: ' + $taskExe)
Write-Output ('Client: ' + $Client + '; shared Thoughts home: ' + $taskHome)
Write-Output 'This registers one local stdio server named thoughts using the client CLI. It does not initialize, migrate or overwrite a library, download or run the server, or change system security settings.'
if (-not $Apply) { Write-Output 'Inspection only. Add -Apply after reviewing the executable, checksum and proposed registration.'; return }
$taskCommand = Get-Command $taskCli -ErrorAction Stop
if ($Client -eq 'codex') {
    $taskSavedPreference = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    $taskExistingRaw = & $taskCommand.Source mcp get thoughts --json 2>$null
    $taskGetExit = $LASTEXITCODE
    $ErrorActionPreference = $taskSavedPreference
    if ($taskGetExit -eq 0) {
        $taskExisting = ($taskExistingRaw | Out-String) | ConvertFrom-Json
        $taskTransport = $taskExisting.transport
        if ($taskTransport.command -eq $taskExe -and @($taskTransport.args).Count -eq 1 -and $taskTransport.args[0] -eq 'mcp' -and $taskTransport.env.THOUGHTS_HOME -eq $taskHome) { Write-Output 'Already registered with the same executable and home.'; return }
        throw 'A thoughts MCP entry already exists. Inspect it before choosing an explicit update; no entry was replaced.'
    }
    $taskArguments = @('mcp','add','thoughts','--env',('THOUGHTS_HOME='+$taskHome),'--',$taskExe,'mcp')
} else {
    $taskClaudeDir = if ($env:CLAUDE_CONFIG_DIR) { $env:CLAUDE_CONFIG_DIR } else { $env:USERPROFILE }
    $taskClaudeConfig = if ($env:CLAUDE_CONFIG_DIR) { Join-Path $taskClaudeDir '.claude.json' } else { Join-Path $taskClaudeDir '.claude.json' }
    if (Test-Path -LiteralPath $taskClaudeConfig) {
        $taskExisting = (Get-Content -LiteralPath $taskClaudeConfig -Raw | ConvertFrom-Json).mcpServers.thoughts
        if ($taskExisting) {
            if ($taskExisting.command -eq $taskExe -and @($taskExisting.args).Count -eq 1 -and $taskExisting.args[0] -eq 'mcp' -and $taskExisting.env.THOUGHTS_HOME -eq $taskHome) { Write-Output 'Already registered with the same executable and home.'; return }
            throw 'A user-level thoughts entry already exists. No entry was replaced.'
        }
    }
    $taskArguments = @('mcp','add','--transport','stdio','--scope','user','thoughts','--env',('THOUGHTS_HOME='+$taskHome),'--',$taskExe,'mcp')
}
& $taskCommand.Source @taskArguments
if ($LASTEXITCODE) { throw 'Client CLI did not confirm registration. Inspect its status before retrying.' }
Write-Output 'Registration command completed. Verify /mcp and thoughts_status in a new client session; this alone is not a connection test.'
