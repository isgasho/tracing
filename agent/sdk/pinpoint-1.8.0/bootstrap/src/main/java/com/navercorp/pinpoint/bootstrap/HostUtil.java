package com.navercorp.pinpoint.bootstrap;

import java.net.InetAddress;
import java.net.UnknownHostException;

public class HostUtil {
    private static String getHostNameForLiunx() {
        try {
            return (InetAddress.getLocalHost()).getHostName();
        } catch (UnknownHostException uhe) {
            String host = uhe.getMessage(); // host = "hostname: hostname"
            if (host != null) {
                int colon = host.indexOf(':');
                if (colon > 0) {
                    return host.substring(0, colon);
                }
            }
            return "UnknownHost";
        }
    }

    public static String getHostName() {
        if (System.getenv("COMPUTERNAME") != null) {
            return System.getenv("COMPUTERNAME");
        } else {
            return getHostNameForLiunx();
        }
    }

    public static String getApplicationName(String hostName) {
        String[] names = hostName.split("-");
        if (names.length == 1) {
            return hostName;
        } else if (names.length == 3) {
            return names[1];
        } else if (names.length == 4) {
            return names[1] + names[2];
        } else {
            return null;
        }
    }

    public static String getAgentId(String hostName) {
        String[] names = hostName.split("-");
        if (names.length == 1) {
            return hostName;
        } else if (names.length == 3) {
            String id = "";
            if (names[2].toLowerCase().equals("vip")) {
                id = "v";
            } else if (names[2].toLowerCase().equals("yf")) {
                id = "y";
            } else {
                id = names[2];
            }
            return names[1] + id;
        } else if (names.length == 4) {
            String id = "";
            if (names[3].toLowerCase().equals("vip")) {
                id = "v";
            } else if (names[3].toLowerCase().equals("yf")) {
                id = "y";
            } else {
                id = names[3];
            }
            return names[1] + names[2] + id;
        } else {
            return null;
        }
    }

    public static void main(String[] args) {
        System.out.println(getApplicationName("jt-web-PressureTest-34"));
        System.out.println(getAgentId("jt-web-PressureTest-34"));
    }
}
