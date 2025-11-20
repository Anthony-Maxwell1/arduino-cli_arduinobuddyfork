## How to integrate Arduino-CLI_Arduinobuddyfork into your own projects
This is a version of the original Arduino CLI designed to run on android devices.
**Important Notes**
- This is designed to be compiled with gomobile's bind function.
- This is compiled into a shared library, and is intended to be used inside an app and not standalone.

### How to compile
A build script is provided for this purpose.

**Prerequisites**: go

Run `build.bat`.
build.bat will run through a repeating build, fixing an issue each time. This should install and configure your gomobile setup to work, and compile. The result will be in dist/arduinocli.aar and arduinocli-sources.jar.