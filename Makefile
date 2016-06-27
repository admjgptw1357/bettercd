all: bettercd.go
	go build bettercd.go

install: bettercd
	@if [ ! -d ~/.bettercd ]; then \
		mkdir ~/.bettercd; \
	fi
	\cp -f bettercd ~/.bettercd
	\cp -f bettercd.sh ~/.bettercd
	echo "source ~/.bettercd/bettercd.sh" >> ~/.zshrc
	@echo "Plese re-login or source"

update: bettercd
	\cp -f bettercd ~/.bettercd
	\cp -f bettercd.sh ~/.bettercd
