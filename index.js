// src/index.ts
async function aesDecrypt(cipherText, key2, iv2) {
    let cipherBuffer = Uint8Array.from(atob(cipherText), (c) => c.charCodeAt(0));
    let keyBuffer = Uint8Array.from(atob(key2), (c) => c.charCodeAt(0));
    let ivBuffer = Uint8Array.from(atob(iv2), (c) => c.charCodeAt(0));
    let cryptoKey = await crypto.subtle.importKey(
      "raw",
      keyBuffer,
      { name: "AES-CBC" },
      false,
      ["decrypt"]
    );
    let decryptedBuffer = await crypto.subtle.decrypt(
      { name: "AES-CBC", iv: ivBuffer },
      cryptoKey,
      cipherBuffer
    );
    let decoder = new TextDecoder();
    let decryptedText = decoder.decode(decryptedBuffer);
    return decryptedText;
  }
  var key = "lNJaJ7DNsIO+V6djr4zYM06O2ERnEHEO9BBhom4bhdI=";
  var iv = "ik3XwT3CRFUVBMAcshtyKw==";
  var src_default = {
    async fetch(request, env, ctx) {
      let resp = { Code: 2e3, Msg: "" };
      var body = await request.text();
      try {
        var decodedString = await aesDecrypt(body, key, iv);
        let param = JSON.parse(decodedString);
        let val = await env.used_key.get(param.Code);
        if (val != null) {
          resp.Code = 4003;
          resp.Msg = "\u6FC0\u6D3B\u7801\u5DF2\u88AB\u4F7F\u7528";
          console.log(param.Code, resp.Msg, val);
          return Response.json(resp);
        }
        await env.used_key.put(param.Code, param.DeviceID);
        return Response.json(resp);
      } catch (error) {
        resp.Code = 4e3;
        resp.Msg = "\u65E0\u6548\u8F93\u5165";
        return Response.json(resp);
      }
    }
  };
  export {
    src_default as default
  };
  //# sourceMappingURL=index.js.map
  