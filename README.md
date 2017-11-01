# Upspin `OpenStack` repository

## Installation

```
upspin-setupstorage-openstack --domain=upspin.example.org --region=GRA1 upspin-container
```

You can now proceed to the regular installation using [setupserver](https://upspin.io/doc/server_setup.md) except that you need to put the exported var of your openrc.sh file into systemd unit file