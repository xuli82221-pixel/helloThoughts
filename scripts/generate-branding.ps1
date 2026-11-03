$ErrorActionPreference='Stop'
Add-Type -AssemblyName System.Drawing
$taskRoot=Split-Path -Parent $PSScriptRoot
$taskBitmap=New-Object System.Drawing.Bitmap 512,512
$taskCanvas=[System.Drawing.Graphics]::FromImage($taskBitmap)
$taskCanvas.SmoothingMode=[System.Drawing.Drawing2D.SmoothingMode]::AntiAlias
$taskCanvas.Clear([System.Drawing.Color]::White)
$taskCanvas.ScaleTransform(2,2)
$taskShape=New-Object System.Drawing.Drawing2D.GraphicsPath
foreach ($taskArc in @(@(8,8,104,104,180,90),@(144,8,104,104,270,90),@(144,144,104,104,0,90),@(8,144,104,104,90,90))) { $taskCanvas.Flush(); $taskShape.AddArc($taskArc[0],$taskArc[1],$taskArc[2],$taskArc[3],$taskArc[4],$taskArc[5]) }
$taskShape.CloseFigure()
$taskBrush=New-Object System.Drawing.SolidBrush ([System.Drawing.Color]::White)
$taskCanvas.FillPath($taskBrush,$taskShape)
$taskPen=New-Object System.Drawing.Pen ([System.Drawing.Color]::FromArgb(24,24,27)),12
$taskPen.StartCap=[System.Drawing.Drawing2D.LineCap]::Round
$taskPen.EndCap=[System.Drawing.Drawing2D.LineCap]::Round
$taskPen.LineJoin=[System.Drawing.Drawing2D.LineJoin]::Round
$taskBook=New-Object System.Drawing.Drawing2D.GraphicsPath
$taskBook.AddLine(147,64,77,64); $taskBook.AddArc(64,64,26,26,270,-90); $taskBook.AddLine(64,77,64,190); $taskBook.AddArc(64,177,26,26,180,-90); $taskBook.AddLine(77,203,171,203); $taskBook.AddArc(158,177,26,26,90,-90); $taskBook.AddLine(184,190,184,148)
$taskCanvas.DrawPath($taskPen,$taskBook)
$taskCanvas.DrawLine($taskPen,91,64,91,203)
foreach ($taskY in @(91,124,157)) { $taskCanvas.DrawLine($taskPen,52,$taskY,76,$taskY) }
$taskPoints=[System.Drawing.Point[]]@((New-Object System.Drawing.Point 133,143),(New-Object System.Drawing.Point 142,111),(New-Object System.Drawing.Point 193,60),(New-Object System.Drawing.Point 216,83),(New-Object System.Drawing.Point 165,134),(New-Object System.Drawing.Point 133,143))
$taskCanvas.DrawLines($taskPen,$taskPoints); $taskCanvas.DrawLine($taskPen,181,72,204,95)
$taskBitmap.Save((Join-Path $taskRoot 'branding/icon.png'),[System.Drawing.Imaging.ImageFormat]::Png)
$taskBook.Dispose(); $taskPen.Dispose(); $taskBrush.Dispose(); $taskShape.Dispose(); $taskCanvas.Dispose(); $taskBitmap.Dispose()
foreach ($taskDir in @('plugins/memo/assets','adapters/claude-desktop/assets','web/public/assets','plugins/memo/skills/memo-setup/assets','plugins/memo/skills/open-thoughts/assets')) {
    New-Item -ItemType Directory -Force -Path (Join-Path $taskRoot $taskDir) | Out-Null
    Copy-Item -LiteralPath (Join-Path $taskRoot 'branding/icon.png') -Destination (Join-Path $taskRoot ($taskDir+'/icon.png'))
}
Copy-Item -LiteralPath (Join-Path $taskRoot 'branding/icon.svg') -Destination (Join-Path $taskRoot 'web/public/assets/icon.svg')
