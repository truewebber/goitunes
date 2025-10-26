# Purchase Parameters Guide

This document describes all parameters required for purchasing applications from the Apple App Store.

## Overview

Purchasing an app requires multiple parameters that authenticate the user, identify the device, and prove the request is legitimate. These parameters are sent in the buy product request.

## Required Parameters

### 1. Authentication Parameters

#### `xToken` (Password Token)
- **Type**: String (query parameter)
- **Description**: Authentication token obtained after successful login
- **Example**: `"AgAAAABVxM9NAAABVxM9NQAAAAFXEM9N..."`
- **How to get**: Call `client.Auth().Login(ctx, password)` → `authResp.PasswordToken`
- **Lifetime**: Session-based, expires after logout or timeout

#### `X-Dsid` (Directory Services ID)
- **Type**: String (HTTP header)
- **Description**: Unique identifier for the Apple account
- **Example**: `"123456789"`
- **How to get**: From login response → `authResp.DSID`
- **Note**: Required for most authenticated operations

### 2. Device Parameters

#### `guid` (Device GUID)
- **Type**: String (in request body)
- **Description**: Globally Unique Identifier for the device
- **Example**: `"00000000-0000-0000-0000-000000000000"`
- **Format**: UUID v4 format
- **How to set**: `goitunes.WithDevice(guid, machineName, userAgent)`
- **Note**: Should remain consistent for the same "device"

#### `machineName`
- **Type**: String (in request body)
- **Description**: Human-readable name for the machine
- **Example**: `"MyMachine"`, `"Johns-MacBook-Pro"`
- **How to set**: Via `goitunes.WithDevice(guid, machineName, userAgent)`

#### User Agent
- **Type**: String (HTTP header)
- **Description**: Identifies the client application
- **Examples**:
  - Windows: `"iTunes/12.10.8 (Windows; Microsoft Windows 10 x64 Professional Edition (Build 19041); x64) AppleWebKit/7608.1017.0.25.4"`
  - macOS: `"iTunes/12.10.8 (Macintosh; OS X 10.15.6) AppleWebKit/605.1.15"`
- **How to set**: `goitunes.WithDevice(guid, machineName, goitunes.UserAgentWindows)`
- **Available constants**: `UserAgentWindows`, `UserAgentMacOS`

### 3. The Critical Parameter: kbsync

#### `kbsync`
- **Type**: Base64-encoded binary data
- **Description**: Cryptographic certificate proving device legitimacy
- **Size**: ~800-1200 characters (base64)
- **Format**: `<key>kbsync</key><data>BASE64_DATA</data>`
- **How to set**: `goitunes.WithKbsync(kbsyncBase64)`
- **Details**: See [kbsync Guide](KBSYNC.md) for complete information

**What is kbsync?**
- Cryptographic signature tied to device and account
- Generated using iTunes Core library functions
- Different for login (`flag=0xB`) vs purchase (`flag=1`)
- Requires DSID for purchase operations

**Two types:**
```
1. kbsync_auth (for login)
   - Generated with flag 0xB
   - No DSID required
   - Used in authentication flow

2. kbsync (for purchase)
   - Generated with flag 1
   - Requires valid DSID
   - Used in buy operations
```

**How to obtain:**
- Option A: Use external service (iTunes Core wrapper)
- Option B: Extract from successful iTunes/App Store purchase

**Lifetime:**
- No official TTL documented
- Observed to be stable for extended periods (weeks to years)
- Validity model: **"Until revoked"** rather than fixed expiration
- Can be invalidated by:
  - Account password change
  - Device de-authorization in iTunes/App Store
  - Apple server-side revocation
  - Security policy changes

**Limitations:**
- Tied to specific device GUID
- May need different certificate for re-downloads (STDRDL)

### 4. Application Parameters

#### `salableAdamId` (Adam ID)
- **Type**: String/Integer
- **Description**: Unique identifier for the app in App Store
- **Example**: `"284882215"` (Facebook app)
- **How to get**: From app info endpoint or App Store URL
- **Note**: Also called "trackId" in some contexts

#### `appExtVrsId` (Version ID)
- **Type**: Integer
- **Description**: Specific version identifier for the app
- **Example**: `850109516`
- **How to get**: From `client.Applications().GetByBundleID()` → `app.VersionID`
- **Note**: Required to purchase specific version

#### `price`
- **Type**: String
- **Description**: Price of the application
- **Example**: `"0"` (free), `"0.99"`, `"4.99"`
- **Format**: Decimal string
- **Note**: Used for display and validation

#### `pricingParameters`
- **Type**: String
- **Description**: Type of purchase operation
- **Values**:
  - `"STDQ"` - Standard purchase (first-time buy)
  - `"STDRDL"` - Standard re-download (already owned)
- **Current support**: Only `STDQ` is supported
- **Note**: STDRDL requires different kbsync certificate

### 5. Metadata Parameters

#### `mtClientId`
- **Type**: String
- **Description**: Metrics client identifier
- **Example**: `"3z30dhYIz29Wz4gvz9AEz1NIUDKelm"`
- **Generated**: Random identifier for tracking
- **Note**: Used for Apple's analytics

#### `mtEventTime`
- **Type**: String (Unix timestamp in milliseconds)
- **Description**: Event timestamp
- **Example**: `"1706789123456"`
- **Format**: 13-digit Unix timestamp (ms)

#### `mtRequestId`
- **Type**: String
- **Description**: Unique request identifier
- **Example**: `"3z30dhYIz29Wz4gvz9AEz1NIUDKelmzJ4H6DIUSz1HZC"`
- **Generated**: Random identifier

#### `mtPageContext`, `mtPageType`, `mtTopic`
- **Type**: String
- **Description**: Metadata about user's navigation context
- **Default values**:
  - `mtPageContext`: `"App Store"`
  - `mtPageType`: `"Software"`
  - `mtTopic`: `"xp_its_main"`

### 6. Store Parameters

#### `X-Apple-Store-Front`
- **Type**: String (HTTP header)
- **Description**: Store region identifier
- **Example**: `"143441-1,29"` (US store)
- **Format**: `"{storeFrontId}-{version},{language}"`
- **How to set**: Via region in `goitunes.New("us", ...)`

#### `X-Apple-Tz`
- **Type**: String (HTTP header)  
- **Description**: Timezone offset
- **Example**: `"7200"` (UTC+2)
- **Default**: `"0"` (UTC)

### 7. Additional Headers

#### `Content-Type`
- **Value**: `"application/x-apple-plist"`
- **Description**: Indicates plist format for request body

#### `Referer`
- **Value**: `"http://itunes.apple.com/app/id{adamId}"`
- **Description**: Referrer URL showing source page

## Complete Request Example

### HTTP Request
```http
POST /WebObjects/MZBuy.woa/wa/buyProduct?xToken=PASSWORD_TOKEN HTTP/1.1
Host: buy.itunes.apple.com
Content-Type: application/x-apple-plist
User-Agent: iTunes/12.10.8 (Windows; ...)
X-Apple-Store-Front: 143441-1,29
X-Apple-Tz: 0
X-Dsid: 123456789
X-Token: PASSWORD_TOKEN
Referer: http://itunes.apple.com/app/id284882215

<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
<dict>
	<key>appExtVrsId</key><string>850109516</string>
	<key>guid</key><string>00000000-0000-0000-0000-000000000000</string>
	<key>kbsync</key><data>BASE64_ENCODED_KBSYNC_DATA</data>
	<key>machineName</key><string>MyMachine</string>
	<key>mtApp</key><string>com.apple.iTunes</string>
	<key>mtClientId</key><string>3z30dhYIz29Wz4gvz9AEz1NIUDKelm</string>
	<key>mtEventTime</key><string>1706789123456</string>
	<key>mtPageContext</key><string>App Store</string>
	<key>mtPageId</key><string>1140828062</string>
	<key>mtPageType</key><string>Software</string>
	<key>mtPrevPage</key><string>Genre_134583</string>
	<key>mtRequestId</key><string>3z30dhYIz29Wz4gvz9AEz1NIUDKelmzJ4H6DIUSz1HZC</string>
	<key>mtTopic</key><string>xp_its_main</string>
	<key>needDiv</key><string>0</string>
	<key>pg</key><string>default</string>
	<key>price</key><string>0</string>
	<key>pricingParameters</key><string>STDQ</string>
	<key>rebuy</key><string>false</string>
	<key>productType</key><string>C</string>
	<key>salableAdamId</key><string>284882215</string>
	<key>uuid</key><string>353F3F00-9D87-5BB1-9055-B7761CCD57AA</string>
</dict>
</plist>
```

### Go Code Example
```go
// 1. Create client with all required parameters
client, err := goitunes.New("us",
    goitunes.WithCredentials(appleID, "", ""),
    goitunes.WithKbsync(kbsyncBase64),
    goitunes.WithDevice(deviceGUID, machineName, goitunes.UserAgentWindows),
)

// 2. Authenticate
authResp, err := client.Auth().Login(ctx, password)
// Now have: authResp.DSID, authResp.PasswordToken

// 3. Get app information
apps, err := client.Applications().GetByBundleID(ctx, "com.facebook.Facebook")
app := apps[0]
// Now have: app.AdamID, app.VersionID, app.Price

// 4. Purchase
downloadInfo, err := client.Purchase().Buy(ctx, app.AdamID, app.VersionID)
```

## Parameter Dependencies

```
Purchase Flow:
┌─────────────────────────────────────────────────────┐
│ 1. Device Setup (before any requests)              │
│    ├─ GUID (device identifier)                     │
│    ├─ Machine Name                                 │
│    ├─ User Agent                                   │
│    └─ kbsync (pre-generated)                       │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│ 2. Authentication (Login)                          │
│    Input:  AppleID, Password                       │
│    Output: DSID, PasswordToken                     │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│ 3. App Discovery (GetByBundleID)                  │
│    Input:  Bundle ID                               │
│    Output: AdamID, VersionID, Price                │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│ 4. Purchase (Buy)                                  │
│    Required:                                       │
│    ├─ AdamID (from step 3)                        │
│    ├─ VersionID (from step 3)                     │
│    ├─ DSID (from step 2)                          │
│    ├─ PasswordToken (from step 2)                 │
│    ├─ kbsync (from step 1)                        │
│    ├─ GUID (from step 1)                          │
│    └─ Store Front (from client initialization)    │
│                                                     │
│    Output: DownloadInfo (URL, keys, metadata)     │
└─────────────────────────────────────────────────────┘
```

## Parameter Validation

### Client-side validation
The library performs these checks:
- `adamID`: Must not be empty
- `versionID`: Must be positive integer
- `kbsync`: Must be set (via `WithKbsync()`)
- `credentials`: Must have DSID and PasswordToken (from login)

### Server-side validation
Apple servers verify:
- `kbsync` signature is valid and matches device
- `DSID` belongs to authenticated account
- `PasswordToken` is valid and not expired
- `versionID` matches the `adamID`
- Account has permission to purchase (payment method, etc.)
- No rate limiting violations

## Common Issues

### 1. "credentials do not support purchasing"
```go
// Problem: Missing kbsync
client, _ := goitunes.New("us",
    goitunes.WithCredentials(appleID, "", ""),
    // ❌ Missing: goitunes.WithKbsync(kbsync)
)

// Solution: Add kbsync
client, _ := goitunes.New("us",
    goitunes.WithCredentials(appleID, "", ""),
    goitunes.WithKbsync(kbsyncBase64), // ✅
)
```

### 2. "application requires re-download (STDRDL)"
```go
// Problem: App already owned, needs STDRDL
// Error: MZCommerceSoftware.OwnsSupersededMinorSoftwareApplicationForUpdate

// Current limitation: STDRDL requires different kbsync certificate
// Not supported in current implementation
```

### 3. "unauthorized" / "invalid credentials"
```
Causes:
- PasswordToken expired (re-login required)
- DSID doesn't match account
- kbsync doesn't match device GUID
- Store region mismatch
```

### 4. "invalid kbsync"
```
Causes:
- kbsync is for different device GUID
- kbsync was generated for login (flag=0xB) not purchase (flag=1)
- kbsync expired or invalidated
- Wrong base64 encoding
```

## Parameter Generation Reference

For those interested in generating parameters dynamically:

### kbsync Generation (Requires iTunes Core)
```cpp
// From ASOios project (Windows + iTunes Core)
BOOL GetEncryptKbsyncValue(
    IN LPVOID lpGsaServices,      // GSA service object
    LPVOID lpBuffer,               // Output buffer
    char* lpDsid,                  // DSID string (or NULL for login)
    LONGLONG dsid,                 // DSID as int64 (or 0 for login)
    bool bAuthLogin                // true=login(0xB), false=purchase(1)
);
```

### Device GUID Generation
```go
// Simple approach: Generate once and reuse
import "github.com/google/uuid"

guid := uuid.New().String()
// Store this and use consistently for the same "device"
```

### Timestamps
```go
// mtEventTime: milliseconds since epoch
timestamp := time.Now().UnixNano() / 1000000
```

### Random IDs
```go
// mtClientId, mtRequestId: random alphanumeric strings
func generateRandomID(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    // ... random generation
}
```

## Security Considerations

### Sensitive Parameters
These parameters should be treated as secrets:
- ❌ **kbsync**: Device-specific certificate
- ❌ **PasswordToken**: Session authentication
- ❌ **DSID**: Account identifier
- ⚠️ **GUID**: Device identifier (less sensitive but should be consistent)

### Best Practices
1. **Store kbsync securely**: Use environment variables or secure vaults
2. **Rotate credentials**: Re-login periodically to get fresh PasswordToken
3. **Consistent device info**: Keep GUID and machineName stable
4. **Handle errors gracefully**: Don't expose sensitive data in error messages
5. **Rate limiting**: Don't spam purchase requests

## References

- [Example: Purchase App](../examples/purchase_app/README.md)
- [Client Options](../pkg/goitunes/options.go)
- [Purchase Service](../pkg/goitunes/service_purchase.go)
- [Constants (Store Fronts, User Agents)](../pkg/goitunes/constants.go)

## See Also

- `README.md` - Main documentation
- `examples/purchase_app/` - Complete working example
- `examples/authentication/` - Authentication flow example
- Technical analysis: Research notes on kbsync (removed for brevity)

