#!/usr/bin/env node

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');
const os = require('os');

// Define target platforms and architectures
const targets = [
  { platform: 'darwin', arch: 'amd64', goarch: 'amd64' },
  { platform: 'darwin', arch: 'arm64', goarch: 'arm64' },
  { platform: 'linux', arch: 'x64', goarch: 'amd64' },
  { platform: 'linux', arch: 'arm64', goarch: 'arm64' },
  { platform: 'win32', arch: 'x64', goarch: 'amd64' }
];

// Get package version from package.json
const packageJson = JSON.parse(fs.readFileSync('package.json', 'utf8'));
const version = packageJson.version;

// Create bin directory if it doesn't exist
if (!fs.existsSync('bin')) {
  fs.mkdirSync('bin');
}

// Build for each target
for (const target of targets) {
  const { platform, arch, goarch } = target;
  
  // Set binary name based on platform
  const binaryName = platform === 'win32' ? 'rules-cli.exe' : 'rules-cli';
  
  // Create output directory
  const outputDir = path.join('bin', `${platform}-${arch}`);
  if (!fs.existsSync(outputDir)) {
    fs.mkdirSync(outputDir, { recursive: true });
  }
  
  const outputPath = path.join(outputDir, binaryName);
  
  // Set environment variables for cross-compilation
  const env = {
    ...process.env,
    GOOS: platform === 'win32' ? 'windows' : platform,
    GOARCH: goarch,
    CGO_ENABLED: '0'
  };
  
  console.log(`Building for ${platform}-${arch}...`);
  
  try {
    // Build the binary with version info
    execSync(
      `go build -ldflags="-s -w -X main.Version=${version}" -o ${outputPath}`,
      { env, stdio: 'inherit' }
    );
    
    console.log(`Built ${outputPath}`);
    
    // Make binary executable (not needed for Windows)
    if (platform !== 'win32') {
      fs.chmodSync(outputPath, '755');
    }
  } catch (error) {
    console.error(`Failed to build for ${platform}-${arch}: ${error.message}`);
    process.exit(1);
  }
}

console.log('Build completed successfully!');