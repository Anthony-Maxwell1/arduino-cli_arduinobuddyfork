set "CGO_LDFLAGS=-Wl,-z,max-page-size=16384 -Wl,-z,common-page-size=16384"
./scripts/buildgomobile.bat ./mobile 25 dist/arduinocli.aar
