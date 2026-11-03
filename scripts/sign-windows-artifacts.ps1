param(
    [Parameter(Mandatory=$true)][string[]]$Files,
    [Parameter(Mandatory=$true)][string]$CertificateThumbprint,
    [Parameter(Mandatory=$true)][string]$TimestampUrl,
    [Parameter(Mandatory=$true)][string]$SignTool
)
$ErrorActionPreference = 'Stop'
if ($CertificateThumbprint -notmatch '^[A-Fa-f0-9]{40}$') { throw 'Use an existing code-signing certificate thumbprint' }
$TaskTimestamp = [uri]$TimestampUrl
if (-not $TaskTimestamp.IsAbsoluteUri -or $TaskTimestamp.Scheme -ne 'https' -or $TaskTimestamp.UserInfo) { throw 'A trusted HTTPS RFC3161 timestamp URL is required' }
foreach ($TaskFile in $Files) {
    & $SignTool sign /sha1 $CertificateThumbprint /fd SHA256 /tr $TimestampUrl /td SHA256 $TaskFile
    if ($LASTEXITCODE) { throw 'Code signing failed' }
    if ((Get-AuthenticodeSignature -LiteralPath $TaskFile).Status -ne 'Valid') { throw 'The signed artifact does not have a valid trusted signature' }
}
