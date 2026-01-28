set dotenv-load := true

test:
	go test ./...

gover:
	goveralls -repotoken ${GOVERALLS_TOKEN}

updsum SEMVER:
	sleep 3
	curl https://sum.golang.org/lookup/github.com/pyrorhythm/fn@{{SEMVER}}

[parallel]
upload-coverage-and-fetch SEMVER: gover (updsum SEMVER)

tag-push SEMVER:
	git tag {{SEMVER}}
	git push origin {{SEMVER}}

release SEMVER: test (tag-push SEMVER) (upload-coverage-and-fetch SEMVER)
