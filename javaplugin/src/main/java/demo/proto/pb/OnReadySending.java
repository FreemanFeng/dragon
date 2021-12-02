package demo.proto.pb;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;

public class OnReadySending {
    private static final Logger LOGGER = LoggerFactory.getLogger(OnReadySending.class);
    public static void call(HashMap<String, Object> m, HashMap<String, Object> c, HashMap<String, Object> r) {
        LOGGER.info("OnReadySending");
    }
}
