
ahfs.exe: 
	go build -o $@
	mv --force $@ bin/$@
	xcopy templates bin\templates\ /q /Y /t