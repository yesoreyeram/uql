# Change Log

Change history of the project. All the feature updates, bug fixes, breaking changes will be documented here.

## [0.0.13]

- Feature: More math functions added. (floor/ceil/round/sign/pow/sin/cos/tan/log/log2/log10)
- Feature: New function `kv` which converts object into array of key value pairs
- Feature: New root level command `distinct` which returns uniq values

## [0.0.12]

- Feature: `parse_url` support added
- Feature: `parse_urlquery` support added
- Bug fix: Fixed a bug by which the `project` command should work correctly on non-array objects

## [0.0.8]

- Feature: new lines starting with `#` will be considered comments
- New command: `mv-expand` added. This expand multi value array into multiple records
- Bug fix: fixed a bug where new lines in windows (`\r\n`) doesn't work

## [0.0.7]

- `parse-yaml` support added

## [0.0.6]

- Summarize by multiple props
- New summarize functions added - `first`, `last`

## [ 0.0.5 ]

- Date utils added (`tounixtime`,`format_datetime`,`add_datetime`,`startofminute`,`startofhour`,`startofday`,`startofmonth`,`startofweek`,`startofyear`)
- Bug fixes

## [ 0.0.4 ]

- Basic parsing functions added via `parse-json`, `parse-csv` and `parse-xml`
- Bug fixes

## [ 0.0.3 ]

- Basic conversion functions added in `project` and `extend` commands. (`toint`, `tolong`, `tobool`, `tostring`,`todouble`, `tofloat`, `todatetime`, `unixtime_seconds_todatetime`, `unixtime_nanoseconds_todatetime`, `unixtime_milliseconds_todatetime` and `unixtime_microseconds_todatetime` )

## [ 0.0.2 ]

- Basic commands added. (`hello`,`ping`,`echo`,`count`,`project`,`project-away`,`limit` and `order by`)

## [ 0.0.1 ]

- Initial placeholder version
