# Utilities

There are a few scripts in the `scripts` folder that will be of use:

Command | Description
--- | ---
`$ source scripts/bosh-env.sh` | Set the BOSH environment variables for the current Director into the environment.
`$ ./scripts/get-ssh-keys.sh`  | Extract SSH keys. These will create two PEM files: `vcap-private-key.pem` and `jumpbox-private-key.pem`. Usage is the usual SSH mechanisms like `ssh -i jumpbox-private-key.pem jumpbox@your-director`. If you change the director, any trusts need to be resolved as usual.
`$ ./scripts/credhub-env.sh`<br/>`$ credhub login`   | Set the CredHub environment variables.
