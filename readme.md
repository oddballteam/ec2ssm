<details>
<summary>Requirements</summary>

- vpn access
- [AWS Session Manager plugin for AWS CLI](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)

</details>

# Installation
> assuming you have ~/go/bin in your path

```
go install github.com/oddballteam/ec2ssm@latest
```

alternatively you can download a [release](https://github.com/oddballteam/ec2ssm/releases). Unzip the archive and add the ec2ssm binary to a folder in your path.

# Usage
```
ec2ssm
```

run `ec2ssm` and select Web or CLI. Then select an account using the arrow keys.


### Issues
This is a new tool and likely will have issues. 
You can message `@Doug Moore` if you run into an issue.

### Build
builds are automated from GitHub Actions. A token with repo write permission is required for GoRelease to work.

To create a new build:

1. push any code changes
2. run `git tag -a v1.0.1 -m "tag description"` with an appropriate incremented version (e.g. `v1.0.1` -> `v1.0.2`) following major versioning. A `latest` tag will be added for you by GitHub.
3. push the tag `git push origin v1.0.1` which triggers an automated builds
