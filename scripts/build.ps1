param(
    [Parameter(Mandatory=$true)][string]$Python,
    [string]$Go='',
    [string]$Output='',
    [switch]$PublicRelease,
    [switch]$UnsignedPreview,
    [string]$SourceCommit='',
    [string]$BuildUrl='',
    [string]$CertificateThumbprint='',
    [string]$TimestampUrl='',
    [string]$SignTool='',
    [Parameter(Mandatory=$true)][string]$InnoCompiler
)
# Keep a single Windows package path; historical M0 output is never overwritten here.
$ErrorActionPreference='Stop'
& (Join-Path $PSScriptRoot 'build-desktop-adapters.ps1') @PSBoundParameters
