# Arduino CLI - Ported to android.
Arduino CLI with serial features stripped and replaced with a bridge for any mobile platform (android or ios) to implement through gomobile.
Also exposes interface to access arduino cli features without going through the cli, compressing data formats into dumber formats that gomobile can push through the more primitive formats that java supports compared to go.

> [!WARNING]
> **This is designed to work as a shared library with a serial library binded by a android (or iOS, however harder that may be) app binding them together. This is not standalone.**

> [!NOTE]
> Forked from the Arduino CLI by arduino, view here: **Arduino CLI (original, standalone, desktop):** [here](https://github.com/arduino/arduino-cli)

> [!NOTE]
> This is a component of ArduinoBuddy, the parent project. **ArduinoBuddy project:** [here](https://github.com/Anthony-Maxwell1/ArduinoBuddy)

For more information, architecture and a fuller readme, go to the ArduinoBuddy project above.
