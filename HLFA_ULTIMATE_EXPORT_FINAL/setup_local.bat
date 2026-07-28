@echo off
echo [*] Menyiapkan Lingkungan HLFA Lokal...
python -m venv venv
call venv\Scripts\activate
echo [*] Menginstall Dependensi...
pip install -e .
echo [*] Verifikasi Instalasi...
python -m hlfa --help
echo.
echo [OK] HLFA Siap! Gunakan command: python -m hlfa scan [URL]
pause
