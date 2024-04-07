
# Philips Hue Bridge Exporter

This Go program discovers Philips Hue bridges on your local network, optionally uses a provided API key to authenticate, and exports detailed information about the lights and rooms (group) into separate CSV files.

## Features

- Discovers Philips Hue bridges within the local network.
- Allows the use of a pre-existing API key or guides through creating a new user on the Hue bridge.
- Fetches and exports detailed information about lights and rooms (groups) into two separate CSV files: `hue_lights_data.csv` and `hue_groups_data.csv`.
- Filters out entertainment areas to avoid listing lights under multiple rooms.

## Prerequisites

Before running the program, ensure you have:

- A Philips Hue bridge connected to your local network and internet.
- Go installed on your system.

## Usage

1. **With Automatic User Creation:**

   First, ensure your Philips Hue bridge is connected and accessible. Then, run the program without any arguments. You'll be prompted to press the link button on your Hue bridge to allow the program to create a new user (API key).

   ```sh
   go run huebridgeexport.go
   ```

   Follow the on-screen instructions to complete the process.

2. **With a Pre-existing API Key:**

   If you already have an API key, you can provide it as a command-line argument to bypass the user creation step.

   ```sh
   go run huebridgeexport.go YOUR_API_KEY_HERE
   ```

## Output

The program will generate two CSV files in the current directory:

- `hue_lights_data.csv`: Contains detailed information about each light, including its name, model ID, type, manufacturer, product name, and group membership.
- `hue_groups_data.csv`: Contains details about each room, including its name and type, excluding any defined as "Entertainment" areas.

## Notes

- The program is designed for use within a local network environment. Ensure your computer and Philips Hue bridge are on the same network.
- API keys are sensitive information. If generated through this program, ensure to keep it secure.

## Troubleshooting

If you encounter issues discovering your Hue bridge or generating an API key, ensure your bridge is powered on and connected to the same local network as your computer. Additionally, verify your internet connection, as bridge discovery relies on Philips' online service.
