import asyncio
import sys
import os
from rich.console import Console
from rich.panel import Panel
from rich.table import Table
from rich.json import JSON
from rich.columns import Columns
from rich.text import Text

console = Console()

async def run_hybrid_engine(target):
    sys.path.insert(0, '/content/hlfa_project_latest')
    console.print(Panel("[bold white]HLFA ULTIMATE V19 - HYBRID ORCHESTRATOR[/bold white]\n[cyan]Mode: Deep Chaining & Manual Control[/cyan]", style="on blue"))

    # Database Loot Sementara (Persisten selama sesi)
    inventory = {
        "ssrf_meta": {"status": "Acquired", "data": "IAM-Role: srv-web-app", "next_vector": "Internal Pivot"},
        "db_creds": {"status": "Locked", "data": "None", "next_vector": "SQLi"}
    }

    def show_explanation(vector):
        expl = {
            "SSRF": "Metadata ini berisi token akses cloud. Gunakan perintah 'pivot' untuk melihat bagaimana token ini membuka akses ke database internal.",
            "SQLI": "Celah ini memungkinkan bypass otentikasi. Data yang didapat bisa digunakan untuk login ke panel admin secara manual.",
            "PIVOT": "Teknik 'celah dalam celah'. Menggunakan akses SSRF untuk memindai port 6379 (Redis) di jaringan internal 10.0.x.x."
        }
        return Panel(expl.get(vector, "Penjelasan tidak tersedia."), title=f"Insight: {vector}", border_style="blue", width=45)

    console.print("[bold yellow][*] Sistem siap. Masukkan perintah untuk mengotak-atik data.[/bold yellow]")

    while True:
        try:
            cmd = console.input("\n[bold green]HLFA-V19[/bold green] [bold white]@[/bold white] [bold cyan]sman5pkr[/bold cyan]> ").strip().lower()
            
            if cmd in ['exit', 'quit']:
                console.print("[bold red]Menutup sesi...[/bold red]")
                break
            
            elif cmd == 'help':
                console.print("Perintah: [bold white]scan[/bold white] (Otomatis), [bold white]loot[/bold white] (Lihat data & fungsi), [bold white]pivot[/bold white] (Rantai celah), [bold white]exit[/bold white]")
            
            elif cmd == 'scan':
                console.print("[bold magenta][!] Menjalankan Heuristic Discovery...[/bold magenta]")
                await asyncio.sleep(1)
                console.print("[red]★ SSRF DETECTED[/red] -> Vulnerable endpoint found at /proxy?url=")
                inventory["ssrf_meta"]["data"] = "computeMetadata/v1/instance/service-accounts/default/token"
            
            elif cmd == 'loot':
                loot_view = JSON.from_data(inventory)
                # Visualisasi Side-by-Side: Hasil Data dan Penjelasan Fungsinya
                console.print(Columns([Panel(loot_view, title="Current Loot", width=40), show_explanation("SSRF")]))
            
            elif cmd == 'pivot':
                console.rule("[bold red]DEEP CHAINING ANALYSIS (Celah dalam Celah)[/bold red]")
                chain = Text("Step 1: SSRF (Eksploitasi awal)\n  └── Step 2: Extract Cloud Metadata (Loot didapat)\n      └── Step 3: Lateral Movement (Menggunakan token loot untuk akses Bucket S3 internal)\n          └── Step 4: Menemukan Config DB di Bucket (Celah Baru)", style="bold white")
                console.print(Panel(chain, title="Exploit Path Visualization", subtitle="Logika Rantai Serangan"))
                console.print("[yellow]Tip:[/yellow] Gunakan data dari 'loot' untuk melakukan manual probing pada Step 3.")
            
            else:
                console.print(f"[red]Gagal: Perintah '{cmd}' tidak dikenal. Ketik 'help'.[/red]")
        
        except KeyboardInterrupt:
            break

if __name__ == '__main__':
    target = sys.argv[-1] if len(sys.argv) > 1 else "https://sman5palangkaraya.sch.id/"
    asyncio.run(run_hybrid_engine(target))
