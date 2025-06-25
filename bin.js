#!/usr/bin/env node

const { spawn } = require("child_process");
const path = require("path");
const fs = require("fs");
const os = require("os");

// Determine platform and architecture
const platform = os.platform();
let arch = os.arch();

// Map architecture for consistency with binary naming
if (platform === "darwin" && arch === "x64") {
  arch = "amd64";
}

// Map OS and architecture to binary name
let binaryName;
if (platform === "win32") {
  binaryName = "rules-cli.exe";
} else {
  binaryName = "rules-cli";
}

// Construct path to the binary
const binaryPath = path.join(
  __dirname,
  "bin",
  `${platform}-${arch}`,
  binaryName
);

// Check if binary exists
if (!fs.existsSync(binaryPath)) {
  console.error(`Binary not found for your platform (${platform}-${arch})`);
  process.exit(1);
}

// Make binary executable (not needed on Windows)
if (platform !== "win32") {
  try {
    fs.chmodSync(binaryPath, "755");
  } catch (error) {
    console.error(`Failed to make binary executable: ${error.message}`);
    process.exit(1);
  }
}

// Execute the binary with all arguments passed through
const childProcess = spawn(binaryPath, process.argv.slice(2), {
  stdio: "inherit",
});

childProcess.on("error", (error) => {
  console.error(`Failed to start binary: ${error.message}`);
  process.exit(1);
});

childProcess.on("close", (code) => {
  process.exit(code);
});
