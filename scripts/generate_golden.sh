#!/bin/bash

# Script to generate or update golden files for CLI testing
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Generating golden files for CLI testing...${NC}"

# Config file containing commands and golden file paths
CONFIG_FILE="tests/golden_commands.txt"

if [ ! -f "$CONFIG_FILE" ]; then
  echo -e "${YELLOW}Error: Configuration file $CONFIG_FILE not found${NC}"
  exit 1
fi

# Create a temporary directory
TEMP_DIR=$(mktemp -d)
echo -e "${GREEN}Created temporary directory:${NC} $TEMP_DIR"

# Clean up temp directory on exit
cleanup() {
  echo -e "${GREEN}Cleaning up temporary directory...${NC}"
  rm -rf "$TEMP_DIR"
}
trap cleanup EXIT

# Build the latest version of the CLI
echo -e "${GREEN}Building the CLI...${NC}"
go build -o ./rules-cli

# Copy the CLI to the temporary directory
cp ./rules-cli "$TEMP_DIR/"
cd "$TEMP_DIR"

echo -e "${GREEN}Working in temporary directory:${NC} $TEMP_DIR"

# Read commands from the config file
while IFS= read -r line || [ -n "$line" ]; do
  # Skip comments and empty lines
  [[ "$line" =~ ^#.*$ || -z "$line" ]] && continue
  
  # Split the line into command and output file
  cmd=$(echo "$line" | cut -d'|' -f1)
  output_file=$(echo "$line" | cut -d'|' -f2)
  
  # Convert output file path to absolute path based on original working directory
  abs_output_file="$(cd - > /dev/null && pwd)/$output_file"
  
  # Ensure the directory exists
  mkdir -p "$(dirname "$abs_output_file")"
  
  echo -e "${GREEN}Generating golden file for:${NC} $cmd"
  echo -e "${GREEN}Output file:${NC} $abs_output_file"
  
  # Execute the command and capture output (including stderr)
  if [ -z "$cmd" ]; then
    # If command is empty, just run the CLI without arguments
    ./rules-cli > "$abs_output_file" 2>&1 || true
  else
    # Otherwise run with the specified arguments
    ./rules-cli $cmd > "$abs_output_file" 2>&1 || true
  fi
done < "$(cd - > /dev/null && pwd)/$CONFIG_FILE"

# Return to original directory
cd - > /dev/null

echo -e "${YELLOW}Golden files generated successfully!${NC}"
echo -e "${YELLOW}Review changes before committing.${NC}"