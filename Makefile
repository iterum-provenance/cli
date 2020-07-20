.PHONY: FORCE

# Build the CLI
build: FORCE 
	go build -o ./build/iterum

# Symlink the CLI, enabling `iterum` from the terminal
link: FORCE 
	@echo "Trying to link the executable to your path:"
	sudo ln -fs "${PWD}/build/iterum" /usr/bin/iterum
	@echo "Use iterum [command] to test the CLI and make clean to remove"

# Remove the created symlink and build folder
clean: FORCE
	sudo rm /usr/bin/iterum
	rm -rf ./build