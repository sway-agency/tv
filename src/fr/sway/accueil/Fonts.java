package fr.sway.accueil;

import android.content.res.AssetManager;
import android.graphics.Typeface;

import java.io.IOException;
import java.io.InputStream;

final class Fonts {
    private static final String BRAND_FONT = "fonts/satoshi.ttf";
    private static final String FALLBACK_FAMILY = "sans-serif";

    private Fonts() {
    }

    static Typeface satoshi(AssetManager assets, int weight) {
        if (!exists(assets, BRAND_FONT)) {
            return Typeface.create(Typeface.create(FALLBACK_FAMILY, Typeface.NORMAL), weight, false);
        }
        return new Typeface.Builder(assets, BRAND_FONT)
                .setFontVariationSettings("'wght' " + weight)
                .build();
    }

    private static boolean exists(AssetManager assets, String path) {
        try (InputStream ignored = assets.open(path)) {
            return true;
        } catch (IOException e) {
            return false;
        }
    }
}
