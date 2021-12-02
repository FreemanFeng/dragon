package com.dragon.httpserver.helper;

import com.dragon.httpserver.enums.Decoder;
import com.google.gson.Gson;
import org.apache.commons.codec.binary.Base64;

import java.io.UnsupportedEncodingException;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.HashMap;

public class HttpHelper {
    private static final String CONFIG_KEY_URI = "uri";
    private static final String CONFIG_KEY_URL = "url";
    private static final String CONFIG_KEY_METHOD = "method";
    private static final String NAME_ARGS = "args";
    private static final String DEFAULT_KEY_ARGS = "a";
    private static final String NAME_HEADERS = "headers";
    private static final String DEFAULT_KEY_HEADERS = "h";
    private static final String NAME_BODY = "body";
    private static final String DEFAULT_KEY_BODY = "msg";
    private static final String RESP_KEY_CODE = "code";
    private static final String RESP_KEY_CODE2 = "_code";
    private static final String RESP_KEY_HEAD = "head";
    private static final String RESP_KEY_RESP = "resp";

    public static HashMap<String, String> getArguments(HashMap<String, Object> messages, HashMap<String, Object> config) {
        return getParams(DEFAULT_KEY_ARGS, NAME_ARGS, messages, config);
    }

    public static HashMap<String, String> getHeaders(HashMap<String, Object> messages, HashMap<String, Object> config) {
        return getParams(DEFAULT_KEY_HEADERS, NAME_HEADERS, messages, config);
    }

    public static HashMap<String, String> getParams(String key, String name, HashMap<String, Object> messages, HashMap<String, Object> config) {
        String np = key;
        if (config.containsKey(name)) {
            np = Decoder.INSTANCE.decodeStr(config.get(name));
        }
        HashMap<String, Object> mp = Decoder.INSTANCE.decodeMap(messages.get(np));
        HashMap<String, String> params = new HashMap<>();
        for (String k : mp.keySet()) {
            params.put(k, Decoder.INSTANCE.decodeStr(mp.get(k)));
        }
        return params;
    }

    public static String getBody(HashMap<String, Object> messages, HashMap<String, Object> config) {
        String np = getName(DEFAULT_KEY_BODY, NAME_BODY, config);
        return Decoder.INSTANCE.decodeStr(messages.get(np));
    }

    public static HashMap<String, Object> getBodyMap(HashMap<String, Object> messages, HashMap<String, Object> config) {
        String np = getName(DEFAULT_KEY_BODY, NAME_BODY, config);
        return Decoder.INSTANCE.decodeMap(messages.get(np));
    }

    public static ArrayList<Object> getBodyList(HashMap<String, Object> messages, HashMap<String, Object> config) {
        String np = getName(DEFAULT_KEY_BODY, NAME_BODY, config);
        return Decoder.INSTANCE.decodeList(messages.get(np));
    }

    public static String getName(String key, String name, HashMap<String, Object> config) {
        String np = key;
        if (config.containsKey(name)) {
            np = Decoder.INSTANCE.decodeStr(config.get(name));
        }
        return np;
    }

    public static void setBadRequest(HashMap<String, Object> response) {
        HashMap<String, String> head = new HashMap<>();
        HashMap<String, Object> resp = new HashMap<>();
        response.put(RESP_KEY_CODE, 400);
        response.put(RESP_KEY_CODE2, 400);
        response.put(RESP_KEY_HEAD, head);
        response.put(RESP_KEY_RESP, resp);
    }

    public static String getURI(HashMap<String, Object> config) {
        return Decoder.INSTANCE.decodeStr(config.get(CONFIG_KEY_URI));
    }

    public static String getMethod(HashMap<String, Object> config) {
        return Decoder.INSTANCE.decodeStr(config.get(CONFIG_KEY_METHOD));
    }

    public static String getURL(HashMap<String, Object> messages, HashMap<String, Object> config) {
        if(config.containsKey(CONFIG_KEY_URL)){
            return Decoder.INSTANCE.decodeStr(config.get(CONFIG_KEY_URL));
        }
        String uri = getURI(config);
        HashMap<String, String> args = HttpHelper.getArguments(messages, config);
        if (args.size() > 0) {
            StringBuilder path = new StringBuilder("?");
            int i = 0;
            for (String key : args.keySet()) {
                try {
                    String v = Decoder.INSTANCE.decodeStr(args.get(key));
                    String data = java.net.URLEncoder.encode(v, "UTF-8");
                    path.append(key).append("=").append(data);
                } catch (UnsupportedEncodingException e) {
                    e.printStackTrace();
                    return uri;
                }
                if (i < args.size() - 1) {
                    path.append("&");
                }
                i++;
            }
            uri += path.toString();
        }
        return uri;
    }
}
