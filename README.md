# Portainer reset password helper

This helper container is designed to update the password of the original administrator account of a Portainer instance.

**Note**: it is designed to reset the password associated to the original administrator account (UserID == 1). If the account is removed, the helper will try to create a user named admin with UserID = 1. If the name admin is already taken, the helper will try to use "admin-" + a random string.

# Usage

## Portainer running as a container

```
# stop the existing Portainer container
docker container stop portainer

# run the helper using the same bind-mount/volume for the data volume
docker run --rm -v portainer_data:/data portainer/helper-reset-password
2020/06/04 00:13:58 Password succesfully updated for user: admin
2020/06/04 00:13:58 Use the following password to login: &_4#\3^5V8vLTd)E"NWiJBs26G*9HPl1

# restart portainer and use the password above to login
docker container start portainer

```

## Portainer running as a stack/service

```
# scale down to zero the existing Portainer service
docker service scale portainer_portainer=0

# run the helper using the same bind-mount/volume for the data volume
docker run --rm -v portainer_portainer_data:/data portainer/helper-reset-password
2020/06/04 00:13:58 Password succesfully updated for user: admin
2020/06/04 00:13:58 Use the following password to login: &_4#\3^5V8vLTd)E"NWiJBs26G*9HPl1

# scale back to one the existing Portainer service and use the password above to login
docker service scale portainer_portainer=1
```

## Command line arguments

The helper accepts the following optional command line arguments:

- `--password` string  
  The new admin password to set. If omitted, the helper will generate a secure random password and print it to stdout.

- `--password-hash` string  
  A pre-computed password hash to set for the admin user (the helper will not re-hash this value). This is useful when you already have a bcrypt/hash value to apply.

- `--data-path` string  
  The path to the Portainer data store inside the container. Defaults to `/data`. Use this to point to a different mount or directory if needed.

Notes:
- `--password` and `--password-hash` are mutually exclusive; do not provide both at the same time.

### Examples:

#### run and let the helper generate a password
```docker run --rm -v portainer_data:/data portainer/helper-reset-password```

#### set a specific password
```docker run --rm -v portainer_data:/data portainer/helper-reset-password --password "MyNewP@ssw0rd"```

#### set a precomputed password hash
```docker run --rm -v portainer_data:/data portainer/helper-reset-password --password-hash "$2y$..."```

## Security

For information about reporting security vulnerabilities, please see our [Security Policy](SECURITY.md).

## Licensing

Portainer reset password helper is licensed under the zlib license. See [LICENSE](./LICENSE) for reference.
