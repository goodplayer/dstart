# dstart

start a program as daemon service. for linux &amp; mac


####help

```
dstart -h
dstart --help
```

####usage

```
dstart [options] executable [params...]
```

####option usage

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

####example

#####tail a file then write to another
```
./dstart -in ~/input_file -out ~/output_file -wd /usr/bin tail
```
