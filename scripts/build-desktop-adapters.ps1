param([Parameter(Mandatory=$true)][string]$Python, [string]$Go='', [string]$Output='', [switch]$PublicRelease, [string]$CertificateThumbprint='', [string]$TimestampUrl='', [string]$SignTool='', [Parameter(Mandatory=$true)][string]$InnoCompiler)
$ErrorActionPreference='Stop'
if ([Environment]::OSVersion.Platform -ne [PlatformID]::Win32NT) { throw 'This release entry runs on Windows only; macOS work is deferred' }
$taskRoot=Split-Path -Parent $PSScriptRoot
if (-not $Go) { $Go=Join-Path $taskRoot '.tools/go/bin/go.exe' }
$taskVersion=(Get-Content -LiteralPath (Join-Path $taskRoot 'adapters/claude-desktop/manifest.json') -Raw | ConvertFrom-Json).version
if (-not $Output) { $Output=Join-Path $taskRoot ('dist/windows/'+$taskVersion) }
if (Test-Path -LiteralPath $Output) { throw 'Output exists; use a new directory. No automatic deletion.' }
$taskBuildId=[guid]::NewGuid().ToString('N')
$taskBins=Join-Path $taskRoot ('.artifacts/windows-adapter-build-'+$taskBuildId)
New-Item -ItemType Directory -Force -Path (Join-Path $taskBins 'common') | Out-Null
& $Python -X utf8 (Join-Path $PSScriptRoot 'release_metadata.py')
if ($LASTEXITCODE) { throw 'Publisher metadata synchronization failed' }
& $Python -X utf8 (Join-Path $PSScriptRoot 'test_legal_notices.py')
if ($LASTEXITCODE) { throw 'Localized notice checks failed' }
if ($PublicRelease) {
    & $Python -X utf8 (Join-Path $PSScriptRoot 'release_metadata.py') --check-public
    if ($LASTEXITCODE) { throw 'Public publisher metadata is incomplete; no release package produced' }
}
Copy-Item -LiteralPath (Join-Path $taskRoot 'branding/icon.svg') -Destination (Join-Path $taskRoot 'web/public/assets/icon.svg')
Push-Location (Join-Path $taskRoot 'web')
try {
    & npm test; if ($LASTEXITCODE) { throw 'Frontend tests failed' }
    & npm run build; if ($LASTEXITCODE) { throw 'Frontend build failed' }
} finally { Pop-Location }
Push-Location $taskRoot
$taskOldOS=$env:GOOS; $taskOldArch=$env:GOARCH; $taskOldCGO=$env:CGO_ENABLED
try {
    $env:GOOS='windows'; $env:GOARCH='amd64'; $env:CGO_ENABLED='0'
    & $Go test ./... -count=1; if ($LASTEXITCODE) { throw 'Go tests failed or were blocked; no packages produced' }
    $env:CGO_ENABLED='0'
    foreach ($taskTarget in @('windows-amd64')) {
        $taskParts=$taskTarget.Split('-'); $env:GOOS=$taskParts[0]; $env:GOARCH=$taskParts[1]
        $taskFolder=Join-Path $taskBins $taskTarget
        New-Item -ItemType Directory -Force -Path $taskFolder | Out-Null
        $taskName=if ($env:GOOS -eq 'windows') { 'memo.exe' } else { 'memo' }
        & $Go build -trimpath -ldflags '-s -w' -o (Join-Path $taskFolder $taskName) ./cmd/memo
        if ($LASTEXITCODE) { throw ('Compile failed: '+$taskTarget) }
    }
    $env:GOOS=$taskOldOS; $env:GOARCH=$taskOldArch; $env:CGO_ENABLED=$taskOldCGO
    $taskBinary=Join-Path $taskBins 'windows-amd64/memo.exe'
    $taskInstaller=Join-Path $taskBins 'common/register-mcp.ps1'
    Copy-Item -LiteralPath (Join-Path $PSScriptRoot 'register-mcp.ps1') -Destination $taskInstaller
    if ($CertificateThumbprint) {
        & (Join-Path $PSScriptRoot 'sign-windows-artifacts.ps1') -Files @($taskBinary,$taskInstaller) -CertificateThumbprint $CertificateThumbprint -TimestampUrl $TimestampUrl -SignTool $SignTool
    }
    $taskSetupFolder=Join-Path $taskBins 'setup'
    & (Join-Path $PSScriptRoot 'build-windows-installer.ps1') -BinDir $taskBins -Output $taskSetupFolder -Compiler $InnoCompiler -CertificateThumbprint $CertificateThumbprint -TimestampUrl $TimestampUrl -SignTool $SignTool
    if ($LASTEXITCODE) { throw 'Offline installer build failed' }
    $taskSetup=Join-Path $taskSetupFolder ('thoughts-'+$taskVersion+'-windows-amd64-setup.exe')
    $taskSignatures=@(foreach ($taskFile in @($taskBinary,$taskInstaller,$taskSetup)) {
        $taskSignature=Get-AuthenticodeSignature -LiteralPath $taskFile
        [ordered]@{ file=[IO.Path]::GetFileName($taskFile);status=[string]$taskSignature.Status;sha256=(Get-FileHash -LiteralPath $taskFile -Algorithm SHA256).Hash.ToLowerInvariant();signer=if($taskSignature.SignerCertificate){$taskSignature.SignerCertificate.Subject}else{$null} }
    })
    $taskSignatures | ConvertTo-Json -Depth 4 | Set-Content -LiteralPath (Join-Path $taskBins 'common/signatures.json') -Encoding utf8
    if ($PublicRelease -and @($taskSignatures | Where-Object {$_.status -ne 'Valid'}).Count) { throw 'Valid core, registration helper and setup signatures are required for public distribution' }
    $taskGoRoot=(& $Go env GOROOT | Out-String).Trim()
    & node (Join-Path $PSScriptRoot 'notices.mjs') (Join-Path $taskBins 'common') $taskGoRoot
    if ($LASTEXITCODE) { throw 'Notices generation failed' }
    $taskPackArgs=@('-X','utf8',(Join-Path $PSScriptRoot 'pack-desktop-adapters.py'),'--bin-dir',$taskBins,'--output',$Output,'--setup-file',$taskSetup)
    if ($PublicRelease) { $taskPackArgs+='--public-release' }
    & $Python @taskPackArgs
    if ($LASTEXITCODE) { throw 'Packaging failed' }
    & $Python -X utf8 (Join-Path $PSScriptRoot 'validate-desktop-adapters.py') $Output
    if ($LASTEXITCODE) { throw 'Package validation failed' }
} finally {
    $env:GOOS=$taskOldOS; $env:GOARCH=$taskOldArch; $env:CGO_ENABLED=$taskOldCGO
    Pop-Location
}
