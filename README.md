# HAR File Session Cookie Scanner
## Background
This program was developed in response to a [security incident disclosed by Okta](https://krebsonsecurity.com/2023/10/hackers-stole-access-tokens-from-oktas-support-unit/), where an adversary was able to access HAR (HTTP Archive) files shared by Okta's customers with their Customer Support team. These HAR files may have included sensitive session cookies, which could be exploited to hijack user sessions. The incident highlights the importance of scrutinizing the contents of HAR files before sharing them with third parties, even for debugging or customer support purposes.

## Purpose
This tool scans HAR files to identify potential session cookies that may be unsafe to share with third parties. By flagging these cookies, the program aims to prevent the inadvertent sharing of sensitive information.

## Installation
### From Source
No dependencies are needed. Simply clone the repo and compile the `cmd/main.go` file.

### Using `go install`

Run go install: 
```go install github.com/nmelo/har-sanitize/cmd/har-sanitize@latest```

the `har-sanitize` command will be built and installed in your `~go/bin`. Now you can run `har-sanitze <har_file_name>`. 

### Pre-Built Release
Alternatively, you can download a pre-built release from the [Releases](https://github.com/nmelo/ok-to-send-this-har-to-okta/releases) page. Once downloaded, change permissions to executable and run: 

```bash
santize_har <har_file_name>
```

## Usage
Save the code to a file, for example `main.go`.

Place a HAR file named `example.har` in the same directory as `main.go`.

Run the program with the following command:

```bash
go run cmd/main.go example.har
```

The program will scan all the cookies in all the requests contained in the HAR file, flag and scramble potential session cookies that could be risky to share, and then save a sanitized version of the HAR file.

## Example Output

The program will print any risky session cookies: 

```bash
go run cmd/main.go example.har
Unsafe to share: JSESSIONID=CBF969ABF6B1101DC5A9636425425272
Unsafe to share: JSESSIONID=B6DCo89873987234JSDJLHCK32323233
Modified HAR file has been saved as safe_to_share.har
```