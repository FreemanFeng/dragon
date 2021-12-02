package demo.proto.pb;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;

public class OnError {
    private static final Logger LOGGER = LoggerFactory.getLogger(OnError.class);
    public static void call(HashMap<String, Object> m, HashMap<String, Object> c, HashMap<String, Object> r) {
        LOGGER.info("OnError");
    }
}
