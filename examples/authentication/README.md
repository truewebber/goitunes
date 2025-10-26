# Authentication Example

This example demonstrates how to authenticate with the App Store using Apple ID credentials.

## Features

- Login with Apple ID and password
- Manage authentication tokens
- Configure device information
- Check authentication status

## Usage

Set environment variables:

```bash
export APPLE_ID="your@email.com"
export APPLE_PASSWORD="your_password"
export DEVICE_GUID="00000000-0000-0000-0000-000000000000"  # Optional
export MACHINE_NAME="MyMachine"  # Optional
```

Run the example:

```bash
cd examples/authentication
go run main.go
```

## Code Explanation

### Creating Client with Credentials

```go
client, err := goitunes.New("us",
    goitunes.WithCredentials(appleID, "", ""),
    goitunes.WithDevice(guid, machineName, goitunes.UserAgentWindows),
)
```

Initialize the client with your Apple ID. The password tokens will be empty initially and filled after login.

### Performing Login

```go
authResp, err := client.Auth().Login(ctx, password)
```

Authenticate with your password. This will:
1. Send authentication request to Apple
2. Receive password token and DSID
3. Update the client's credentials automatically

### Checking Authentication Status

```go
isAuth := client.IsAuthenticated()
canPurchase := client.CanPurchase()
```

Check if the client is authenticated and whether it can make purchases.

## Security Notes

- **Never** hardcode credentials in your code
- Use environment variables or secure credential storage
- Password tokens are temporary and will expire
- Store DSID and password token securely if you want to reuse them

## Device Information

Device information is required for authentication:
- **GUID**: Unique device identifier (UUID format)
- **Machine Name**: Human-readable machine name
- **User Agent**: Browser/iTunes user agent string

These values should remain constant for the same "device" to avoid authentication issues.

## What's Next?

After authentication, you can:
- Use all public API methods
- Purchase applications (requires kbsync certificate)
- Download purchased applications

See the `purchase_app` example for purchasing applications.

