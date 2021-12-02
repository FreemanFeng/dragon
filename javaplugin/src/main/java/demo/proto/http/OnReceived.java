package demo.proto.http;

import com.dragon.httpserver.enums.Decoder;
import com.dragon.httpserver.enums.Encoder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;

public class OnReceived  {
    private static final Logger LOGGER = LoggerFactory.getLogger(OnReceived.class);
    public static void call(HashMap<String, Object> m, HashMap<String, Object> c, HashMap<String, Object> r) {
        Object a = r.get("resp");
        HashMap<String, Object> ma = Decoder.INSTANCE.decodeMap(a);
        ma.put("hello2", "world");
        r.put("resp", ma);
        LOGGER.info("OnReceived");
    }
}
