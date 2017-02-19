# Discord-bot-for-learning
Discord bot for learning git, go, linux, and to have fun!

## Compile
go get github.com/BlueAgent/Discord-bot-for-learning/hello
* Replace "hello" with the program you want

## Steps to get Go working
Document how I installed Go (I'm still new to unix) with reference to: https://golang.org/doc/install
1. Download latest go from https://golang.org/dl/
	* wget https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz
2. Extract to folder
	* tar -xzf go1.8.linux-amd64.tar.gz
	* mv go ~/go1.8
3. Set environment variables in your profile
	* export GOROOT=$HOME/go1.8
	* export PATH=$PATH:$GOROOT/bin
4. Optional: Configure Go workspace folder (defaults to ~/go or C:/go)
	* export GOPATH=/mnt/c/go
    * might better to create a symlink to /mnt/c/go in the home folder instead
5. Optinal: Add Go compiled programs to the path
	* export PATH=$PATH:$GOPATH/bin
