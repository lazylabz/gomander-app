# Gomander

[![codecov](https://codecov.io/gh/lazylabz/gomander-app/branch/main/graph/badge.svg?token=D4LYOARMC5)](https://codecov.io/gh/lazylabz/gomander-app)

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
