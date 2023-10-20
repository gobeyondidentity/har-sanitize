# HAR File Session Cookie Scanner
## Background
This program was developed in response to a [security incident disclosed by Okta](https://krebsonsecurity.com/2023/10/hackers-stole-access-tokens-from-oktas-support-unit/), where an adversary was able to access HAR (HTTP Archive) files shared by Okta's customers with their Customer Support team. These HAR files may have included sensitive session cookies, which could be exploited to hijack user sessions. The incident highlights the importance of scrutinizing the contents of HAR files before sharing them with third parties, even for debugging or customer support purposes.

## Purpose
This tool scans HAR files to identify potential session cookies that may be unsafe to share with third parties. By flagging these cookies, the program aims to prevent the inadvertent sharing of sensitive information.

## Requirements
Go 1.x
Optional: gjson package for extended functionality. To install, run:

```bash
go get -u github.com/tidwall/gjson
```

## Usage
Save the code to a file, for example scan_har.go.

Place a HAR file named `example.har` in the same directory as `main.go`.

Run the program with the following command:

```bash
go run main.go
```

The program will scan all the cookies in all requests contained in the HAR file and flag potential session cookies that could be risky to share.

## Example Output

The program will print any risky session cookies: 

```bash
go run main.go
Unsafe to share: JSESSIONID=CBF969ABF6B1101DC5A9636425425272
Unsafe to share: JSESSIONID=B6DCo89873987234JSDJLHCK32323233
...
```