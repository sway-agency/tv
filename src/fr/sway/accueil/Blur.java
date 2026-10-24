package fr.sway.accueil;

import android.graphics.Bitmap;
import android.graphics.Canvas;
import android.view.View;

final class Blur {
    private static final int WIDTH = 480;
    private static final int RADIUS = 2;
    private static final int PASSES = 2;

    private Blur() {
    }

    static Bitmap of(View view) {
        if (view.getWidth() == 0 || view.getHeight() == 0) {
            return null;
        }
        float scale = (float) WIDTH / view.getWidth();
        int height = Math.max(1, Math.round(view.getHeight() * scale));
        Bitmap snapshot = Bitmap.createBitmap(WIDTH, height, Bitmap.Config.ARGB_8888);
        Canvas canvas = new Canvas(snapshot);
        canvas.scale(scale, scale);
        view.draw(canvas);
        int[] pixels = new int[WIDTH * height];
        snapshot.getPixels(pixels, 0, WIDTH, 0, 0, WIDTH, height);
        for (int pass = 0; pass < PASSES; pass++) {
            pixels = boxPass(pixels, WIDTH, height, 1, 0);
            pixels = boxPass(pixels, WIDTH, height, 0, 1);
        }
        snapshot.setPixels(pixels, 0, WIDTH, 0, 0, WIDTH, height);
        return snapshot;
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
