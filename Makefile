.PHONY: FORCE


build: FORCE 
	go build -o ./build/iterum


link: FORCE 
	@echo "Trying to link the executable to your path:"
	sudo ln -fs "${PWD}/build/iterum" /usr/bin/iterum
	@echo "Use iterum [command] to test the CLI and make clean to remove"

clean: FORCE
	sudo rm /usr/bin/iterum
	
