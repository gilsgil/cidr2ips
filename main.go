package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"strings"
)

// extractIPsFromCIDR parses the CIDR and prints each host IP.
func extractIPsFromCIDR(cidr string) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	// Check if IPv4
	if ip4 := ip.To4(); ip4 != nil {
		network := ipNet.IP.To4()
		start := ip4ToUint32(network)
		mask := ip4ToUint32(net.IP(ipNet.Mask).To4())
		broadcast := start | ^mask
		total := broadcast - start + 1
		// For IPv4, exclude network and broadcast addresses if more than two addresses exist.
		if total > 2 {
			for i := start + 1; i < broadcast; i++ {
				fmt.Println(uint32ToIP(i).String())
			}
		}
	} else {
		// IPv6: hosts() returns all addresses except the network address.
		ones, bits := ipNet.Mask.Size()
		total := new(big.Int).Lsh(big.NewInt(1), uint(bits-ones))
		// If total addresses > 1, iterate from 1 to total-1.
		if total.Cmp(big.NewInt(1)) > 0 {
			networkInt := ipToBigInt(ipNet.IP)
			one := big.NewInt(1)
			i := big.NewInt(1)
			for i.Cmp(total) < 0 {
				ipInt := new(big.Int).Add(networkInt, i)
				fmt.Println(bigIntToIP(ipInt, false).String())
				i.Add(i, one)
			}
		}
	}
}

// ip4ToUint32 converts an IPv4 address to a uint32.
func ip4ToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

// uint32ToIP converts a uint32 back to an IPv4 address.
func uint32ToIP(n uint32) net.IP {
	return net.IPv4(byte(n>>24), byte((n>>16)&0xFF), byte((n>>8)&0xFF), byte(n&0xFF))
}

// ipToBigInt converts an IP address to a big.Int.
func ipToBigInt(ip net.IP) *big.Int {
	ip = ip.To16()
	return new(big.Int).SetBytes(ip)
}

// bigIntToIP converts a big.Int to an IP address.
// If isIPv4 is true, it returns an IPv4 address.
func bigIntToIP(i *big.Int, isIPv4 bool) net.IP {
	b := i.Bytes()
	if isIPv4 {
		// Ensure 4 bytes
		if len(b) < 4 {
			padding := make([]byte, 4-len(b))
			b = append(padding, b...)
		}
		return net.IPv4(b[len(b)-4], b[len(b)-3], b[len(b)-2], b[len(b)-1])
	}
	// For IPv6, ensure 16 bytes.
	if len(b) < 16 {
		padding := make([]byte, 16-len(b))
		b = append(padding, b...)
	}
	return net.IP(b)
}

// readCIDRsFromStdinOrFile reads CIDRs from a file (if filePath is provided) or from stdin.
func readCIDRsFromStdinOrFile(filePath string) ([]string, error) {
	var reader io.Reader
	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	scanner := bufio.NewScanner(reader)
	var cidrs []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			cidrs = append(cidrs, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cidrs, nil
}

func main() {
	target := flag.String("t", "", "CIDR to extract IPs (example: 192.168.0.0/24)")
	listFile := flag.String("l", "", "File containing a list of CIDRs")
	flag.Parse()

	var cidrList []string
	var err error

	// If a list file is provided, read CIDRs from it.
	if *listFile != "" {
		cidrList, err = readCIDRsFromStdinOrFile(*listFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Check if data is being piped in via stdin.
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			cidrList, err = readCIDRsFromStdinOrFile("")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
				os.Exit(1)
			}
		} else if *target != "" {
			cidrList = []string{*target}
		} else {
			fmt.Fprintln(os.Stderr, "Error: You must provide a CIDR with -t, or a list of CIDRs with -l or via stdin.")
			os.Exit(1)
		}
	}

	// Process each CIDR.
	for _, cidr := range cidrList {
		extractIPsFromCIDR(cidr)
	}
}
