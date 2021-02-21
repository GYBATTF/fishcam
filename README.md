# fishcam

A simple application that runs in the background regularly setting your background to the current image from [fishcam.com](fishcam.com).

## Build

In the root directory run:

```
make
```

# Usage
Set as an autorun application, and let it run in the background.
Settings are by default stored in `$HOME/.local/etc/fishcam/config.json`.

The default `config.json` looks like this:

```
{
  "picdir": "/home/lana/.local/etc/fishcam",
  "updateinterval": {
    "days": 0,
    "hours": 0,
    "minutes": 30
  }
}
```

Any edits will be read automatically, no need to restart the application.