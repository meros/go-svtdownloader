# go-svtdownloader
A tool to find and download all available episodes of svt play series

Make sure you follow any licenses that regulate how you can use the streams that you download

To install:
```bash
go get github.com/meros/go-svtdownloader/cmd/svtdownloader
```

```bash
usage: svtdownloader --series=SERIES --outDir=OUTDIR [<flags>]

Flags:
      --help           Show context-sensitive help (also try --help-long and --help-man).
  -s, --series=SERIES  Name of series
  -o, --outDir=OUTDIR  Base directory to put files
  -p, --pushbulletToken=PUSHBULLETTOKEN  
                       Pushbullet token for notifications
  -d, --pushbulletDevice=PUSHBULLETDEVICE  
                       Pushbullet device for notifications
  -f, --forever        Keep running forever
```
