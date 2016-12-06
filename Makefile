

export GOOS=windows
export GOARCH=amd64

server:;    go build -o output/leo_server.exe ./core/cmd/leo_server/server.go
msys2_starter:; go build -o output/msys2_starter.exe ./core/cmd/msys2_starter/start.go
shell_collector:; go build -o output/shell_collector.exe ./core/cmd/shell_collector/shell_collector.go
inputer:; go build -o output/leo.exe ./core/cmd/inputer/leo.go

build: server msys2_starter shell_collector inputer

install_to_win:; cp output/leo_server.exe output/shell_collector.exe /c/leo/;sc create leo_server binPath= "C:\leo\leo_server.exe" start= auto
