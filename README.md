# Gomander

[![release](https://img.shields.io/github/v/release/lazylabz/gomander-app)](https://github.com/lazylabz/gomander-app/releases/latest)
[![license](https://img.shields.io/github/license/lazylabz/gomander-app)](https://github.com/lazylabz/gomander-app/blob/main/LICENSE)
[![codecov](https://codecov.io/gh/lazylabz/gomander-app/branch/main/graph/badge.svg?token=D4LYOARMC5)](https://codecov.io/gh/lazylabz/gomander-app)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/lazylabz/gomander-app)

A simple GUI to launch, monitor and organize your commands. Because nobody should have to juggle 12+ terminal windows just to run their project.

We started building this for ourselves when we got tired of the daily terminal chaos, and figured other developers might find it useful too.

## Main Features

- Keep all your commands organized by project so you never lose track of what belongs where
- Bundle related commands into groups and run them all at once
- See what's running, what's not, and check the logs without switching windows
- Configure working directories however you want - use relative or absolute paths, whatever fits your workflow
- Export and import project configurations to get your whole team on the same page
- Works on macOS, Linux and Windows

---

## Known Issues
### TUI support
Right now commands that use TUI (e.g. ngrok) are not supported, as that require PTY support and the current code relies on commands with stdout

## macOS Users - Important Notice

**⚠️ Code Signing Disclaimer**

Since Gomander is a new project, we haven't purchased the $100 USD Apple Developer certificate required to sign the application. This means macOS will show a warning that the app is "damaged" or from an "unidentified developer."

**To use Gomander on macOS, you'll need to manually remove it from quarantine:**

If you downloaded the .app directly:
```bash
sudo xattr -d com.apple.quarantine /PATH/TO/YOUR/GOMANDER.app
```

If you installed via the DMG (most common):
```bash
sudo xattr -d com.apple.quarantine /Applications/gomander.app
```

This is a one-time process (every time you install a new release 😂). Once you run this command, Gomander will launch normally.

---

## Contributing

We'd love your help making Gomander better! Whether you've found a bug, have an idea for a new feature, or want to contribute code - all contributions are welcome.
- **Found a bug?** Open an issue and let us know what happened
- **Have a feature idea?** We're always curious to hear what would make your workflow smoother
- **Want to contribute code?** Pull requests are always appreciated

Don't hesitate to jump in, even if it's your first contribution to an open source project.

## License

This project is licensed under GPL-3. For more details, check out the [LICENSE](/LICENSE) file.

## Development

### Environment Setup

To get started with development, you'll need to follow the [Wails installation guide](https://wails.io/docs/gettingstarted/installation).

That's it, no extra setup needed beyond what Wails requires.

### Running Dev Mode

Once you've got Wails set up, just run:

```bash
wails dev
```

### Building

Before building, you'll need a couple of prerequisites:
- [create-dmg](https://formulae.brew.sh/formula/create-dmg)
- [makensis](https://formulae.brew.sh/formula/makensis)

To build for all platforms:
```bash
make all
```

Want to see more specific build options? Check them out with:
```bash
make help
```
