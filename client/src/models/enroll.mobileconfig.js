export function generateProfile(scepUrl, nanoUrl, deviceId) {
    return `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>PayloadContent</key>
    <array>
        <dict>
            <key>PayloadContent</key>
            <dict>
                <key>Key Type</key>
                <string>RSA</string>
                <key>Challenge</key>
                <string>nanomdm</string>
                <key>Key Usage</key>
                <integer>5</integer>
                <key>Keysize</key>
                <integer>2048</integer>
                <key>URL</key>
                <string>${scepUrl}/scep</string>
            </dict>
            <key>PayloadIdentifier</key>
            <string>com.github.micromdm.scep</string>
            <key>PayloadType</key>
            <string>com.apple.security.scep</string>
            <key>PayloadUUID</key>
            <string>CB90E976-AD44-4B69-8108-8095E6260978</string>
            <key>PayloadVersion</key>
            <integer>1</integer>
        </dict>
        <dict>
            <key>AccessRights</key>
            <integer>8191</integer>
            <key>CheckOutWhenRemoved</key>
            <true/>
            <key>IdentityCertificateUUID</key>
            <string>CB90E976-AD44-4B69-8108-8095E6260978</string>
            <key>PayloadIdentifier</key>
            <string>com.github.micromdm.nanomdm.mdm</string>
            <key>PayloadType</key>
            <string>com.apple.mdm</string>
            <key>PayloadUUID</key>
            <string>96B11019-B54C-49DC-9480-43525834DE7B</string>
            <key>PayloadVersion</key>
            <integer>1</integer>
            <key>ServerCapabilities</key>
            <array>
                <string>com.apple.mdm.per-user-connections</string>
            </array>
            <key>ServerURL</key>
            <string>${nanoUrl}/mdm</string>
            <key>SignMessage</key>
            <true/>
            <key>Topic</key>
            <string>com.apple.mgmt.External.${deviceId}</string>
        </dict>
    </array>
    <key>PayloadDisplayName</key>
    <string>Enrollment Profile</string>
    <key>PayloadIdentifier</key>
    <string>com.github.micromdm.nanomdm</string>
    <key>PayloadScope</key>
    <string>System</string>
    <key>PayloadType</key>
    <string>Configuration</string>
    <key>PayloadUUID</key>
    <string>F9760DD4-F2D1-4F29-8D2C-48D52DD0A9B3</string>
    <key>PayloadVersion</key>
    <integer>1</integer>
</dict>
</plist>`
}
