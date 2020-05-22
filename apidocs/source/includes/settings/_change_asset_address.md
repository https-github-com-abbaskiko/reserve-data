## Pending change asset address

```shell
curl -X POST "https://gateway.local/v3/setting-change-main" \
-H 'Content-Type: application/json' \
-d '{
    "change_list": [
        {
            "type": "change_asset_addr",
            "data": {
                "id": 1,
                "address": "0xC7DC5C95728d9ca387239Af0A49b7BCe8927d309"
            }
        }
    ]
}'
```

> sample json

```json
{
    "id": 1,
    "success": true
}
```

### HTTP Request

`POST https://gateway.local/v3/setting-change-main`
<aside class="notice">Write key is required</aside>

###Data fields:

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
id | int | true | nil | id of asset
address | string (ethereum address) | true | nil | new address