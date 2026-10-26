# Metaverse Postman workspace

Import only `metaverse.postman_collection.json`. It includes the local API and WebSocket URLs, sample credentials, and runtime variables.

```bash
make docker-up
make migrate-up
make api-run
make realtime-run
```

Run user/admin authentication first; its scripts save the tokens as collection variables. Create an avatar, element, map, and then a space to populate the identifiers used by later requests. Open the WebSocket request and send the documented `join` and `move` frames in Postman's message composer.

HTTP authentication uses a `token` cookie, not an Authorization header. Realtime uses the same user token inside the `join` event payload.

The collection uses Postman Collection Format v2.1. [Official schema](https://schema.postman.com/json/collection/v2.1.0/docs/index.html)
