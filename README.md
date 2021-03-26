# sshmon-check-smtp
Nagios/Checkmk-compatible SSHMon-check for SMTP-Servers

## Installation
* Download [latest Release](https://github.com/indece-official/sshmon-check-smtp/releases/latest)
* Move binary to `/usr/local/bin/sshmon_check_smtp`


## Usage
```
$> sshmon_check_smtp -host mail.testdomain.com -username test@testdomain.com -password mypassword
```

```
Usage of sshmon_check_smtp:
  -dns string
        Use alternate dns server
  -host string
        SMTP-Server Host
  -password string
        Password for SMTP-Server
  -port int
        SMTP-Port (default 465)
  -service string
        Service name (defaults to SMTP_<host>)
  -ssl
        Use SSL for connection (default true)
  -timeout int
        Connect-Timeout in seconds (default 3)
  -username string
        Username for SMTP-Server
  -v    Print the version info and exit
```

Output:
```
0 SMTP_mail.testdomain.com - OK - SMTP-Server on mail.testdomain.com:465 is healthy
```

## Development
### Snapshot build

```
$> make --always-make
```

### Release build

```
$> BUILD_VERSION=1.0.0 make --always-make
```