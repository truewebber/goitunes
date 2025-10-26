# Purchase Application Example

This example demonstrates how to purchase and download applications from the App Store.

## Prerequisites

You need the following to purchase applications:

1. **Apple ID and Password**: Your Apple account credentials
2. **Kbsync Certificate**: A certificate required for purchases (STDQ operations)
3. **Device Information**: GUID and machine name

⚠️ **Important**: Obtaining a valid kbsync certificate is complex and beyond the scope of this example. This is the main barrier to using the purchase functionality.

## Features

- Authenticate with Apple ID
- Purchase applications
- Get download information (URL, keys, SINF, metadata)
- Prepare for IPA installation

## Usage

Set environment variables:

```bash
export APPLE_ID="your@email.com"
export APPLE_PASSWORD="your_password"
export KBSYNC="your_kbsync_certificate_base64"
export DEVICE_GUID="00000000-0000-0000-0000-000000000000"
export MACHINE_NAME="MyMachine"
```

Run the example:

```bash
cd examples/purchase_app
go run main.go
```

## Code Explanation

### Setup Client with Kbsync

```go
client, err := goitunes.New("us",
    goitunes.WithCredentials(appleID, "", ""),
    goitunes.WithKbsync(kbsync),
    goitunes.WithDevice(guid, machineName, goitunes.UserAgentWindows),
)
```

The kbsync certificate is required for purchase operations.

### Authenticate

```go
authResp, err := client.Auth().Login(ctx, password)
```

Login with your password to get authentication tokens.

### Purchase Application

```go
downloadInfo, err := client.Purchase().Buy(ctx, adamID, versionID)
```

Purchase the application and receive download information.

## Download Information

The `DownloadInfoDTO` contains:

- **URL**: Direct download link for the IPA file
- **DownloadKey**: Key for the download cookie
- **Headers**: Required HTTP headers for downloading
- **Sinf**: DRM information (base64 encoded)
- **Metadata**: iTunes metadata plist (base64 encoded)
- **DownloadID**: Unique download identifier
- **VersionID**: Application version ID
- **FileSize**: IPA file size in bytes

## Installing the Downloaded IPA

After downloading, you need to:

1. Download the IPA using the provided URL and headers
2. Decode the SINF and Metadata from base64
3. Inject `Sinf.sinf` into `SC_Info/` directory in the IPA
4. Inject `iTunesMetadata.plist` to the root of the IPA
5. The IPA can then be installed on a device

## Limitations

- **Re-downloads**: If you already own the app, you might need STDRDL pricing parameter, which requires a different kbsync certificate
- **Kbsync**: Obtaining a valid kbsync is the main challenge
- **Device binding**: The kbsync and device info must match

## Security Notes

- Never share your kbsync certificate
- Keep device information consistent
- Store credentials securely
- Be aware of Apple's terms of service

