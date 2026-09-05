$projectRoot = Split-Path -Parent $PSScriptRoot
Set-Location (Join-Path $projectRoot 'web')

& npm.cmd run dev
exit $LASTEXITCODE
