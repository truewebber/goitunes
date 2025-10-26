# kbsync - Key Bag Synchronization Certificate

## What is kbsync?

**kbsync** is a cryptographic certificate that proves device authenticity to Apple's App Store servers. It's required for purchasing and downloading applications.

Think of it as a "device passport" - it tells Apple "this is a legitimate device with proper authorization."

## Where is it used?

### Authentication Flow (Login)
```
Login Request → kbsync_auth (flag 0xB, no DSID)
└─ Proves device legitimacy during authentication
```

### Purchase Flow (Buy App)
```
Purchase Request → kbsync (flag 1, with DSID)
└─ Proves device has right to purchase applications
```

### Request Example
```xml
<dict>
    <key>kbsync</key>
    <data>BASE64_ENCODED_CERTIFICATE</data>
    <key>guid</key>
    <string>00000000-0000-0000-0000-000000000000</string>
    ...
</dict>
```

## How to obtain kbsync?

### Option 1: Extract from iTunes/App Store (Recommended for testing)

**macOS/Windows:**
1. Install iTunes (legacy version with iTunes Core library)
2. Perform a legitimate purchase
3. Capture the network traffic (e.g., with Wireshark/Charles Proxy)
4. Extract kbsync from the `buyProduct` request body
5. Decode from base64 to verify format
6. Use this kbsync in your application

**Characteristics:**
- ✅ Simple one-time extraction
- ✅ No coding required
- ✅ Works for weeks/years
- ❌ Manual process
- ❌ Tied to specific device GUID

### Option 2: Generate using iTunes Core library

**Requirements:**
- Windows with iTunes 12.x installed
- Access to iTunes Core DLL functions
- GSA (Grand Slam Authentication) services initialized

**Process:**
```cpp
// Pseudo-code from ASOios project
1. Load iTunes Core library
2. Initialize GSA services
3. Perform device provisioning
4. Call kbsync generation function:
   
   GetEncryptKbsyncValue(
       gsaService,
       outputBuffer,
       dsid,           // NULL for login, DSID for purchase
       dsidAsInt64,    // 0 for login, DSID value for purchase
       isAuthLogin     // true for login, false for purchase
   );
```

**Implementation:**
- See [ASOios project](https://github.com/Ilnur-m/ASOios) (C++, Windows)
- Requires reverse engineering iTunes Core
- Complex setup and maintenance

### Option 3: Use external generation service

**Architecture:**
```
Your Go App → HTTP/gRPC → kbsync Generator Service (C++/iTunes Core)
                              ↓
                        Return kbsync
```

**Benefits:**
- ✅ Dynamic generation
- ✅ Separate concerns
- ❌ Additional infrastructure
- ❌ Windows/iTunes dependency

## Using kbsync in goitunes

```go
import "github.com/truewebber/goitunes/pkg/goitunes"

// Provide pre-generated kbsync
client, err := goitunes.New("us",
    goitunes.WithCredentials(appleID, "", ""),
    goitunes.WithKbsync(kbsyncBase64),  // Base64-encoded kbsync
    goitunes.WithDevice(guid, "MyMachine", goitunes.UserAgentWindows),
)

// Now you can purchase apps
downloadInfo, err := client.Purchase().Buy(ctx, adamID, versionID)
```

## Technical Details

### Format
- **Encoding**: Base64
- **Size**: ~800-1200 characters (base64)
- **Content**: Binary cryptographic data

### Types
| Type | Flag | DSID Required | Use Case |
|------|------|---------------|----------|
| kbsync_auth | 0xB | No | Login/Authentication |
| kbsync | 1 | Yes | Purchase operations |

### Generation Algorithm
```
1. Convert DSID to two 32-bit integers
2. Call iTunes Core function with device info
3. Generate cryptographic signature
4. Encode to base64
```

### Lifetime
- **Model**: Until revoked (no fixed TTL)
- **Stability**: Weeks to years
- **Invalidated by**:
  - Password change
  - Device de-authorization
  - Apple server-side revocation
  - Security policy updates

## Limitations

### Current Implementation
- ❌ **Cannot generate dynamically in pure Go** - requires iTunes Core
- ❌ **STDRDL not supported** - re-downloads need different kbsync
- ✅ **Works for STDQ** - first-time purchases

### Device Binding
- kbsync is tied to specific device GUID
- Changing GUID requires new kbsync
- Must keep device info consistent

### Platform Dependency
- Windows: iTunes Core DLL
- macOS: Deprecated in macOS 10.15+ (iTunes split into Music/TV/Podcasts)
- Linux: Not available (no iTunes)

## Security Considerations

### Treat as Secret
```go
// ❌ Don't commit to version control
const kbsync = "ABC123..."  // Bad!

// ✅ Use environment variables
kbsync := os.Getenv("KBSYNC")

// ✅ Or secure vault
kbsync := vault.GetSecret("apple/kbsync")
```

### Best Practices
1. Store kbsync securely (encrypted storage, env vars)
2. Don't share kbsync between different device GUIDs
3. Monitor for invalidation (401/403 errors)
4. Have process to regenerate when needed
5. Keep device info (GUID, machine name) consistent

## Troubleshooting

### "credentials do not support purchasing"
**Cause**: Missing kbsync in client configuration  
**Solution**: Add `goitunes.WithKbsync(kbsyncBase64)`

### "unauthorized" / "invalid kbsync"
**Causes**:
- kbsync for different device GUID
- kbsync type mismatch (auth vs purchase)
- kbsync has been revoked

**Solution**: Extract/generate new kbsync

### "application requires re-download (STDRDL)"
**Cause**: App already owned, needs STDRDL pricing parameter  
**Limitation**: Current implementation only supports STDQ (first purchase)

## References

- [Purchase Parameters Guide](PURCHASE_PARAMETERS.md) - Complete parameter documentation
- [ASOios Project](https://github.com/Ilnur-m/ASOios) - Reference implementation (C++)
- [Examples](examples/purchase_app/) - Working purchase example

## Quick Start

```bash
# 1. Extract kbsync from iTunes (one-time)
# ... use network capture tool ...

# 2. Set environment variable
export KBSYNC="your_base64_encoded_kbsync"
export DEVICE_GUID="00000000-0000-0000-0000-000000000000"

# 3. Use in your code
go run main.go
```

## Summary

| Aspect | Details |
|--------|---------|
| **What** | Cryptographic device certificate |
| **Why** | Proves device authenticity to Apple |
| **When** | Login (kbsync_auth) and Purchase (kbsync) |
| **How** | Extract from iTunes or generate via iTunes Core |
| **Lifetime** | Until revoked (weeks to years) |
| **Challenge** | Cannot generate in pure Go |
| **Solution** | Use pre-generated kbsync (current approach) |

