// src/utils.ts
function base64Decode(str) {
    return Uint8Array.from(atob(str), (c) => c.charCodeAt(0));
}

function bytesToBase64(bytes) {
    const binString = Array.from(
        bytes,
        (byte) => String.fromCodePoint(byte)
    ).join("");
    return btoa(binString);
}

function hexToArrayBuffer(hexString) {
    var typedArray = new Uint8Array(hexString.match(/[\da-f]{2}/gi).map(function (h) {
        return parseInt(h, 16);
    }));
    return typedArray.buffer;
}

async function loadKey(key2) {
    let keyBuffer = hexToArrayBuffer(key2);
    let cryptoKey = await crypto.subtle.importKey(
        "raw",
        keyBuffer,
        {name: "AES-GCM"},
        false,
        ["encrypt", "decrypt"]
    );
    return cryptoKey;
}

async function encryptAES_GCM(key2, plaintext, iv2) {
    const encodedText = new TextEncoder().encode(plaintext);
    let k = await loadKey(key2);
    let nonce = hexToArrayBuffer(iv2);
    const encryptedData = await crypto.subtle.encrypt(
        {
            name: "AES-GCM",
            iv: nonce
        },
        k,
        encodedText
    );
    return {
        ciphertext: new Uint8Array(encryptedData),
        iv: iv2
    };
}

async function decryptAES_GCM(key2, ciphertext, iv2) {
    let k = await loadKey(key2);
    let nonce = hexToArrayBuffer(iv2);
    const decryptedData = await crypto.subtle.decrypt(
        {
            name: "AES-GCM",
            iv: nonce
        },
        k,
        ciphertext
    );
    return new TextDecoder().decode(decryptedData);
}

// src/index.ts
var key = "57227176a09c27191875e85ce2ccea571e415fd98038ccb21e892c4d7182bc3e";
var iv = "ff5097cd1d355f6d6f8d9225";

async function validActiveCode(request, env, ctx) {
    let resp = {code: 4003, msg: "无效激活码"};
    var body = await request.text();
    try {
        var decodedString = await decryptAES_GCM(key, base64Decode(body), iv);
        let param = JSON.parse(decodedString);
        console.log(param);
        if (param.uuid === '') {
            return resp;
        }
        let val = await env.used_key.get(param.uuid);
        if (val == null) {
            return resp;
        }
        const data = JSON.parse(val);
        //第一次激活
        if (data.device_id === undefined || data.device_id === null || data.device_id === ''){
            data.device_id = param.device_id
            await env.used_key.put(param.uuid, JSON.stringify(data));
            resp.code = 2000;
            resp.msg = "激活成功";
            resp.data = {uuid: param.uuid, exp: data.expire_at, duration, timestamp};
            return resp;
        }

        if (param.device_id !== data.device_id) {
            resp.msg = "激活码已被使用";
            return resp
        }

        const timestamp = Math.floor(Date.now() / 1000);
        const duration = data.expire_at - timestamp;
        if (duration < 0) {
            resp.msg = "激活码已过期";
            return resp;
        }
        //await env.used_key.put(param.uuid, JSON.stringify(code));
        resp.code = 2000;
        resp.msg = "激活成功";
        resp.data = {uuid: param.uuid, exp: data.expire_at, duration, timestamp};
        return resp;
    } catch (error) {
        resp.code = 4000;
        resp.msg = "无效输入";
        console.error(error);
        return resp;
    }
}

var src_default = {
    async fetch(request, env, ctx) {
        const u = URL.parse(request.url);
        if (!u) {
            return Response.json({Code: 4000});
        }
        switch (u.pathname) {
            case "/acv":
                let resp = await validActiveCode(request, env, ctx);
                const data = await encryptAES_GCM(key, JSON.stringify(resp), iv);
                return new Response(bytesToBase64(data.ciphertext));
            default :
                return new Response(htmlContent, {
                    headers: {
                        'Content-Type': 'text/html'
                    }
                });
        }
    }
};
export {
    src_default as default
};
const htmlContent = `
  <!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dancing Robot</title>
    <style>
        body {
            margin: 0;
            padding: 0;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            background-color: #f5f5f7;
            color: #333;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            height: 100vh;
            text-align: center;
        }
        h1 {
            font-size: 2.5rem;
            font-weight: 600;
            margin-bottom: 20px;
            color: #1d1d1f;
        }
        p {
            font-size: 1.2rem;
            margin-bottom: 40px;
            color: #6e6e73;
        }
        img {
            max-width: 100%;
            height: auto;
            border-radius: 12px;
            box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
            transition: transform 0.3s ease;
        }
        img:hover {
            transform: scale(1.05);
        }
        .button {
            margin-top: 20px;
            padding: 10px 20px;
            font-size: 1rem;
            color: #fff;
            background-color: #007aff;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            text-decoration: none;
            display: inline-flex;
            align-items: center;
            gap: 8px;
            transition: background-color 0.3s ease;
        }
        .button:hover {
            background-color: #005bb5;
        }
        .button img {
            width: 20px;
            height: 20px;
        }
    </style>
</head>
<body>
    <h1>个人号机器人租用</h1>
    <p>消息群发、转发等</p>
    <p>解决官方机器人无法互相识别消息</p>
    <img src="https://media.giphy.com/media/26xBwdIuRJiAIqHwA/giphy.gif" alt="Dancing Robot">
    <a href="https://t.me/ZuLinpWuYu" class="button">
        <img src="https://upload.wikimedia.org/wikipedia/commons/8/82/Telegram_logo.svg" alt="Telegram Icon">
        在Telegram上联系我
    </a>
</body>
</html>
`;
