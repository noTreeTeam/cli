#!/usr/bin/env node

// Ref 1: https://github.com/sanathkr/go-npm
// Ref 2: https://medium.com/xendit-engineering/how-we-repurposed-npm-to-publish-and-distribute-our-go-binaries-for-internal-cli-23981b80911b
"use strict";

import binLinks from "bin-links";
import fs from "fs";
import fetch from "node-fetch";
import { Agent } from "https";
import { HttpsProxyAgent } from "https-proxy-agent";
import path from "path";
import { extract } from "tar";

// Mapping from Node's `process.arch` to Golang's `$GOARCH`
const ARCH_MAPPING = {
  x64: "amd64",
  arm64: "arm64",
};

// Mapping between Node's `process.platform` to Golang's
const PLATFORM_MAPPING = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const arch = ARCH_MAPPING[process.arch];
const platform = PLATFORM_MAPPING[process.platform];

// TODO: import pkg from "../package.json" assert { type: "json" };
const readPackageJson = async () => {
  const contents = await fs.promises.readFile("package.json");
  return JSON.parse(contents);
};

// Build the download url from package.json
const getDownloadUrl = (repo, version) => {
  // Always use 'supabase' as the binary name regardless of package name
  const binaryName = 'supabase';
  const url = `https://github.com/${repo}/releases/download/v${version}/${binaryName}_${platform}_${arch}.tar.gz`;
  return url;
};

const fetchAndParseCheckSumFile = async (repo, version, agent) => {
  // Use 'supabase' as the binary name for checksum file
  const binaryName = 'supabase';
  const checksumFileUrl = `https://github.com/${repo}/releases/download/v${version}/${binaryName}_${version}_checksums.txt`;

  // Fetch the checksum file
  console.info("Downloading", checksumFileUrl);
  const response = await fetch(checksumFileUrl, { agent });
  if (response.ok) {
    const checkSumContent = await response.text();
    const lines = checkSumContent.split("\n");

    const checksums = {};
    for (const line of lines) {
      const [checksum, packageName] = line.split(/\s+/);
      checksums[packageName] = checksum;
    }

    return checksums;
  } else {
    console.error(
      "Could not fetch checksum file",
      response.status,
      response.statusText
    );
  }
};

const errGlobal = `Installing Supabase CLI as a global module is not supported.
Please use one of the supported package managers: https://github.com/noTreeTeam/cli#install-the-cli
`;
const errChecksum = "Checksum mismatch. Downloaded data might be corrupted.";
const errUnsupported = `Installation is not supported for ${process.platform} ${process.arch}`;

/**
 * Reads the configuration from application's package.json,
 * downloads the binary from package url and stores at
 * ./bin in the package's root.
 *
 *  See: https://docs.npmjs.com/files/package.json#bin
 */
async function main() {
  const yarnGlobal = JSON.parse(
    process.env.npm_config_argv || "{}"
  ).original?.includes("global");
  if (process.env.npm_config_global || yarnGlobal) {
    throw errGlobal;
  }
  if (!arch || !platform) {
    throw errUnsupported;
  }

  // Read from package.json and prepare for the installation.
  const pkg = await readPackageJson();
  
  // Use the actual bin key from package.json (should be "supabase")
  const binKey = Object.keys(pkg.bin)[0]; // Get the first (and likely only) bin key
  
  if (platform === "windows") {
    // Update bin path in package.json
    pkg.bin[binKey] += ".exe";
  }

  // Prepare the installation path by creating the directory if it doesn't exist.
  const binPath = pkg.bin[binKey];
  const binDir = path.dirname(binPath);
  await fs.promises.mkdir(binDir, { recursive: true });

  // Create the agent that will be used for all the fetch requests later.
  const proxyUrl =
    process.env.npm_config_https_proxy ||
    process.env.npm_config_http_proxy ||
    process.env.npm_config_proxy;
  // Keeps the TCP connection alive when sending multiple requests
  // Ref: https://github.com/node-fetch/node-fetch/issues/1735
  const agent = proxyUrl
    ? new HttpsProxyAgent(proxyUrl, { keepAlive: true })
    : new Agent({ keepAlive: true });

  // Resolve repository and version from environment or package.json
  const repo = process.env.SUPABASE_REPO || pkg.repository;
  let version = process.env.SUPABASE_VERSION || pkg.version;

  if (version === "latest") {
    const api = `https://api.github.com/repos/${repo}/releases/latest`;
    console.info("Querying", api);
    const resp = await fetch(api, {
      agent,
      headers: { "User-Agent": "supabase-cli" },
    });
    if (!resp.ok) {
      throw new Error(`Failed to get latest release: ${resp.status} ${resp.statusText}`);
    }
    const data = await resp.json();
    version = data.tag_name.replace(/^v/, "");
  }

  // First, fetch the checksum map.
  const checksumMap = await fetchAndParseCheckSumFile(repo, version, agent);

  // Then, download the binary.
  const url = getDownloadUrl(repo, version);
  console.info("Downloading", url);
  const resp = await fetch(url, { agent });
  if (!resp.ok) {
    throw new Error(`Failed to download binary: ${resp.status} ${resp.statusText}`);
  }

  console.info("Download successful, extracting binary...");

  // Download to a temporary file first
  const tempTarPath = path.join(binDir, 'temp_download.tar.gz');
  const tempDir = path.join(binDir, 'temp_extract');
  
  // Ensure directories exist
  await fs.promises.mkdir(binDir, { recursive: true });
  await fs.promises.mkdir(tempDir, { recursive: true });

  // Write the downloaded content to a file
  const buffer = await resp.arrayBuffer();
  await fs.promises.writeFile(tempTarPath, Buffer.from(buffer));
  
  console.info("Downloaded tar file, now extracting...");

  // Extract the tar file
  await extract({
    file: tempTarPath,
    cwd: tempDir,
  });

  console.info("Extraction completed. Looking for binary...");
  
  const binaryName = 'supabase';
  const binName = path.basename(binPath);
  
  // Find and move the binary from the extracted directory structure
  const extractedDirName = `${binaryName}_${platform}_${arch}`;
  const possiblePaths = [
    path.join(tempDir, extractedDirName, binName),  // Standard structure
    path.join(tempDir, binName),  // Flat structure
    path.join(tempDir, extractedDirName, `${binaryName}.exe`),  // Windows with .exe
    path.join(tempDir, `${binaryName}.exe`)  // Flat Windows
  ];

  let foundBinary = false;
  for (const sourcePath of possiblePaths) {
    try {
      await fs.promises.access(sourcePath);
      console.info("Found binary at:", sourcePath);
      await fs.promises.copyFile(sourcePath, binPath);
      console.info("Binary copied to:", binPath);
      foundBinary = true;
      break;
    } catch (error) {
      console.info("Binary not found at:", sourcePath);
    }
  }

  if (!foundBinary) {
    // List what was actually extracted
    console.error("Binary not found! Listing extracted contents:");
    try {
      const walk = async (dir, prefix = '') => {
        const files = await fs.promises.readdir(dir);
        for (const file of files) {
          const fullPath = path.join(dir, file);
          const stat = await fs.promises.stat(fullPath);
          if (stat.isDirectory()) {
            console.error(prefix + file + '/');
            await walk(fullPath, prefix + '  ');
          } else {
            console.error(prefix + file);
          }
        }
      };
      await walk(tempDir);
    } catch (e) {
      console.error("Could not list extracted files:", e.message);
    }
    throw new Error("Could not find extracted binary");
  }

  // Clean up temporary files
  try {
    await fs.promises.unlink(tempTarPath);
    await fs.promises.rm(tempDir, { recursive: true });
    console.info("Cleaned up temporary files");
  } catch (error) {
    console.warn("Could not clean up temporary files:", error.message);
  }

  // Link the binaries in postinstall to support yarn
  await binLinks({
    path: path.resolve("."),
    pkg: pkg,
  });

  console.info("Installed Supabase CLI successfully");
}

await main();
