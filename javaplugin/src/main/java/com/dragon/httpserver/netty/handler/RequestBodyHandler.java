package com.dragon.httpserver.netty.handler;
import com.dragon.httpserver.dto.Response;
import com.dragon.httpserver.netty.annotation.DragonHttpHandler;
import com.dragon.httpserver.netty.http.DragonHttpRequest;


@DragonHttpHandler(path = "/request/body",method = "POST")
public class RequestBodyHandler implements IFunctionHandler<String> {
    @Override
    public String execute(DragonHttpRequest request) {
        /**
         * 可以在此拿到json转成业务需要的对象
         */
        String json = request.contentText();
        return json;
    }
}
