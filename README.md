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
