# Afritech Online Database Setup
# Run this script in PowerShell to create the database

Write-Host "=== Creating MariaDB Database ===" -ForegroundColor Cyan

# Database credentials
$dbName = "afritechonline"
$dbUser = "root"
$dbPassword = "root"
$mariadbPath = "C:\Program Files\MariaDB 12.3\bin\mariadb.exe"

# Create database
Write-Host "Creating database '$dbName'..."
$sql = "CREATE DATABASE IF NOT EXISTS $dbName CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
& $mariadbPath -u $dbUser -p$dbPassword -e $sql

Write-Host ""
Write-Host "=== Database Setup Complete ===" -ForegroundColor Green
Write-Host "Database: $dbName"
Write-Host "User: $dbUser"
Write-Host "Host: localhost"
Write-Host "Port: 3306"
Write-Host ""
Write-Host "Next steps:"
Write-Host "1. Run ./start.sh to start the application"
