# Security Policy

## Supported Versions

The following table outlines which versions of Fuego are actively supported with security updates. Please ensure that you are using a supported version to benefit from the latest patches and improvements.

| Version                                         | Supported              |
| ----------------------------------------------- | ---------------------- |
| 0.x.y (x being the latest version released)     | :white_check_mark: Yes |
| 0.x.y (x being NOT the latest version released) | :x: No                 |

## Reporting a Vulnerability

Fuego relies on its community to ensure its security. Here is how to report a vulnerability:

1. **Send a Pull Request (PR):** If possible, immediately send a PR addressing the vulnerability and tag the maintainers for a quick review.
2. **Dependency Issues:** For supply chain or dependency-related vulnerabilities, update the all modules with `make check-all-modules` and submit a PR.
3. **Direct Contact:** If you cannot send a PR or the issue requires further discussion, please contact the maintainers directly by email.

### Important Notes

- Please do not publicly disclose the vulnerability until it has been addressed and patched.
- We are committed to transparency and will publicly acknowledge reporters in the release notes unless requested otherwise.

Your cooperation helps ensure Fuego remains a secure and reliable framework for everyone.
