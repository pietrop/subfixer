
# Subfixer

Subfixer is a golang program with minimal dependencies for processing subtitles.
It presently accepts subtitles in only Subrip / SRT format.

It operates in two modes -

 1. **Normal Operation**
In this mode , the program processes all subtitles, either within a range you specify or from start to end. Any subtitle found to not comply with the default or provided parameters is adjusted either in duration or split / joined depending on what is 
At the end all changes are automatically saved back to disk on the same file as input.
There is **NO BACKUP** so you should make sure to take a backup if you are unsure of the parameters or program operation.

 2. **Pefection Check**
In this mode , the program processes all subtitles, either within a range you specify or from start to end. Each perfection check is performed on each subtitle based on the parameters you supply via command line or the defaults for the same.
This can be used in a script also as below -

```bash
./subfixer -file /path/to/file.srt -mode perfection || echo Failed
```
The exit code is 0 in case of no errors. Also the program will say `Perfection check passed succesfully`

## Build

Download the source code and build using go.

```bash
git clone https://github.com/chetan-prime/subfixer.git
cd subfixer
./make
```

You could make a Makefile if you prefer that instead of using a shell script to build

The output is a binary "subfixer" which is self container and should work without any external dependencies like most golang binaries.

## Usage

Some of the parameters are available only in `-mode normal` or `-mode perfection`


```bash
./subfixer -help

Usage of ./subfixer:
  -chars_per_line int
    	Perfection Check - No. of characters/line (default 42)
  -expand_closer_than float
    	Expand two subtitles closer than n seconds (default 0.5)
  -file string
    	Subtitle Input File (Required)
  -forbidden_chars string
    	Perfection Check - Forbidden Characters (default ""{./;/!/?/,:}"")
  -join_shorter_than int
    	Join two lines shorter in length than (default 42)
  -limit_to string
    	Limit to range or list of subtitle id''s (1-2,4-10,14-16,18)
  -line_balance float
    	Perfection Check - Length Balance (%) (default 50)
  -max_lines int
    	Perfection Check - Max. lines (default 2)
  -min_length float
    	Minimum Length for each subtitle (default 1)
  -mode string
    	Operation Mode (normal/perfection) (default "normal")
  -newlines_as_chars
    	Perfection Check - Treat newlines as characters
  -prefer_compact
    	Perfection Check - Prefer Compact Subtitles (default true)
  -reading_speed float
    	Perfection Check - Reading Speed (ch/sec) (default 21)
  -shrink_longer_than float
    	Shrink a single line subtitle longer than n seconds (default 7)
  -spaces_as_chars
    	Perfection Check - Treat Spaces as characters (default true)
  -speed float
    	Desired Characters Per Second (default 21)
  -speed_epsilon float
    	Epsilon in % of Speed value (default 1)
  -split_longer_than float
    	Proportionately split a two line subtitle longer than n seconds (default 7)
  -trim_spaces int
    	Trim space to left & right of each subtitle (default 1)
```

## Dependencies
The program currently includes source for a modified version of [astisub](https://github.com/asticode/go-astisub) . I have removed code for other subtitle formats we don't use and added a new file `subtitles_utils.go` . This contains new helper functions used by subfixer to the existing library.

Also included is strip.go from [html-strip-tags-go](https://github.com/grokify/html-strip-tags-go)

## Contributing
Pull requests will be considered, time providing. For major changes, please open an issue first to discuss what you would like to change.

There are no tests for now.

## License
[AGPL](https://www.gnu.org/licenses/agpl-3.0.en.html)
