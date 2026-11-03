param(
    [Parameter(Mandatory=$true)][string]$BinDir,
    [Parameter(Mandatory=$true)][string]$Output,
    [Parameter(Mandatory=$true)][string]$Compiler,
    [string]$InstallerId = 'Thoughts.Desktop.Windows',
    [string]$CertificateThumbprint = '',
    [string]$TimestampUrl = '',
    [string]$SignTool = ''
)
$ErrorActionPreference = 'Stop'
$TaskRoot = Split-Path -Parent $PSScriptRoot
$TaskBinary = Join-Path $BinDir 'windows-amd64/memo.exe'
$TaskVersion = (Get-Content -LiteralPath (Join-Path $TaskRoot 'adapters/claude-desktop/manifest.json') -Raw | ConvertFrom-Json).version
if ($TaskVersion -notmatch '^[0-9]+\.[0-9]+\.[0-9]+(?:-[A-Za-z0-9.-]+)?$') { throw 'Invalid version' }
if ($InstallerId -notmatch '^Thoughts\.[A-Za-z0-9.-]+$') { throw 'Invalid installer ID' }
if (Test-Path -LiteralPath $Output) { throw 'Use a new installer output directory; no overwrite' }
$TaskCompilerSignature = Get-AuthenticodeSignature -LiteralPath $Compiler
if ($TaskCompilerSignature.Status -ne 'Valid' -or $TaskCompilerSignature.SignerCertificate.Subject -notmatch 'Pyrsys B.V.') { throw 'Use the signed official Inno Setup 7 compiler' }
$TaskHash = (Get-FileHash -LiteralPath $TaskBinary -Algorithm SHA256).Hash.ToLowerInvariant()
$TaskCompilerArgs = @(
    '/Qp', '/NI', ('/DCoreDir=' + (Split-Path -Parent (Resolve-Path -LiteralPath $TaskBinary).Path)),
    ('/DProjectRoot=' + $TaskRoot), ('/DCoreSHA256=' + $TaskHash),
    ('/DAppVersion=' + $TaskVersion), ('/DOutputFolder=' + [IO.Path]::GetFullPath($Output)),
    ('/DInstallerId=' + $InstallerId), (Join-Path $TaskRoot 'release/windows/Thoughts.iss')
)
if ($CertificateThumbprint) {
    if ($CertificateThumbprint -notmatch '^[A-Fa-f0-9]{40}$' -or -not $SignTool -or $SignTool.Contains('"')) { throw 'A valid certificate thumbprint and signing tool path are required' }
    $TaskTimestamp = [uri]$TimestampUrl
    if (-not $TaskTimestamp.IsAbsoluteUri -or $TaskTimestamp.Scheme -ne 'https' -or $TaskTimestamp.UserInfo -or $TimestampUrl.Contains('"')) { throw 'A trusted HTTPS timestamp URL is required' }
    # Inno signs both its generated uninstaller and final installer.
    $TaskSigningCommand = '$q' + $SignTool + '$q sign /sha1 ' + $CertificateThumbprint + ' /fd SHA256 /tr $q' + $TimestampUrl + '$q /td SHA256 $f'
    $TaskCompilerArgs = @('/DEnableSigning=1',('/Sthoughts=' + $TaskSigningCommand)) + $TaskCompilerArgs
}
& $Compiler @TaskCompilerArgs
if ($LASTEXITCODE) { throw 'Inno Setup compilation failed' }
$TaskSetup = Join-Path $Output ('thoughts-' + $TaskVersion + '-windows-amd64-setup.exe')
if (-not (Test-Path -LiteralPath $TaskSetup -PathType Leaf)) { throw 'Installer output is missing' }
if ($CertificateThumbprint -and (Get-AuthenticodeSignature -LiteralPath $TaskSetup).Status -ne 'Valid') { throw 'Installer signature verification failed' }
[ordered]@{
    version=$TaskVersion; file=[IO.Path]::GetFileName($TaskSetup)
    sha256=(Get-FileHash -LiteralPath $TaskSetup -Algorithm SHA256).Hash.ToLowerInvariant()
    executableSha256=$TaskHash; installerId=$InstallerId
    compilerSha256=(Get-FileHash -LiteralPath $Compiler -Algorithm SHA256).Hash.ToLowerInvariant()
    publicDistributionReady=$false
} | ConvertTo-Json | Set-Content -LiteralPath (Join-Path $Output 'installer-build.json') -Encoding UTF8
