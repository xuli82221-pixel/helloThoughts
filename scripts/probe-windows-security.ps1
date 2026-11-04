param(
    [Parameter(Mandatory=$true)][string]$Candidate,
    [Parameter(Mandatory=$true)][string]$Output,
    [ValidateSet('Any','On','Off')][string]$ExpectedSAC='Any'
)
# Read-only security probe. Never changes policy, certificates, exclusions or services.
# Candidate must be an unpacked Thoughts distribution produced by the build script.
$ErrorActionPreference='Stop'
$taskRoot=(Resolve-Path -LiteralPath $Candidate).Path
$taskIndex=Get-Content -LiteralPath (Join-Path $taskRoot 'index.json') -Raw | ConvertFrom-Json
$taskRecords=@()
$taskHashes=@()
function Get-BytesHash([byte[]]$Bytes) {
    $taskHasher=[Security.Cryptography.SHA256]::Create()
    try { return ([BitConverter]::ToString($taskHasher.ComputeHash($Bytes))).Replace('-','').ToLowerInvariant() }
    finally { $taskHasher.Dispose() }
}
foreach ($taskPackage in $taskIndex.packages) {
    $taskArchive=Join-Path $taskRoot $taskPackage.file
    $taskArchiveHash=(Get-FileHash -LiteralPath $taskArchive -Algorithm SHA256).Hash.ToLowerInvariant()
    Add-Type -AssemblyName System.IO.Compression.FileSystem
    $taskZip=[IO.Compression.ZipFile]::OpenRead($taskArchive)
    try {
        $taskInner=if ($taskPackage.file.EndsWith('.mcpb')) {'server/thoughts.exe'} else {'thoughts.exe'}
        $taskStream=$taskZip.GetEntry($taskInner).Open()
        $taskMemory=[IO.MemoryStream]::new()
        try { $taskStream.CopyTo($taskMemory); $taskBytes=$taskMemory.ToArray() }
        finally { $taskStream.Dispose(); $taskMemory.Dispose() }
        $taskHash=Get-BytesHash $taskBytes
        $taskHashes+=$taskHash
        $taskPE=[BitConverter]::ToInt32($taskBytes,0x3c)
        if ([BitConverter]::ToUInt16($taskBytes,$taskPE+24) -ne 0x20b) { throw 'Expected PE32+ executable' }
        $taskSecurityDirectory=$taskPE+24+112+4*8
        $taskSource=if ($taskInner.StartsWith('server/')) {'sources/claude-desktop-windows-amd64/server/thoughts.exe'} else {'sources/windows-amd64/thoughts.exe'}
        $taskExe=Join-Path $taskRoot $taskSource
        if ((Get-FileHash -LiteralPath $taskExe).Hash.ToLowerInvariant() -ne $taskHash) { throw 'Archive and unpacked executable differ' }
        $taskSignature=Get-AuthenticodeSignature -LiteralPath $taskExe
        $taskRecords+=[ordered]@{
            package=$taskPackage.file; packageSha256=$taskArchiveHash
            packageHashMatchesIndex=($taskArchiveHash -eq $taskPackage.sha256)
            innerExecutable=$taskInner; executableSha256=$taskHash
            executableHashMatchesIndex=($taskHash -eq $taskPackage.executableSha256)
            authenticode=[string]$taskSignature.Status
            signerPresent=($null -ne $taskSignature.SignerCertificate)
            peCertificateOffset=[BitConverter]::ToUInt32($taskBytes,$taskSecurityDirectory)
            peCertificateSize=[BitConverter]::ToUInt32($taskBytes,$taskSecurityDirectory+4)
            archiveMarkOfTheWeb=($null -ne (Get-Item -LiteralPath $taskArchive -Stream Zone.Identifier -ErrorAction SilentlyContinue))
            executableMarkOfTheWeb=($null -ne (Get-Item -LiteralPath $taskExe -Stream Zone.Identifier -ErrorAction SilentlyContinue))
        }
    } finally { $taskZip.Dispose() }
}
$taskMP=Get-MpComputerStatus
$taskOS=Get-CimInstance Win32_OperatingSystem
$taskProcesses=@(Get-CimInstance Win32_Process -Filter "name = 'thoughts.exe'" | ForEach-Object {
    if ($_.ExecutablePath) {
        $taskRunningHash=(Get-FileHash -LiteralPath $_.ExecutablePath).Hash.ToLowerInvariant()
        [ordered]@{pid=$_.ProcessId;sha256=$taskRunningHash;matchesCandidate=($taskHashes -contains $taskRunningHash);authenticode=[string](Get-AuthenticodeSignature -LiteralPath $_.ExecutablePath).Status}
    }
})
$taskEvents=@()
$taskLogStatus=@()
foreach ($taskLog in @('Microsoft-Windows-CodeIntegrity/Operational','Microsoft-Windows-Windows Defender/Operational')) {
    try {
        $taskRead=@(Get-WinEvent -FilterHashtable @{LogName=$taskLog;StartTime=(Get-Date).AddDays(-7)} -MaxEvents 2000 -ErrorAction Stop)
        $taskMatching=@($taskRead | Where-Object {$_.Message -match '(?i)memo\.exe|thoughts[^a-z]'})
        foreach ($taskEvent in $taskMatching) {
            # Do not export raw messages: they may contain workstation paths.
            $taskEvents+=[ordered]@{log=$taskLog;id=$taskEvent.Id;utc=$taskEvent.TimeCreated.ToUniversalTime().ToString('o');recordId=$taskEvent.RecordId}
        }
        $taskLogStatus+=[ordered]@{log=$taskLog;read='OK';scanned=$taskRead.Count;matching=$taskMatching.Count;bounded=$true}
    } catch {
        $taskLogStatus+=[ordered]@{log=$taskLog;read='UNAVAILABLE_OR_NO_EVENTS';errorType=$_.Exception.GetType().Name}
    }
}
$taskResult=[ordered]@{
    capturedUtc=[DateTime]::UtcNow.ToString('o');version=$taskIndex.version
    os=[ordered]@{caption=$taskOS.Caption;version=$taskOS.Version;architecture=$taskOS.OSArchitecture}
    security=($taskMP | Select-Object SmartAppControlState,AntivirusEnabled,RealTimeProtectionEnabled,IsTamperProtected,AntivirusSignatureVersion,AMProductVersion)
    expectedSAC=$ExpectedSAC;expectedSACMatches=($ExpectedSAC -eq 'Any' -or [string]$taskMP.SmartAppControlState -eq $ExpectedSAC)
    packages=$taskRecords;allPackagesSameCore=(@($taskHashes | Select-Object -Unique).Count -eq 1)
    runningExecutables=$taskProcesses;eventLogQueries=$taskLogStatus;matchingEventMetadata=$taskEvents
    vmTools=@(Get-Command Get-VM,WindowsSandbox.exe,vmrun.exe,VBoxManage.exe -ErrorAction SilentlyContinue | ForEach-Object {$_.Name})
    securitySettingsChanged=$false;cleanEnvironmentAcceptance=$false
}
$taskResult | ConvertTo-Json -Depth 8 | Set-Content -LiteralPath $Output -Encoding utf8
if (-not $taskResult.expectedSACMatches) { throw 'SAC does not match requested test profile. No policy was changed. Stop this test.' }
Write-Output ('Read-only evidence saved. Version: '+$taskIndex.version+'; SAC: '+$taskMP.SmartAppControlState)
