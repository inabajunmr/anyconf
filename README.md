# anyconf

anyconf is local configuration file finder.
## Install

```
go get github.com/inabajunmr/anyconf
```
## Usage

When you want to edit config file about AWS.

```
anyconf aws
```

### Example

![](./doc/anyconf.gif)

### Local extension

When you want to add config path to only your local environment, you can use local config file.

```
vim ~/.anyconf/configs.txt
```

When you add line to `~/.anyconf/configs.txt` like following example, you can find yourtool on anyconf.

```
yourtool ~/yourtool/config
```

You can find and edit it via anyconf.

### Editor Configuration

Default editor is vim.
If you want to use vscode and you can launch it by `code` command, you can configure it on `~/.anyconf/configs.yml` like following example.

```
editor: code
```

You can find and edit it via anyconf too.

## Development

### Prerequisite
anyconf use [statik](https://github.com/rakyll/statik) for build process.

### Build

After editing static/configs.txt, following command is needed.
```
statik -src static
```

## Contribution

Please add new config to [static/configs.txt](static/configs.txt).
When you want to add `~/.aws/config`, you can add line like following example.

```diff
aws/credentials ~/.aws/credentials
+ aws/config ~/.aws/config
fish ~/.config/fish/config.fish
git/config ~/.gitconfig
```
