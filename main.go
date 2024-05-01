package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

// const filename = "C:\\Windows\\system32\\drivers\\etc\\hosts"
const filename = "hosts"
const searchTextStart = "# Custom Subdomains"
const searchTextEnd = "# END: Custom Subdomains"

func main() {
	clearScreen()

	for {
		fmt.Println("Please select an option:")
		displayListLine("Add subdomain", 1)
		displayListLine("List subdomains", 2)
		displayListLine("Delete subdomain", 3)
		displayListLine("Exit", 4)
		fmt.Print(">> ")

		option := getUserInput()

		switch option {
		case "1":
			addSubdomain()
		case "2":
			listSubdomains()
		case "3":
			deleteSubdomain()
		default:
			fmt.Println("Invalid option selected. Please select a valid option.")
		}

		restart()
	}
}

func listSubdomains() {
	lines, err := readFileLines(filename)
	customDomainNumber := 1
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	inCustomSubDomainsSection := false
	for _, line := range lines {
		if line == searchTextStart {
			inCustomSubDomainsSection = true
			continue
		} else if line == searchTextEnd {
			break
		}
		if inCustomSubDomainsSection {
			displayListLine(line, customDomainNumber)
			customDomainNumber++
		}
	}

	pauseBeforeExist()
}

func deleteSubdomain() {

}

func displayListLine(lineText string, lineNumber int) {
	lineTextTrimmed := strings.TrimSpace(lineText)

	terminalWidth, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Error getting terminal width:", err)
		return
	}

	remainingWidth := terminalWidth - len(lineTextTrimmed) - 2 - len(fmt.Sprintf("%d", lineNumber))
	if remainingWidth < 0 {
		remainingWidth = 0
	}

	dots := strings.Repeat(".", remainingWidth)

	fmt.Printf("%s %s %d", lineTextTrimmed, dots, lineNumber)
}

func addSubdomain() {
	// Prompt the user to enter a domain name
	fmt.Print("Enter the domain name: ")
	domain := getUserInput()

	// Prompt the user to enter an IP address
	fmt.Print("Enter the IP address (default is 127.0.0.1): ")
	ip := getUserInput()

	if ip == "" {
		ip = "127.0.0.1"
	}

	replaceText := fmt.Sprintf("%s %s", ip, domain)

	if textExistsInFile(filename, replaceText) {
		fmt.Println("The domain already exists.")
		fmt.Println()
		fmt.Println()
		fmt.Print("Press the Enter Key to terminate the console screen!")
		getUserInput()
		return
	}

	lines, err := readFileLines(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	updatedLines := insertTextAfter(lines, searchTextStart, replaceText)

	err = writeLinesToFile(filename, updatedLines)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("File updated successfully!")
}

func pauseBeforeExist() {
	fmt.Println()
	fmt.Println()
	fmt.Print("Press the Enter Key to terminate the console screen!")
	getUserInput()
}

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	// Remove newline character from the input
	return strings.TrimSpace(input)
}

func readFileLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func insertTextAfter(lines []string, searchText, newText string) []string {
	var updatedLines []string
	for _, line := range lines {
		updatedLines = append(updatedLines, line)
		if line == searchText {
			updatedLines = append(updatedLines, newText)
		}
	}
	return updatedLines
}

func writeLinesToFile(filename string, lines []string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	return writer.Flush()
}

func textExistsInFile(filename, searchText string) bool {
	lines, err := readFileLines(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return false
	}

	for _, line := range lines {
		if line == searchText {
			return true
		}
	}
	return false
}

func restart() {
	clearScreen()

	switch runtime.GOOS {
	case "windows":
		exec.Command("cmd", "/c", "cls").Run()
	case "linux", "darwin":
		exec.Command("clear").Run()
	default:
		panic("Unsupported OS")
	}

	// Close any open files or resources if necessary
	// Optionally, you can reset any global state

	// Execute the program again with the same arguments
	args := os.Args
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func clearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
