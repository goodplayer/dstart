# dstart

Start a program as daemon service. for linux &amp; mac , if you don't want to write complex systemd service file.

Will be good when use this tool together with @reboot command in crontab :) .

#### help

```
dstart -h
dstart --help
```

#### usage

```
dstart [options] executable [params...]
```

#### A. Option usage

##### 1. Direct start a program

```
Usage of ./dstart:
  -env value
        -env key=value [-env key=value ...] (default array flags.)
  -err string
        -err error_output_file (currently overwriten file not create a new one)
  -etoo
        -etoo : redirect error output to standard output
  -in string
        -in input_file
  -out string
        -out output_file (currently overwriten file not create a new one)
  -u string
        -u username
  -wd string
        -wd working_directory
```

##### 2. Start programs from a configuration file

```
./dstart -c config.toml
```

Note: The config file can be either found in working directory or application file directory.

See [config.toml.example](config.toml.example) file for config example

#### B. Example

##### tail a file then write to another
```
./dstart -in ~/input_file -out ~/output_file -wd /usr/bin tail
```
