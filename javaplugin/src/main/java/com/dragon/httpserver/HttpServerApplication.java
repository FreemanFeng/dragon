package com.dragon.httpserver;

import com.dragon.httpserver.netty.annotation.DragonHttpHandler;
import org.springframework.boot.WebApplicationType;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.builder.SpringApplicationBuilder;
import org.springframework.context.annotation.ComponentScan;

@SpringBootApplication()
@ComponentScan(includeFilters = @ComponentScan.Filter(DragonHttpHandler.class))

public class HttpServerApplication {

    public static void main(String[] args) {
        new SpringApplicationBuilder(HttpServerApplication.class).web(WebApplicationType.NONE).run(args);
    }

}
