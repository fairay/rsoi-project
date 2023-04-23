# Что должно лежать в секретах

**DOCKERHUB_LOGIN** - имя (не почта) профиля с https://hub.docker.com/

**DOCKERHUB_PASSWORD** - пароль для входа в https://hub.docker.com/

**K8S_HOST** - External endpoing из
``` bash
yc managed-kubernetes cluster list
```

**K8S_TOKEN** - "token" из
``` bash
yc managed-kubernetes create-token 
```
Имеет свойство протухать


# Прочее

**Okta RSA keys**
``` bash
curl https://dev-35033098.okta.com/oauth2/default/v1/keys
```
