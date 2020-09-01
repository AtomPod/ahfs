
ahfs.exe: 
	go build -o $@
	mv --force $@ bin/$@
