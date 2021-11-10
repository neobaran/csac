# CSAC (Cloud Service Auto Cert)

## Features

- [x] [Tencent Cloud](https://cloud.tencent.com/)
  - [x] Auto DNS Certification
  - [x] Auto Upload Cert File
  - [x] Auto Update CDN HTTPS Cert
  - [ ] Auto Update CLB HTTPS Cert
  - [ ] Auto Update API Gateway HTTPS Cert
  - [ ] Other...
- [ ] Alibaba Cloud
- [ ] Other...

## Config

```yaml
# config.yaml
Email: your@email.com
Tencent:
  SecretId: your-tencent-cloud-secret-id
  SecretKey: your-tencent-cloud-secret-key
Domains:
  - *.example.com
```

Or Set Docker environment variables

```conf
# .env
CSAC_EMAIL=your@email.com
CSAC_TENCENT_SECRET_ID=your-tencent-cloud-secret-id
CSAC_TENCENT_SECRET_KEY=your-tencent-cloud-secret-key
CSAC_DOMAIN=*.example.com
```

## Usage

```sh
docker run -v /path/to/config.yaml:/config.yaml neobaran/csac -c /config.yaml
```

```sh
docker run --env-file .env neobaran/csac
```
