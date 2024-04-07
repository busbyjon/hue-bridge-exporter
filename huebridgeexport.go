package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Bridge Discovery and User Creation
type Bridge struct {
	ID                string `json:"id"`
	InternalIPAddress string `json:"internalipaddress"`
}

type CreateUserRequest struct {
	DeviceType string `json:"devicetype"`
}

type CreateUserResponse struct {
	Success struct {
		Username string `json:"username"`
	} `json:"success"`
	Error struct {
		Description string `json:"description"`
	} `json:"error"`
}

// Hue Resources
type Light struct {
	Name            string `json:"name"`
	ModelID         string `json:"modelid"`
	Type            string `json:"type"`
	Manufacturer    string `json:"manufacturername"`
	ProductName     string `json:"productname"`
	GroupMembership string // Custom field to track group memberships
}

type Group struct {
	Name  string   `json:"name"`
	Type  string   `json:"type"`
	Lights []string `json:"lights"` // List of light IDs belonging to the group
}

func discoverBridges() ([]Bridge, error) {
	resp, err := http.Get("https://discovery.meethue.com/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var bridges []Bridge
	err = json.Unmarshal(body, &bridges)
	return bridges, err
}

func createHueBridgeUser(ipAddress string) (string, error) {
	url := fmt.Sprintf("http://%s/api", ipAddress)
	requestBody := CreateUserRequest{DeviceType: "my_hue_app#go program"}
	jsonValue, _ := json.Marshal(requestBody)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var responses []CreateUserResponse
	json.Unmarshal(body, &responses)

	if len(responses) > 0 {
		if responses[0].Success.Username != "" {
			return responses[0].Success.Username, nil
		} else if responses[0].Error.Description != "" {
			return "", fmt.Errorf(responses[0].Error.Description)
		}
	}

	return "", fmt.Errorf("unknown error occurred")
}

func fetchResource(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func main() {
    var userToken string
    var bridges []Bridge // Declare bridges here to ensure it's accessible throughout the function
    var err error // Declare error here to reuse it throughout the function

    // Check if an API key was provided as a command-line argument
    if len(os.Args) > 1 {
        userToken = os.Args[1]
        fmt.Println("Using provided API key.")
    } else {
        bridges, err = discoverBridges()
        if err != nil || len(bridges) == 0 {
            fmt.Println("No Hue bridges found or there was an error in discovery:", err)
            return
        }

        fmt.Printf("Found Hue Bridge: %s. Please press the link button on the bridge.\n", bridges[0].InternalIPAddress)
        fmt.Println("Press 'Enter' after pressing the link button...")
        fmt.Scanln() // Wait for user to press enter

        userToken, err = createHueBridgeUser(bridges[0].InternalIPAddress)
        if err != nil {
            fmt.Println("Error creating user on the Hue bridge:", err)
            return
        }
        fmt.Println("Successfully created user on the Hue bridge. User Token:", userToken)
    }

    // Ensure bridges were discovered before proceeding (when not using a provided API key)
    if len(bridges) == 0 && userToken != "" {
        bridges, err = discoverBridges()
        if err != nil || len(bridges) == 0 {
            fmt.Println("Failed to discover Hue bridges with the provided API key.", err)
            return
        }
    }
	// Fetch lights and groups
	lightsURL := fmt.Sprintf("http://%s/api/%s/lights", bridges[0].InternalIPAddress, userToken)
	groupsURL := fmt.Sprintf("http://%s/api/%s/groups", bridges[0].InternalIPAddress, userToken)

	lightsJSON, err := fetchResource(lightsURL)
	groupsJSON, err := fetchResource(groupsURL)
	if err != nil {
		fmt.Println("Error fetching Hue data:", err)
		return
	}

	var lights map[string]Light
	var groups map[string]Group
	json.Unmarshal(lightsJSON, &lights)
	json.Unmarshal(groupsJSON, &groups)

	// Map lights to their groups, excluding 'Entertainment' type groups
	for _, group := range groups {
		if group.Type == "Room" { // Only consider groups that are rooms
			for _, lightID := range group.Lights {
				if light, ok := lights[lightID]; ok {
					if light.GroupMembership != "" {
						light.GroupMembership += ", "
					}
					light.GroupMembership += group.Name
					lights[lightID] = light
				}
			}
		}
	}

	// Write lights data to CSV
	lightsFile, err := os.Create("hue_lights_data.csv")
	if err != nil {
		fmt.Println("Error creating lights CSV file:", err)
		return
	}
	defer lightsFile.Close()

	lightsWriter := csv.NewWriter(lightsFile)
	defer lightsWriter.Flush()
	lightsWriter.Write([]string{"ID", "Name", "Model ID", "Type", "Manufacturer", "Product Name", "Group Membership"})

	for id, light := range lights {
		lightsWriter.Write([]string{id, light.Name, light.ModelID, light.Type, light.Manufacturer, light.ProductName, light.GroupMembership})
	}

	// Write groups data to CSV
	groupsFile, err := os.Create("hue_groups_data.csv")
	if err != nil {
		fmt.Println("Error creating groups CSV file:", err)
		return
	}
	defer groupsFile.Close()

	groupsWriter := csv.NewWriter(groupsFile)
	defer groupsWriter.Flush()
	groupsWriter.Write([]string{"ID", "Name", "Type"})

	for id, group := range groups {
		groupsWriter.Write([]string{id, group.Name, group.Type})
	}

	fmt.Println("Data has been successfully written to hue_lights_data.csv and hue_groups_data.csv")
}
