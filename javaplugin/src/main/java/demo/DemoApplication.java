package demo;

import com.dragon.httpserver.HttpServerApplication;
import com.dragon.httpserver.enums.Callback;
import demo.proto.http.OnSending;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class DemoApplication extends HttpServerApplication {
    private static final Logger LOGGER = LoggerFactory.getLogger(DemoApplication.class);

    public static void init() {
        Callback.INSTANCE.setCallback("http.OnReadySending", demo.proto.http.OnReadySending.class);
        Callback.INSTANCE.setCallback("http.OnSending", demo.proto.http.OnSending.class);
        Callback.INSTANCE.setCallback("http.OnReceived", demo.proto.http.OnReceived.class);
        Callback.INSTANCE.setCallback("http.OnError", demo.proto.http.OnError.class);
        Callback.INSTANCE.setCallback("pb.OnReadySending", demo.proto.pb.OnReadySending.class);
        Callback.INSTANCE.setCallback("pb.OnReceived", demo.proto.pb.OnReceived.class);
        Callback.INSTANCE.setCallback("pb.OnError", demo.proto.pb.OnError.class);
        Callback.INSTANCE.setCallback("Sign", demo.common.Sign.class);
    }
    public static void main(String[] args) {
        init();
        LOGGER.info(Callback.INSTANCE.getCallbacks().toString());
        HttpServerApplication.main(args);
    }
}
