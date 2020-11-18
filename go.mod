module github.com/paketo-buildpacks/google-cloud

go 1.15

require (
	github.com/buildpacks/libcnb v1.18.1
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/onsi/gomega v1.10.3
	github.com/paketo-buildpacks/libpak v1.50.0
	github.com/sclevine/spec v1.4.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/sys v0.0.0-20201117170446-d9b008d0a637 // indirect
)

replace github.com/paketo-buildpacks/libpak => ../../paketo-buildpacks/libpak
