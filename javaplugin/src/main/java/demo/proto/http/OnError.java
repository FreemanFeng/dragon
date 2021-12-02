package demo.proto.http;

import com.dragon.httpserver.enums.Decoder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;

public class OnError {
    private static final Logger LOGGER = LoggerFactory.getLogger(OnError.class);
    public static void call(HashMap<String, Object> m, HashMap<String, Object> c, HashMap<String, Object> r) {
        if (r.containsKey("code")) {
            Object a = r.get("code");
            Integer code = Decoder.INSTANCE.decodeInt(a);
            if (code == 400) {
                r.put("code", code);
            }
        }
        LOGGER.info("OnError");
    }
}
