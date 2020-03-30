# GCP Firewall API

This repository provides an API to create and manage Firewall rules in a GCP host project using API for an application.

Want to go further ?

- [ ] Add Authentication
- [ ] Manage RBAC
- [ ] Add acceptance criterias on rules
- [ ] Force targetTags as we force rule Name

## Disclamer

DEMO only. Do not use it on production. _Done over a long night during covid lockdown...._

## Test it !

Create rules for an applications

```bash
$ curl -X POST 127.0.0.1:8080/project/cka-jnu/service_project/foo-sp/application/kubernetes-the-hard-way --data '[{"CustomName": "test-ssh", "Rule": {"name": "dummy","network": "global/networks/default","allowed": [{"IPProtocol": "TCP", "ports": ["22"]}],"targetTags": ["foo"]}}]'
```

Verify rules for the application created

```bash
$ curl 127.0.0.1:8080/project/cka-jnu/service_project/foo-sp/application/kubernetes-the-hard-way | jq
```

Delete rules for the application created

```bash
$ curl -X DELETE 127.0.0.1:8080/project/cka-jnu/service_project/foo-sp/application/kubernetes-the-hard-way | jq
```

## Rules

Rules are based on Google compute API [rest/v1/firewalls](https://cloud.google.com/compute/docs/reference/rest/v1/firewalls)

The tool erase the rule name (if provided) to set a custom name like `serviceProject-applicationName-customName` to avoid dupplicated name and make easier list, update of deletion.
