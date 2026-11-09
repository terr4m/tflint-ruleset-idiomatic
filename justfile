default:
    just --list

tidy:
    go mod tidy

update:
    go get -u ./...

fmt:
    golangci-lint fmt ./...

lint:
    golangci-lint run --fix ./...

test:
    go test -v ./...

build:
    go build -o ./dist/ -v ./...

md-fmt:
    rumdl fmt --fix .

md-lint:
    rumdl check .

[working-directory("docs")]
docs-build:
    hugo build --gc  --minify --environment production

[working-directory("docs")]
docs-serve:
    hugo serve --buildDrafts --environment development

[working-directory("docs")]
docs-upgrade:
    hugo mod get -u github.com/alex-shpak/hugo-book@latest && hugo mod tidy
