> [!WARNING]
> The tool is actively in development, so expect bugs and missing features. However, it should do its thing and probably wont break your CS2 fonts.  

# Even Better Font Manager
EBFM - Cross-platform TUI tool that makes managing CS2 fonts easy.

## How it works
EBFM basically just edits font config files for you. It takes the font you chose and copies it into `/game/core/panorama/fonts/custom/` and then edits `/game/core/panorama/fonts/conf.d/42-repl-global.conf` to replace all the fonts game uses to your font of choice.  
If you're interested in how cs2 does its font configs, you can check [fonts-conf documentation](https://fontconfig.pages.freedesktop.org/fontconfig/fontconfig-user.html).

## Can i get banned
No, it doesn't inject into the game or change any memory. It only changes publically available files that even say that they can be edited.  
[Unmodified 42-repl-global.conf](https://domitori.xyz/5T0aE.png)

## Installation
Other installation methods will be added after first release is out.

### Building from source

1. Install GoLang -> https://go.dev/doc/install
2. Clone the repository: `git clone https://github.com/dom1torii/even-better-font-manager.git`
4. `cd` into the folder
5. Build the binary: `go build ./cmd/ebfm/`

### Planned features/improvements
- Rework how font collection is stored
