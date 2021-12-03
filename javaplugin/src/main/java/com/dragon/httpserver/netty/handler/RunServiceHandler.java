//   Copyright 2019 Freeman Feng<freeman@nuxim.cn>
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
package com.dragon.httpserver.netty.handler;

import com.dragon.httpserver.enums.*;
import com.dragon.httpserver.netty.annotation.DragonHttpHandler;
import com.dragon.httpserver.netty.http.DragonHttpRequest;
import com.google.gson.Gson;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.google.gson.reflect.TypeToken;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.apache.commons.codec.binary.Base64;
import org.springframework.util.StringUtils;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.nio.charset.StandardCharsets;
import java.util.*;


@DragonHttpHandler(path = "/run/", equal = false, method = "POST")
public class RunServiceHandler implements IFunctionHandler<String> {
    private static final Logger LOGGER = LoggerFactory.getLogger(RunServiceHandler.class);

    @Override
    public String execute(DragonHttpRequest request) {

        String service = request.getStringPathValue(2);
        LOGGER.info("run service {}", service);
        Class<?> clazz = Callback.INSTANCE.getCallback(service);
        String content = request.contentText();
        Gson gson = new Gson();
        String[] ss = gson.fromJson(content, String[].class);
        LOGGER.info("data {}", Arrays.toString(ss));
        int n = ss.length / 2;
        Class[] clazzs = new Class[n];
//        ArrayList<Object> ms = new ArrayList<>();
        Object[] ms = new Object[n];
        for (int i = 0; i < n; i++) {
            String k = ss[2 * i];
            String v = ss[2 * i + 1];
            LOGGER.info("matched {}", k);
            if (k.equals("MapBytes")) {
                ms[i] = Decoder.INSTANCE.decodeMapBytes(v);
                continue;
            }
            if (k.equals("MAP")) {
                String s = Decoder.INSTANCE.decodeBase64(v);
                ms[i] = Decoder.INSTANCE.decodeMap(s);
                continue;
            }
            if (k.equals("ListBytes")) {
                ms[i] = Decoder.INSTANCE.decodeListBytes(v);
                continue;
            }
            if (k.equals("LIST")) {
                String s = Decoder.INSTANCE.decodeBase64(v);
                ms[i] = Decoder.INSTANCE.decodeList(s);
                continue;
            }
            if (k.equals("STR")) {
                ms[i] = v;
                continue;
            }
            if (k.equals("BYTES")) {
                ms[i] = v.getBytes(StandardCharsets.UTF_8);
                continue;
            }
            if (k.equals("INT") || k.equals("UINT")) {
                Integer ii = new Integer(v);
                ms[i] = ii;
                continue;
            }
            if (k.equals("FLOAT")) {
                Float ii = new Float(v);
                ms[i] = ii;
                continue;
            }
            if (k.equals("BOOL")) {
                if (v.equals("True") || v.equals("TRUE") || v.equals("true")) {
                    ms[i] = true;
                } else {
                    ms[i] = false;
                }
            }
        }
        try {
            Method[] methods = clazz.getDeclaredMethods();
            LOGGER.info("call methods {}", Arrays.toString(methods));
//            Method method = clazz.getDeclaredMethod("run", clazzs);
            Method method = methods[0];
            String[] sk = StringUtils.split(service, ".");
            if (sk != null && sk.length > 1) {
                method.invoke(null, ms);
                LOGGER.info("after callback {}", ms);
                String r = Encoder.INSTANCE.encodeListBytes(ms);
                LOGGER.info("return {}", r);
                return r;
            }
            Object b = method.invoke(null, ms);
            LOGGER.info("return {}", b.toString());
            return b.toString();
        } catch (IllegalAccessException | InvocationTargetException e) {
            e.printStackTrace();
        }
        return "NOK";
    }
}
