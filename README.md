
### What is this?
MateSSH is a new SSH server built with memory-safe, secure, lightweight, minimalistic, and modern in mind.

### WARNING
This is still a work in progress project. It is missing important features such as rate limiting and sftp server. Do not use it outside of a development environment.

### Motivations

#### Secure by default
Used in nearly every Linux server environment, openssh-server is distributed by default in nearly every distribution with the system login password set to be used for authentication. The login passwords that are entered frequently on a daily basis are not very strong in most environments. Often, login passwords are set to be extremely weak.

MateSSH does not provide password authentication and enforces public key authentication by default.

#### Easy Setup
You may be concerned that the absence of password authentication makes setup difficult. No need to worry, however.

Upon initial startup, you will be presented with an easy-to-remember and relatively secure one-time passphrase to use for setup. Connect to the server from any SSH client, enter this passphrase, and you will be prompted to enter your SSH public key. This completes the setup. You no longer need to edit sshd_config to enter templated settings or run the ssh-copy-id command.

#### Provides minimal functionality.
This means no plug-in functionality is implemented, reducing the attack surface.

#### Does not run as root.
MateSSH is programmed to refuse to run as root. This cannot be changed.
