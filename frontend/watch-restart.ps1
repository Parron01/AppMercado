# watch-restart.ps1 para frontend Angular

Write-Host "🚀 Monitor de alterações para reinício automático do container Angular"
Write-Host "📂 Monitorando arquivos Angular em: $PWD"
Write-Host "🔄 Pressione Ctrl+C para encerrar o monitor"
Write-Host "⏱️ Iniciando monitoramento..."

$containerName = "appmercado-frontend-dev"
$lastChange = Get-Date

while ($true) {
    # Procura por arquivos Angular modificados após o último reinício
    $files = Get-ChildItem -Path . -Recurse -Include "*.ts", "*.html", "*.scss", "*.css", "*.json" | 
             Where-Object { $_.LastWriteTime -gt $lastChange -and -not $_.FullName.Contains("node_modules") -and -not $_.FullName.Contains("dist") }
    
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
        Write-Host "⏳ Aguardando inicialização do servidor Angular..."
        Start-Sleep -Seconds 8
        
        Write-Host "✅ Container reiniciado com sucesso!" -ForegroundColor Green
        Write-Host "🕒 $(Get-Date -Format 'HH:mm:ss'): Monitorando novamente..." -ForegroundColor Gray
        Write-Host "🌐 Acesse a aplicação em: http://localhost:4200" -ForegroundColor Cyan
    }
    
    # Verifica a cada 2 segundos
    Start-Sleep -Seconds 2
}