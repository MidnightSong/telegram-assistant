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
        }
        return Response.json({Code: 5000});
    }
};
export {
    src_default as default
};
//# sourceMappingURL=index.js.map
