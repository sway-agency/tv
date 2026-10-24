package fr.sway.accueil;

import android.content.res.AssetManager;
import android.graphics.Typeface;
import android.util.SparseArray;

import java.io.IOException;

final class Fonts {
    private static final String BRAND_FONT = "fonts/satoshi.ttf";
    private static final String FALLBACK_FAMILY = "sans-serif";
    private static final SparseArray<Typeface> BY_WEIGHT = new SparseArray<>();
    private static Boolean brandFontPresent;

    private Fonts() {
    }

    static Typeface satoshi(AssetManager assets, int weight) {
        Typeface cached = BY_WEIGHT.get(weight);
        if (cached == null) {
            cached = load(assets, weight);
            BY_WEIGHT.put(weight, cached);
        }
        return cached;
    }

    private static Typeface load(AssetManager assets, int weight) {
        if (!hasBrandFont(assets)) {
            return Typeface.create(Typeface.create(FALLBACK_FAMILY, Typeface.NORMAL), weight, false);
        }
        return new Typeface.Builder(assets, BRAND_FONT)
                .setFontVariationSettings("'wght' " + weight)
                .build();
    }

    private static boolean hasBrandFont(AssetManager assets) {
        if (brandFontPresent == null) {
            try {
                assets.open(BRAND_FONT).close();
                brandFontPresent = true;
            } catch (IOException e) {
                brandFontPresent = false;
            }
        }
        return brandFontPresent;
    }
}
