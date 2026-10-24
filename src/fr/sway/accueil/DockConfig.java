package fr.sway.accueil;

import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.text.TextUtils;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.regex.Pattern;

final class DockConfig {
    static final String EXTRA_DOCK = "dock";

    private static final String PREFERENCES = "dock";
    private static final String KEY_PACKAGES = "packages";
    private static final Pattern PACKAGE_NAME = Pattern.compile("[a-zA-Z][\\w]*(\\.[a-zA-Z][\\w]*)+");
    private static final List<String> DEFAULT_PACKAGES = Arrays.asList(
            "net.oqee.androidtv.store",
            "com.netflix.ninja",
            "com.disney.disneyplus",
            "com.apple.atve.androidtv.appletv",
            "com.canal.android.canal");

    private final SharedPreferences preferences;

    DockConfig(Context context) {
        preferences = context.getSharedPreferences(PREFERENCES, Context.MODE_PRIVATE);
    }

    void applyIntent(Intent intent) {
        String[] packages = intent == null ? null : intent.getStringArrayExtra(EXTRA_DOCK);
        if (packages == null) {
            return;
        }
        List<String> valid = new ArrayList<>();
        for (String packageName : packages) {
            if (PACKAGE_NAME.matcher(packageName.trim()).matches()) {
                valid.add(packageName.trim());
            }
        }
        preferences.edit().putString(KEY_PACKAGES, TextUtils.join(",", valid)).apply();
    }

    List<String> packages() {
        String stored = preferences.getString(KEY_PACKAGES, null);
        if (stored == null) {
            return DEFAULT_PACKAGES;
        }
        List<String> packages = new ArrayList<>();
        for (String packageName : stored.split(",")) {
            if (!packageName.trim().isEmpty()) {
                packages.add(packageName.trim());
            }
        }
        return packages;
    }
}
