package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"
)

type Email struct {
	MessageID 		string   `json:"message_id"`
	Date      		string   `json:"date"`
	From      		string   `json:"from"`
	To        		[]string `json:"to"`
	Subject   		string   `json:"subject"`
	MimeVersion  	string   `json:"mime_version"`
	ContentType  	string   `json:"content_type"`
	ContentTransfer string   `json:"content_transfer"`
	XFrom      		string   `json:"x_from"`
	XTo        		[]string `json:"x_to"`
	XCC        		[]string `json:"x_cc"`
	XBCC       		[]string `json:"x_bcc"`
	XFolder    		string   `json:"x_folder"`
	XOrigin    		string   `json:"x_origin"`
	XFileName  		string   `json:"x_filename"`
	Content      	string   `json:"content"`
}

func main() {
	// Open an HTTP server to generate application profiles
	go func() {
		http.ListenAndServe(":8080", nil)
	}()
	
	// Access command-line arguments
	args := os.Args

	// Check the number of arguments
	if len(args) < 2 {
		fmt.Println("Usage: indexer <database>")
		return
	}

	// Extract the name argument
	dir := args[1]

	// Declare a variable to store user input
	var userEmail string
	var password string
	var stream string
	
	// Prompt the user for input
	fmt.Print("Enter user email: ")
	_, err := fmt.Scan(&userEmail)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Print("Enter password: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Print("Enter desired stream: ")
	_, err = fmt.Scan(&stream)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Generating JSON...")
	
	createJSONUsers(userEmail, password, stream, dir)
	
	fmt.Println("Done!")
}

// Appends data from all users together
func createJSONUsers(userEmail string, password string, stream string, dir string) {
	limit := 1000 // Desired slice limit per JSON file
	var userDocumentsParts [][]Email // Slice to hold the split parts

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			fileName := file.Name()
			subDir := dir+"/"+fileName
			subFiles, err := ioutil.ReadDir(subDir)
			if err != nil {
				log.Fatal(err)
			}

			for _, file := range subFiles {
				if file.IsDir() {
					dirName := file.Name()
					dirPath := filepath.Join(subDir, dirName)
					userDocuments := getUserDocuments(dirPath)
					partNumber := 0
					for i := 0; i < len(userDocuments); i += limit {
						partNumber++
						
						// Calculate the end index for the current split part
						end := i + limit
						if end > len(userDocuments) {
							end = len(userDocuments)
						}
						
						// Slice the array to create a new split part
						part := userDocuments[i:end]
						userDocumentsParts = append(userDocumentsParts, part)
		
						// Converts data to json
						jsonData, err := json.Marshal(part)
						if err != nil {
							log.Fatal(err)
						}
		
						// Add data to zincsearch
						req, err := http.NewRequest("POST", "http://localhost:5080/api/default/"+stream+"/_json", strings.NewReader(string(jsonData)))
						if err != nil {
							log.Fatal(err)
						}
						req.SetBasicAuth(userEmail, password)
						req.Header.Set("Content-Type", "application/json")
						req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
			
						resp, err := http.DefaultClient.Do(req)
						if err != nil {
							log.Fatal(err)
						}
						defer resp.Body.Close()
						log.Println(resp.StatusCode)
						body, err := io.ReadAll(resp.Body)
						if err != nil {
							log.Fatal(err)
						}
						fmt.Println(string(body))
		
						if err != nil {
							fmt.Println("Error reading file:", err)
							continue
						}
					}
					partNumber = 0;	
				}
			}
		}
	}
}

// Gets data from all documents inside a users directory
func getUserDocuments(dir string) []Email {
	emails := []Email{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		defer file.Close() // Close the file
		if err != nil {
			return err
		}
		content, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		email := parseEmail(string(content))
		emails = append(emails, email)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return emails
}

// Parses the email contents together in an object
func parseEmail(fileContent string) Email {
	email := Email{}

	email.MessageID = getStringValueByPrefix(fileContent, "Message-ID:")
	email.Date = getStringValueByPrefix(fileContent, "Date:")
	email.From = getStringValueByPrefix(fileContent, "From:")
	email.To = getArrayValueByPrefix(fileContent, "To:")
	email.Subject = getStringValueByPrefix(fileContent, "Subject:")
	email.MimeVersion = getStringValueByPrefix(fileContent, "Mime-Version:")
	email.ContentType = getStringValueByPrefix(fileContent, "Content-Type:")
	email.ContentTransfer = getStringValueByPrefix(fileContent, "Content-Transer-Encoding:")
	email.XFrom = getStringValueByPrefix(fileContent, "X-From:")
	email.XTo = getArrayValueByPrefix(fileContent, "X-To:")
	email.XCC = getArrayValueByPrefix(fileContent, "X-cc:")
	email.XBCC = getArrayValueByPrefix(fileContent, "X-bcc:")
	email.XFolder = getStringValueByPrefix(fileContent, "X-Folder:")
	email.XOrigin = getStringValueByPrefix(fileContent, "X-Origin:")
	email.XFileName = getStringValueByPrefix(fileContent, "X-FileName:")
	email.Content = getContent(fileContent)

    return email
}


// Extracts the value for strings
func getStringValueByPrefix(contentFile string, prefix string) string {
    scanner := bufio.NewScanner(strings.NewReader(contentFile))
    var value string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, prefix) {
			value = strings.TrimSpace(line[len(prefix):])
			for scanner.Scan() {
				nextLine := scanner.Text()
				if strings.HasPrefix(nextLine, " ") || strings.HasPrefix(nextLine, "\t") {
					value += nextLine
				} else {
					break
				}
			}
			break
		}
	}
    return value
}


// Extracts the values for arrays
func getArrayValueByPrefix(contentFile string, prefix string) []string {
    scanner := bufio.NewScanner(strings.NewReader(contentFile))
    var values []string
	var currentValue string
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, prefix) {
            currentValue = strings.TrimSpace(line[len(prefix):])
            for scanner.Scan() {
                nextLine := scanner.Text()
                if strings.HasPrefix(nextLine, " ") || strings.HasPrefix(nextLine, "\t") {
                    currentValue += nextLine
                } else {
                    break
                }
            }
            values = append(values, currentValue)
            currentValue = ""
			break
        }
    }
    return values
}

// Extracts the email content
func getContent(contentFile string) string {
	scanner := bufio.NewScanner(strings.NewReader(contentFile))
	var content string
	for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "X-FileName:") {
            for scanner.Scan() {
                content += scanner.Text() + "\n"
            }
            break
        }
    }
	return content
}