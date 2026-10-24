package fr.sway.accueil;

import android.graphics.Bitmap;

final class Blur {
    private static final int WIDTH = 160;
    private static final int RADIUS = 3;
    private static final int PASSES = 3;

    private Blur() {
    }

    static Bitmap of(Bitmap source) {
        if (source == null) {
            return null;
        }
        int height = Math.max(1, WIDTH * source.getHeight() / source.getWidth());
        Bitmap small = Bitmap.createScaledBitmap(source, WIDTH, height, true).copy(Bitmap.Config.ARGB_8888, true);
        int[] pixels = new int[WIDTH * height];
        small.getPixels(pixels, 0, WIDTH, 0, 0, WIDTH, height);
        for (int pass = 0; pass < PASSES; pass++) {
            pixels = boxPass(pixels, WIDTH, height, 1, 0);
            pixels = boxPass(pixels, WIDTH, height, 0, 1);
        }
        small.setPixels(pixels, 0, WIDTH, 0, 0, WIDTH, height);
        return small;
    }

    private static int[] boxPass(int[] source, int width, int height, int stepX, int stepY) {
        int[] target = new int[source.length];
        for (int y = 0; y < height; y++) {
            for (int x = 0; x < width; x++) {
                target[y * width + x] = average(source, width, height, x, y, stepX, stepY);
            }
        }
        return target;
    }

    private static int average(int[] source, int width, int height, int x, int y, int stepX, int stepY) {
        int red = 0;
        int green = 0;
        int blue = 0;
        for (int offset = -RADIUS; offset <= RADIUS; offset++) {
            int sampleX = Math.min(width - 1, Math.max(0, x + offset * stepX));
            int sampleY = Math.min(height - 1, Math.max(0, y + offset * stepY));
            int color = source[sampleY * width + sampleX];
            red += (color >> 16) & 0xff;
            green += (color >> 8) & 0xff;
            blue += color & 0xff;
        }
        int count = 2 * RADIUS + 1;
        return 0xff000000 | (red / count) << 16 | (green / count) << 8 | (blue / count);
    }
}
