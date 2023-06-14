

DST:=.
BIN:=easygpt
GOARCH:=amd64
VERSION:=0.1

FLAGS_VERSION=-X main.version=$(VERSION) -X main.revision=$(git rev-parse --short HEAD)
FLAGS=-tags netgo -installsuffix netgo -trimpath "-ldflags=-buildid=" -ldflags '-s -w -extldflags "-static"'
#FLAGS_WIN=-tags netgo -installsuffix netgo -trimpath "-ldflags=-buildid=" -ldflags '-s -w -extldflags "-static" -H windowsgui'
FLAGS_WIN=-tags netgo -installsuffix netgo -trimpath "-ldflags=-buildid=" -ldflags '-s -w -extldflags "-static"'


all:
	make win
	make linux
	make mac
	make license

run:
	go build
	#./automategpt --input-dir ./input --output-dir ./output
	./automategpt --input-dir ./input --output-dir ./output --prompt ""

test:
	make -C internal test

release-no-tag:
	goreleaser release --snapshot --clean

win:
	rm -rf $(DST)/$(BIN).exe
	#GOARCH=$(GOARCH) GOOS=windows go build -o $(DST)/$(BIN)_windows.exe $(FLAGS_WIN) 
	GOARCH=$(GOARCH) GOOS=windows go build -o $(DST)/$(BIN).exe $(FLAGS_WIN) 
	rm -rf $(DST)/$(BIN).upx.exe && upx $(DST)/$(BIN).exe -o $(DST)/$(BIN).upx.exe
	rm -rf $(DST)/$(BIN).exe
	mv $(DST)/$(BIN).upx.exe $(DST)/$(BIN).exe
	#until cp -f $(DST)/$(BIN).exe /mnt/d/public; do sleep 1; done

linux:
	rm -rf $(DST)/$(BIN)
	GOARCH=$(GOARCH) GOOS=linux go build -o $(DST)/$(BIN) $(FLAGS_UNIX) $(FLAGS)
	#rm -rf $(DST)/$(BIN).upx && upx $(DST)/$(BIN) -o $(DST)/$(BIN).upx
	#rm -rf $(DST)/$(BIN)
	#mv $(DST)/$(BIN).upx $(DST)/$(BIN)

mac:
	rm -rf $(DST)/$(BIN)_mac
	GOOS=darwin go build -o $(DST)/$(BIN)_mac $(FLAGS_UNIX) $(FLAGS)
	#rm -rf $(DST)/$(BIN).upx && upx $(DST)/$(BIN) -o $(DST)/$(BIN).upx
	#rm -rf $(DST)/$(BIN)
	#mv $(DST)/$(BIN).upx $(DST)/$(BIN)


pi:
	rm -rf $(DST)/$(BIN)
	GOARM=6  GOARCH=arm  GOOS=linux go build -o $(DST)/$(BIN) $(FLAGS_UNIX) $(FLAGS)
	#rm -rf $(DST)/$(BIN).upx && upx $(DST)/$(BIN) -o $(DST)/$(BIN).upx
	#rm -rf $(DST)/$(BIN)
	#mv $(DST)/$(BIN).upx $(DST)/$(BIN)

upx:
	until sudo apt install upx -y --fix-missing; do sleep 1; done

license:
	rm -rf licenses
	-@go-licenses save . --save_path licenses
	rm -rf licenses/easygpt


install-upx:
	sudo apt install upx
install-go-licenses:
	go install github.com/google/go-licenses@latest
install-commitlint:
	go install github.com/conventionalcommit/commitlint@latest

clean:
	rm easygpt easygpt.exe easygpt_mac
