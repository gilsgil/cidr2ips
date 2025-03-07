# CIDR2IPs

A command-line utility that converts CIDR notation to individual IP addresses.

## Overview

CIDR2IPs is a Go-based tool that extracts and lists all usable host IP addresses from a given CIDR block. It supports both IPv4 and IPv6 address ranges and can process input from command-line arguments, files, or standard input.

## Features

- Convert CIDR blocks to individual IP addresses
- Support for both IPv4 and IPv6
- Exclude network and broadcast addresses for IPv4 CIDR blocks
- Process multiple CIDR blocks from file or stdin
- Simple command-line interface

## Installation

```bash

# Easy install
go install -v github.com/gilsgil/cidr2ips@latest

# Clone the repository
git clone https://github.com/gilsgil/cidr2ips.git


# Navigate to the repository
cd cidr2ips

# Build the application
go build -o cidr2ips

# Optional: Move to a directory in your PATH
sudo mv cidr2ips /usr/local/bin/
```

## Usage

### Basic Usage with a Single CIDR

```bash
cidr2ips -t 192.168.1.0/30
```

Output:
```
192.168.1.1
192.168.1.2
```

### Process Multiple CIDRs from a File

```bash
cidr2ips -l cidr-list.txt
```

Where `cidr-list.txt` contains:
```
192.168.1.0/30
10.0.0.0/29
```

### Process CIDRs from Stdin (Pipeline)

```bash
echo "192.168.1.0/30" | cidr2ips
```

Or:

```bash
cat cidr-list.txt | cidr2ips
```

## Options

- `-t <cidr>`: Specify a single CIDR block to process
- `-l <file>`: Specify a file containing a list of CIDR blocks, one per line

## Behavior

- For IPv4 CIDR blocks with more than 2 addresses, the network and broadcast addresses are excluded
- For IPv6 CIDR blocks, all host addresses are listed

## Examples

### Extract IPs from a Class C Subnet

```bash
cidr2ips -t 192.168.1.0/24
```

### Process a Small IPv6 Block

```bash
cidr2ips -t 2001:db8::/126
```

### Combining with Other Tools

```bash
cidr2ips -t 10.0.0.0/28 | xargs -I{} nmap -sS {}
```

## License

[MIT License](LICENSE)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
