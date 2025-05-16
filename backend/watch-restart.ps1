# watch-restart.ps1

Write-Host "🚀 Monitor de alterações para reinício automático do container Go"
Write-Host "📂 Monitorando arquivos .go em: $PWD"
Write-Host "🔄 Pressione Ctrl+C para encerrar o monitor"
Write-Host "⏱️ Iniciando monitoramento..."

$containerName = "appmercado-backend-dev"
$lastChange = Get-Date

while ($true) {
    # Procura por arquivos Go modificados após o último reinício
    $files = Get-ChildItem -Path . -Recurse -Include "*.go" | Where-Object { $_.LastWriteTime -gt $lastChange }
    
    if ($files.Count -gt 0) {
        $changedFiles = $files | ForEach-Object { $_.FullName.Replace("$PWD\", "") }
        
        Write-Host ""
        Write-Host "🔔 Alterações detectadas em $($files.Count) arquivo(s):" -ForegroundColor Yellow
        foreach ($file in $changedFiles) {
            Write-Host "   • $file" -ForegroundColor Cyan
        }
        
        Write-Host "🔄 Reiniciando container $containerName..." -ForegroundColor Yellow
        docker restart $containerName
        
        $lastChange = Get-Date
        
        # Espera 5 segundos para o container inicializar
        Write-Host "⏳ Aguardando inicialização do servidor..."
        Start-Sleep -Seconds 5
        
        Write-Host "✅ Container reiniciado com sucesso!" -ForegroundColor Green
        Write-Host "🕒 $(Get-Date -Format 'HH:mm:ss'): Monitorando novamente..." -ForegroundColor Gray
    }
    
    # Verifica a cada 2 segundos
    Start-Sleep -Seconds 2
}