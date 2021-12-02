package com.dragon.httpserver.netty.handler;


import com.dragon.httpserver.dto.Response;
import com.dragon.httpserver.netty.http.DragonHttpRequest;

public interface IFunctionHandler<T> {
    T execute(DragonHttpRequest request);
}
