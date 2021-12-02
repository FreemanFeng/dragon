package com.dragon.httpserver.netty.iohandler;

import com.dragon.httpserver.netty.http.DragonHttpResponse;
import io.netty.channel.ChannelFutureListener;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;
import io.netty.handler.codec.http.FullHttpRequest;
import io.netty.handler.codec.http.HttpResponseStatus;
import io.netty.util.ReferenceCountUtil;
import org.springframework.stereotype.Component;


@ChannelHandler.Sharable
@Component
/**
 * 在这里可以做拦截器，验证一些请求的合法性
 */
public class InterceptorHandler extends ChannelInboundHandlerAdapter {
    @Override
    public void channelRead(ChannelHandlerContext context, Object msg)   {
        if (isPassed((FullHttpRequest) msg)){
            context.fireChannelRead(msg);
            return;
        }

        ReferenceCountUtil.release(msg);
        context.writeAndFlush(DragonHttpResponse.make(HttpResponseStatus.UNAUTHORIZED)).addListener(ChannelFutureListener.CLOSE);
    }

    /**
     * 修改实现来验证合法性
     * @param request
     * @return
     */
    private boolean isPassed(FullHttpRequest request){
        return true;
    }
}
