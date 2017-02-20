# Discord-bot-for-learning

Discord bot for learning git, go, linux, vim, and to have fun!

## Compiling and Running

	go get github.com/BlueAgent/Discord-bot-for-learning/discordbot
	discordbot -t <Your bot token>

## How to install Go

Mainly to document how I installed Go (I'm still new to unix) with reference to: https://golang.org/doc/install

1. Download latest go from https://golang.org/dl/
	* `wget https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz`
	* The above is not working for Bash on Windows, it has been fixed already but not released yet.
	* [@robbiev said the 1.6.4 build works, and it does](https://github.com/Microsoft/BashOnWindows/issues/349#issuecomment-275907363)
	* `wget https://storage.googleapis.com/golang/go1.6.4.linux-amd64.tar.gz`
2. Extract to folder
	* `tar -xzf go1.8.linux-amd64.tar.gz`
	* `mv go ~/go1.8`
3. Set environment variables in your profile
	* `export GOROOT=$HOME/go1.8`
	* `export PATH=$PATH:$GOROOT/bin`
4. Optional: Configure Go workspace folder (defaults to ~/go or C:/go)
	* `export GOPATH=/mnt/c/go`
    * Tried using a symlink but it caused problems with Bash on Ubuntu on Windows.
5. Optional: Add Go compiled programs to the path
	* `export PATH=$PATH:$GOPATH/bin`

## Thank you

[discord](https://discordapp.com/developers/docs/intro) for being the best chat and voip for gamers ever.

[discordgo](https://github.com/bwmarrin/discordgo/) for creating bindings for the DiscordAPI and for making examples too.

[golang](https://golang.org/) for being a fun and interesting language.

[stackoverflow](http://stackoverflow.com) community for the great questions and answers.

[vim](http://www.vim.org/) ,vimtutor and [amix' vimrc](https://github.com/amix/vimrc) for making vim really nice to learn and use.

[Bash on Ubuntu on Windows](https://github.com/Microsoft/BashOnWindows) for access to the amazing unix ecosystem without dual-booting or virtualisation.

[git](https://git-scm.com/) and [GitHub](https://github.com/) for keeping code and code history safe.
