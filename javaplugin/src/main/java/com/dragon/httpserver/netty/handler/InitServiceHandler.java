package com.dragon.httpserver.netty.handler;


import com.dragon.httpserver.dto.Response;
import com.dragon.httpserver.enums.Callback;
import com.dragon.httpserver.netty.annotation.DragonHttpHandler;
import com.dragon.httpserver.netty.http.DragonHttpRequest;
import com.google.gson.Gson;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

@DragonHttpHandler(path = "/init")
public class InitServiceHandler implements IFunctionHandler<String> {

    @Override
    public String execute(DragonHttpRequest request) {
        HashMap<String, Class<?>> services = Callback.INSTANCE.getCallbacks();
        ArrayList<String> r = new ArrayList<>(services.keySet());
        Gson gson = new Gson();
        return gson.toJson(r);
    }
}
