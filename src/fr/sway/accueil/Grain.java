package fr.sway.accueil;

import android.content.res.Resources;
import android.graphics.Bitmap;
import android.graphics.Shader;
import android.graphics.drawable.BitmapDrawable;
import android.graphics.drawable.Drawable;

import java.util.Random;

final class Grain {
    private static final int TILE = 160;
    private static final int MAX_ALPHA = 30;
    private static final long SEED = 42;
    private static Bitmap tile;

    private Grain() {
    }

    static Drawable drawable(Resources resources) {
        BitmapDrawable grain = new BitmapDrawable(resources, tile());
        grain.setTileModeXY(Shader.TileMode.REPEAT, Shader.TileMode.REPEAT);
        return grain;
    }

    private static Bitmap tile() {
        if (tile == null) {
            Random random = new Random(SEED);
            int[] pixels = new int[TILE * TILE];
            for (int i = 0; i < pixels.length; i++) {
                int alpha = random.nextInt(MAX_ALPHA);
                int shade = random.nextBoolean() ? 0xffffff : 0x000000;
                pixels[i] = alpha << 24 | shade;
            }
            tile = Bitmap.createBitmap(pixels, TILE, TILE, Bitmap.Config.ARGB_8888);
        }
        return tile;
    }
}
