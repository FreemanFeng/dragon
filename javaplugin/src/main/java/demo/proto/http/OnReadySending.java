package demo.proto.http;

import com.dragon.httpserver.enums.Decoder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;

public class OnReadySending {
    private static final Logger LOGGER = LoggerFactory.getLogger(OnReadySending.class);
    public static void call(HashMap<String, Object> m, HashMap<String, Object> c, HashMap<String, Object> r) {
        Object a = m.get("a");
        HashMap<String, Object> ma = Decoder.INSTANCE.decodeMap(a);
        ma.put("hello2", "world");
        m.put("a", ma);
        LOGGER.info("OnReadySending");
    }
}
