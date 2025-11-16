# Security Policy

## Supported Versions

We release security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x     | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability in Terminal.FM, please report it responsibly.

### How to Report

**Please DO NOT open a public GitHub issue for security vulnerabilities.**

Instead, report security issues via one of these channels:

1. **GitHub Security Advisories** (Preferred)
   - Go to https://github.com/fulgidus/terminal-fm/security/advisories
   - Click "Report a vulnerability"
   - Fill in the details

2. **Email** (Alternative)
   - Send an email to: security@fulgidus.dev
   - Include "Terminal.FM Security" in the subject line
   - Encrypt sensitive information using our PGP key (if available)

### What to Include

When reporting a vulnerability, please include:

- **Description**: Clear description of the vulnerability
- **Impact**: What an attacker could achieve
- **Steps to Reproduce**: Detailed steps to reproduce the issue
- **Proof of Concept**: Code or commands demonstrating the vulnerability
- **Affected Versions**: Which versions are affected
- **Suggested Fix**: If you have ideas for remediation (optional)
- **Your Details**: Name/handle and how you'd like to be credited (optional)

### Example Report

```
Subject: Terminal.FM Security - SSH Authentication Bypass

Description:
An authentication bypass vulnerability exists in the SSH server that
allows unauthorized access without proper credentials.

Impact:
An attacker could gain full access to the Terminal.FM service without
authentication, potentially accessing user bookmarks and preferences.

Steps to Reproduce:
1. Connect to Terminal.FM via SSH
2. Send malformed authentication packet: [packet details]
3. Access granted without valid credentials

Affected Versions:
v1.0.0 to v1.2.3

Suggested Fix:
Implement proper authentication validation in pkg/ssh/auth.go:42
```

## Response Timeline

- **Acknowledgment**: Within 48 hours
- **Initial Assessment**: Within 7 days
- **Status Updates**: Every 14 days until resolved
- **Fix Development**: As soon as possible based on severity
- **Public Disclosure**: After fix is released and deployed

## Security Disclosure Process

1. **Report Received**: We acknowledge your report within 48 hours
2. **Verification**: We verify and assess the vulnerability
3. **Fix Development**: We develop and test a fix
4. **Fix Release**: We release a security update
5. **Public Disclosure**: We publish a security advisory
6. **Credit**: We credit the reporter (if desired)

## Severity Levels

We classify vulnerabilities using the following severity levels:

### Critical
- Remote code execution
- Authentication bypass
- Data breach affecting multiple users
- **Response Time**: Immediate (within 24 hours)
- **Fix Timeline**: 1-3 days

### High
- Privilege escalation
- SQL injection
- Cross-site scripting (XSS) in TUI
- **Response Time**: Within 48 hours
- **Fix Timeline**: 7 days

### Medium
- Denial of service
- Information disclosure
- Insecure defaults
- **Response Time**: Within 7 days
- **Fix Timeline**: 30 days

### Low
- Minor information leaks
- Best practice violations
- **Response Time**: Within 14 days
- **Fix Timeline**: 60 days

## Security Best Practices for Users

### Self-Hosting

If you're self-hosting Terminal.FM:

1. **Keep Updated**: Always run the latest version
   ```bash
   git pull origin main
   go build -o terminal-fm ./cmd/server
   sudo systemctl restart terminal-fm
   ```

2. **Secure SSH**: Configure SSH security properly
   - Use port 2222 for admin access with key-based auth
   - Implement fail2ban for rate limiting
   - Monitor logs regularly

3. **Firewall**: Only expose necessary ports
   ```bash
   sudo ufw allow 22/tcp    # Terminal.FM
   sudo ufw allow 2222/tcp  # Admin SSH
   sudo ufw enable
   ```

4. **Database Security**: Protect your database
   ```bash
   chmod 600 /var/lib/terminal-fm/data.db
   chown terminal-fm:terminal-fm /var/lib/terminal-fm/data.db
   ```

5. **Regular Backups**: Backup your data
   ```bash
   # Daily backups
   0 3 * * * /usr/local/bin/backup-terminal-fm.sh
   ```

### Using Public Instance

When using `ssh terminal.fm`:

1. **Anonymous Mode**: The service is designed for anonymous use
2. **No Sensitive Data**: Don't store sensitive information in bookmarks
3. **SSH Keys**: Use dedicated SSH keys, not your primary keys
4. **Verify Host**: Check SSH fingerprint on first connect

## Known Security Considerations

### By Design

These are not vulnerabilities but design choices:

1. **Anonymous Access**: Terminal.FM allows anonymous SSH connections
   - **Why**: Accessibility and ease of use
   - **Mitigation**: No shell access, isolated environment

2. **Public Bookmarks**: Bookmarks are tied to SSH fingerprint
   - **Why**: Stateless authentication
   - **Mitigation**: No personal information stored

3. **Direct Stream URLs**: Station URLs are sent to users
   - **Why**: Client-side streaming for scalability
   - **Mitigation**: URLs are from trusted Radio Browser API

### Limitations

1. **Rate Limiting**: Connection rate limiting is basic
   - Working on: Advanced rate limiting in v1.5

2. **Input Validation**: Basic validation on user inputs
   - Working on: Comprehensive validation in v1.5

## Security Features

### Current (v1.0)

- **Isolated Environment**: No shell access, forced command only
- **Input Sanitization**: All user inputs are sanitized
- **Secure Defaults**: Security-focused default configuration
- **Audit Logging**: All connections and actions are logged
- **Resource Limits**: CPU and memory limits via systemd
- **HTTPS Only**: All external API calls use HTTPS

### Planned (v1.5+)

- **Rate Limiting**: Advanced per-IP rate limiting
- **Anomaly Detection**: Detect and block suspicious behavior
- **Encrypted Database**: SQLite encryption with SQLCipher
- **Security Headers**: Enhanced security headers
- **Penetration Testing**: Regular third-party security audits

## Security Updates

Security updates are released as patch versions (e.g., v1.0.1, v1.0.2).

### Subscribing to Updates

- **GitHub Watch**: Watch repository for releases
- **RSS Feed**: Subscribe to GitHub releases RSS
- **Mailing List**: (Coming soon)

### Update Notifications

Security advisories are published:
- GitHub Security Advisories
- GitHub Releases with `security` tag
- Project README

## Compliance

Terminal.FM aims to comply with:

- **GDPR**: Minimal data collection, user privacy
- **Security Best Practices**: OWASP guidelines
- **Open Source Security**: OpenSSF Best Practices

## Hall of Fame

We recognize security researchers who responsibly disclose vulnerabilities:

<!-- Will be populated when we receive reports -->

*No vulnerabilities reported yet. Be the first!*

## Contact

- **Security Email**: security@fulgidus.dev
- **GitHub Security**: https://github.com/fulgidus/terminal-fm/security
- **General Issues**: https://github.com/fulgidus/terminal-fm/issues

## Bug Bounty Program

We currently do not have a bug bounty program. However:
- We provide public recognition for valid reports
- We credit researchers in release notes
- We're exploring bug bounty programs for the future

## Additional Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [CWE Top 25](https://cwe.mitre.org/top25/)
- [Go Security Best Practices](https://go.dev/doc/security/best-practices)
- [SSH Security Best Practices](https://www.ssh.com/academy/ssh/security)

---

Thank you for helping keep Terminal.FM secure!

**Last Updated**: 2024-11-16
