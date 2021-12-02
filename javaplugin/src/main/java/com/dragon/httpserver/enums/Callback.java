package com.dragon.httpserver.enums;

import java.util.HashMap;

public enum Callback {
    INSTANCE;
    private final HashMap<String, Class<?>> callbacks;

    Callback() {
        callbacks = new HashMap<>();
    }

    public HashMap<String, Class<?>> getCallbacks() {
        return callbacks;
    }

    public Class<?> getCallback(String key) {
        return callbacks.get(key);
    }

    public void setCallback(String key, Class<?> handler) {
        callbacks.put(key, handler);
    }
}
