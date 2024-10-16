FROM python:3.9-slim

WORKDIR /app
COPY . .

VOLUME /app/logs:/logs
VOLUME /app/data:/data
VOLUME /app/DemskiyService:/DemskiyService
VOLUME /app/KorolevService:/KorolevService
VOLUME C:\FTP\Archive\OLTC_02aka_ppr02:/app/directories/OLTC_02aka_ppr02
VOLUME C:\FTP\Archive\OLTC_02aka_ppr05:/app/directories/OLTC_02aka_ppr05
VOLUME C:\FTP\Archive\OLTC_13spt_ppr08:/app/directories/OLTC_13spt_ppr08
VOLUME C:\FTP\Archive\OLTC_27moa_ppr34:/app/directories/OLTC_27moa_ppr34
VOLUME C:\FTP\Archive\OLTC_52msa_ppr01:/app/directories/OLTC_52msa_ppr01
VOLUME C:\FTP\Archive\OLTC_52msa_ppr02:/app/directories/OLTC_52msa_ppr02
VOLUME C:\FTP\Archive\OLTC_52msa_ppr03:/app/directories/OLTC_52msa_ppr03
VOLUME C:\FTP\Archive\OLTC_52msa_ppr04:/app/directories/OLTC_52msa_ppr04
VOLUME C:\FTP\Archive\OLTC_52msa_ppr05:/app/directories/OLTC_52msa_ppr05
VOLUME C:\FTP\Archive\OLTC_52msa_ppr06:/app/directories/OLTC_52msa_ppr06
VOLUME C:\FTP\Archive\OLTC_52msa_ppr07:/app/directories/OLTC_52msa_ppr07
VOLUME C:\FTP\Archive\OLTC_52msa_ppr08:/app/directories/OLTC_52msa_ppr08
VOLUME C:\FTP\Archive\OLTC_52msa_ppr09:/app/directories/OLTC_52msa_ppr09
VOLUME C:\FTP\Archive\OLTC_52msa_ppr10:/app/directories/OLTC_52msa_ppr10

CMD ["python", "external_tickets_api.py"]
