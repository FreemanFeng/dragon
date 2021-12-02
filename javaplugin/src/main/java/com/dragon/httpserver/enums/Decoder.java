package com.dragon.httpserver.enums;

import com.google.gson.Gson;
import com.google.gson.reflect.TypeToken;
import org.apache.commons.codec.binary.Base64;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.HashMap;

public enum Decoder {
    INSTANCE;
    private static final Logger LOGGER = LoggerFactory.getLogger(Decoder.class);
    private final Gson gson;

    Decoder() {
        gson = new Gson();
    }

    public Object decode(String json) {
        return gson.fromJson(json, Object.class);
    }
    public String decodeBase64(String data) {
        return new String(Base64.decodeBase64(data.getBytes(StandardCharsets.UTF_8)), StandardCharsets.UTF_8);
    }
    public String decodeStr(Object src) {
        return Encoder.INSTANCE.encode(src);
    }
    public Integer decodeInt(Object src) {
        String json = Encoder.INSTANCE.encode(src);
        return gson.fromJson(json, Integer.class);
    }
    public Boolean decodeBool(Object src) {
        String json = Encoder.INSTANCE.encode(src);
        return gson.fromJson(json, Boolean.class);
    }
    public Float decodeFloat(Object src) {
        String json = Encoder.INSTANCE.encode(src);
        return gson.fromJson(json, Float.class);
    }
    public HashMap<String, Object> decodeMap(String json) {
        return gson.fromJson(json, new TypeToken<HashMap<String, Object>>(){}.getType());
    }
    public HashMap<String, Object> decodeMap(Object src) {
        String json = Encoder.INSTANCE.encode(src);
        return gson.fromJson(json, new TypeToken<HashMap<String, Object>>(){}.getType());
    }
    public ArrayList<Object> decodeList(String json) {
        return gson.fromJson(json, new TypeToken<ArrayList<Object>>(){}.getType());
    }
    public ArrayList<Object> decodeList(Object src) {
        String json = Encoder.INSTANCE.encode(src);
        return gson.fromJson(json, new TypeToken<ArrayList<Object>>(){}.getType());
    }
    public HashMap<String, Object> decodeMapBytes(String json) {
        String s = decodeBase64(json);
        HashMap<String, String> x = gson.fromJson(s, new TypeToken<HashMap<String, String>>(){}.getType());
        HashMap<String, Object> m = new HashMap<>();
        for(String key:x.keySet()){
            String value = x.get(key);
            String sub = decodeBase64(value);
            Object b = gson.fromJson(sub, Object.class);
            m.put(key, b);
        }
        return m;
    }
    public ArrayList<Object> decodeListBytes(String json) {
        String s = decodeBase64(json);
        ArrayList<String> x = gson.fromJson(s, new TypeToken<ArrayList<String>>(){}.getType());
        ArrayList<Object> a = new ArrayList<>();
        for(String value:x){
            String sub = decodeBase64(value);
            Object ma = gson.fromJson(sub, Object.class);
            a.add(ma);
        }
        return a;
    }
}
