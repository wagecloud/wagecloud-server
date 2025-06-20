package nginx

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type ProtocolType string

const (
	http   ProtocolType = "http"
	stream ProtocolType = "stream"
)

func DeleteServerBlock(pathName string, targetPort string) error {
	inputFile, err := os.Open(pathName)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	var output []string
	scanner := bufio.NewScanner(inputFile)

	var (
		inServerBlock = false
		bracketCount  = 0
		blockBuffer   []string
	)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "server {") {
			inServerBlock = true
			bracketCount = 1
			blockBuffer = []string{line}
			continue
		}

		if inServerBlock {
			blockBuffer = append(blockBuffer, line)

			if strings.Contains(line, "{") {
				bracketCount++
			}
			if strings.Contains(line, "}") {
				bracketCount--
				if bracketCount == 0 {
					// Server block ended
					joinedBlock := strings.Join(blockBuffer, "\n")
					if strings.Contains(joinedBlock, "listen "+targetPort+";") {
						// This is the target server block, skip writing
					} else {
						// Not target, keep it
						output = append(output, blockBuffer...)
					}
					inServerBlock = false
					blockBuffer = nil
				}
			}
			continue
		}

		// Not inside server block, keep lines normally
		output = append(output, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Write back to file
	err = os.WriteFile(pathName, []byte(strings.Join(output, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}

type AddOrUpdateServerBlockParams struct {
	PathName     string
	VMIP         string
	InternalPort int
	HostPort     int
	ProtocolType ProtocolType
}

func AddOrUpdateServerBlock(params AddOrUpdateServerBlockParams) error {
	pathName := params.PathName
	vmIP := params.VMIP
	internalPort := params.InternalPort
	hostPort := params.HostPort
	protocolType := params.ProtocolType

	hostPortStr := strconv.Itoa(hostPort)
	// Step 1: Delete old block if exists
	if fileExists(pathName) {
		_ = DeleteServerBlock(pathName, hostPortStr)
	}

	// Step 2: Generate new block
	newBlock, err := generateServerBlock(vmIP, internalPort, hostPort, protocolType)
	if err != nil {
		return err
	}

	// Step 3: Append to file
	f, err := os.OpenFile(pathName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file for append: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString("\n" + newBlock + "\n"); err != nil {
		return fmt.Errorf("error writing new block: %v", err)
	}

	return nil
}

// func main() {
// err := AddOrUpdateServerBlock(filepath.Join(os.Getenv("HOME"), "my-nginx/users.d/stream/test.conf"), "192.168.122.235", 22, 2345, stream)
// 	err := DeleteServerBlock(filepath.Join(os.Getenv("HOME"), "my-nginx/users.d/stream/test.conf"), "2345")
//
// 	if err != nil {
// 		fmt.Println(err)
// 	}
//
// 	err = Reloading()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

func Reloading() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	nginxPrefix := filepath.Join(homeDir, "my-nginx")
	cmd := exec.Command("nginx", "-p", nginxPrefix, "-c", "nginx.conf", "-s", "reload")
	fmt.Println(cmd.String())

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error reloading nginx: %v", err)
	}
	return nil
}

func StartNginx() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	nginxPrefix := filepath.Join(homeDir, "my-nginx")
	cmd := exec.Command("nginx", "-p", nginxPrefix, "-c", "nginx.conf")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error starting nginx: %v", err)
	}

	return nil
}

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func generateServerBlock(vmIP string, internalPort, forwardingPort int, protocolType ProtocolType) (string, error) {
	switch protocolType {
	case "http":
		return fmt.Sprintf(`server {
    listen %d;

    location / {
	proxy_pass http://%s:%d;
    }
}`, forwardingPort, vmIP, internalPort), nil

	case "stream":
		return fmt.Sprintf(`server {
    listen %d;
	proxy_pass %s:%d;
}`, forwardingPort, vmIP, internalPort), nil
	default:
		return "", fmt.Errorf("unsupported protocol type: %s", protocolType)
	}
}
