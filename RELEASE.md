# Making a release #

Compile and test

Then run

  goreleaser --rm-dist --snapshot

To test the build

When happy, tag the release

  git tag -s v0.0.XX -m "Release v0.0.XX"

Then do a release build (set GITHUB token first)

  goreleaser --rm-dist
