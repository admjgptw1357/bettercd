all: bettercd.go
	go build bettercd.go

install: bettercd
	@if [ ! -d ~/.bettercd ]; then \
		mkdir ~/.bettercd; \
	fi
	\cp -f bettercd ~/.bettercd
	\cp -f setting.sh ~/.bettercd
	echo "source ~/.bettercd/setting.sh" >> ~/.zshrc
	echo "source ~/.bettercd/setting.sh" >> ~/.bashrc
	@echo "Plese re-login or source"

update: bettercd
	\cp -f bettercd ~/.bettercd
	\cp -f setting.sh ~/.bettercd
