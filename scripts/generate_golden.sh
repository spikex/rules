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

# Build the latest version of the CLI
echo -e "${GREEN}Building the CLI...${NC}"
go build -o ./rules-cli

# Read commands from the config file
while IFS= read -r line || [ -n "$line" ]; do
  # Skip comments and empty lines
  [[ "$line" =~ ^#.*$ || -z "$line" ]] && continue
  
  # Split the line into command and output file
  cmd=$(echo "$line" | cut -d'|' -f1)
  output_file=$(echo "$line" | cut -d'|' -f2)
  
  # Ensure the directory exists
  mkdir -p "$(dirname "$output_file")"
  
  echo -e "${GREEN}Generating golden file for:${NC} $cmd"
  echo -e "${GREEN}Output file:${NC} $output_file"
  
  # Execute the command and capture output (including stderr)
  if [ -z "$cmd" ]; then
    # If command is empty, just run the CLI without arguments
    ./rules-cli > "$output_file" 2>&1 || true
  else
    # Otherwise run with the specified arguments
    ./rules-cli $cmd > "$output_file" 2>&1 || true
  fi
done < "$CONFIG_FILE"

echo -e "${YELLOW}Golden files generated successfully!${NC}"
echo -e "${YELLOW}Review changes before committing.${NC}"