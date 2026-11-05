set quiet

app := "Twinblox.app"
binary := app / "Contents/MacOS/twinblox"

alias b := build
alias r := run
alias c := check

default:
    just --list

check:
    go build ./... && go vet ./...

build: check
    go build -o {{binary}} .
    codesign --force -s - {{app}}
    touch {{app}}
    echo "built {{app}}"

run: build
    -pkill -x twinblox
    sleep 1
    open {{app}}

test:
    go test $(go list ./internal/... | grep -v /urlhandler)

icon:
    xcrun actool icon/AppIcon.icon --compile {{app}}/Contents/Resources --app-icon AppIcon \\
     --platform macosx --minimum-deployment-target 11.0 --target-device mac \\
     --output-partial-info-plist /tmp/twinblox-icon.plist > /dev/null

    cp icon/MenuBarIcon.png icon/MenuBarIcon@2x.png {{app}}/Contents/Resources/
    just build

install: build
    -pkill -x twinblox
    rm -rf /Applications/{{app}}
    cp -R {{app}} /Applications/
    echo "installed /Applications/{{app}}"
