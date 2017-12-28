# kubernetes-configmap-exporter
Tiny util to export config maps data to installed files.

In some project, we needed to export the Spring-Boot config maps from kubernetes to external spring cloud config server.
Faster option was to write this util. Maybe someone will come in handy.

Linux binary can taken from releases page.

## Usage

To each config map in cubernates set label with file export indifitier:
```yaml
...
  labels:
    dist-applicationm.yml: some-file.json
...
```
- `dist-` - prefix to recognize exported label
- `applicationm.yml` - file from config map to export
- `some-file.json` - file in host to save data from file

```bash
# Export all configs data to files in folder
kubernetes-configmap-exporter -ns myproject -lb dist- -dir configs/
```

To login in cluster uses `.kube/config` dir files. Can override with parameter `kubeconfig`
