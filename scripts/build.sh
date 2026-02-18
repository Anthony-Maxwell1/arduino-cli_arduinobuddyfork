mkdir -p dist
export CGO_LDFLAGS="-Wl,-z,max-page-size=16384 -Wl,-z,common-page-size=16384"
./scripts/buildgomobile.sh ./mobile 25 dist/arduinocli.aar
