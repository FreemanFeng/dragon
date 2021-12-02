package demo.proto.http;

import com.dragon.httpserver.helper.HttpHelper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;

public class OnSending {
    private static final Logger LOGGER = LoggerFactory.getLogger(OnSending.class);
    public static void call(HashMap<String, Object> m, HashMap<String, Object> c, HashMap<String, Object> r) {
        String method = HttpHelper.getMethod(c);
        String url = HttpHelper.getURL(m, c);
        HashMap<String, String> headers = HttpHelper.getHeaders(m, c);
        String body = HttpHelper.getBody(m, c);
        LOGGER.info("OnSending");
    }
}
