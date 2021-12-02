package com.dragon.httpserver.enums;

import com.google.gson.Gson;
import com.google.gson.reflect.TypeToken;
import org.apache.commons.codec.binary.Base64;

import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.HashMap;

public enum Encoder {
    INSTANCE;
    private final Gson gson;

    Encoder() {
        gson = new Gson();
    }

    public String encodeBase64(byte[] binaryData) {
        return new String(Base64.encodeBase64(binaryData), StandardCharsets.UTF_8);
    }

    public String encode(Object src) {
        return gson.toJson(src);
    }

    public String encodeMapBytes(HashMap<String, Object> src) {
        HashMap<String, Object> m = new HashMap<>();
        for(String key:src.keySet()){
            String value = encode(src.get(key));
            String sub = encodeBase64(value.getBytes(StandardCharsets.UTF_8));
            m.put(key, sub);
        }
        return encode(m);
//        return encodeBase64(r.getBytes(StandardCharsets.UTF_8));
    }

    public String encodeListBytes(Object[] src) {
        ArrayList<Object> m = new ArrayList<>();
        for(Object data:src){
            String value = encode(data);
            String sub = encodeBase64(value.getBytes(StandardCharsets.UTF_8));
            m.add(sub);
        }
        return encode(m);
//        return encodeBase64(r.getBytes(StandardCharsets.UTF_8));
    }
}
