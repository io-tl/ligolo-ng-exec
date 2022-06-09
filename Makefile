#GOO = ~/go/bin/garble
GOO = go

all: gui lin con proxy

gui:
	GOARCH=386 GOOS=windows ${GOO} build -ldflags '-H=windowsgui -extldflags " -s -static"' -o agent32_gui.exe cmd/agent/main.go
	GOOS=windows ${GOO} build -ldflags '-H=windowsgui -extldflags " -s -static"' -o agent_gui.exe cmd/agent/main.go
	GOARCH=386 GOOS=windows ${GOO} build -ldflags '-H=windowsgui -extldflags " -s -static"' -o agent32_svc.exe cmd/agent/svc.go
	GOOS=windows ${GOO} build -ldflags '-H=windowsgui -extldflags " -s -static"' -o agent_svc.exe cmd/agent/svc.go

con:
	GOARCH=386 GOOS=windows ${GOO} build -ldflags '-extldflags " -s -static"' -o agent32_console.exe cmd/agent/main.go
	GOOS=windows ${GOO} build -ldflags '-extldflags " -s -static"' -o agent_console.exe cmd/agent/main.go
lin:
	${GOO} build -ldflags '-extldflags " -s -static"' -o agent cmd/agent/main.go

proxy:
	${GOO} build -ldflags '-extldflags " -s -static"' -o proxy cmd/proxy/main.go


clean :
	rm -f agent  agent32_console.exe  agent32_gui.exe  agent_console.exe  agent_gui.exe  proxy agent32_svc.exe agent_svc.exe

