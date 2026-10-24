package fr.sway.accueil;

import android.content.res.AssetManager;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.drawable.BitmapDrawable;
import android.graphics.drawable.Drawable;
import android.os.Handler;
import android.os.Looper;
import android.view.animation.AccelerateDecelerateInterpolator;
import android.widget.ImageView;

import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

final class WallpaperRotator {
    private static final String FOLDER = "wallpapers";
    private static final long DISPLAY_MILLIS = 20 * 60_000;
    private static final long FADE_MILLIS = 2_500;
    private static final float DRIFT_SCALE = 1.08f;
    private static final float DRIFT_PAN = 0.025f;

    private final AssetManager assets;
    private final ImageView backLayer;
    private final ImageView frontLayer;
    private final List<String> wallpapers = new ArrayList<>();
    private final Handler handler = new Handler(Looper.getMainLooper());
    private final Runnable advance = this::showNext;
    private int nextIndex;
    private boolean frontVisible;
    private Bitmap spare;

    WallpaperRotator(AssetManager assets, ImageView backLayer, ImageView frontLayer) {
        this.assets = assets;
        this.backLayer = backLayer;
        this.frontLayer = frontLayer;
        try {
            for (String name : assets.list(FOLDER)) {
                wallpapers.add(FOLDER + "/" + name);
            }
        } catch (IOException ignored) {
        }
        Collections.shuffle(wallpapers);
        frontLayer.setAlpha(0f);
    }

    void start() {
        if (wallpapers.isEmpty()) {
            return;
        }
        if (backLayer.getDrawable() == null && frontLayer.getDrawable() == null) {
            backLayer.setImageBitmap(decode(wallpapers.get(nextIndex++ % wallpapers.size())));
            drift(backLayer);
        }
        handler.removeCallbacks(advance);
        handler.postDelayed(advance, DISPLAY_MILLIS);
    }

    void stop() {
        handler.removeCallbacks(advance);
        spare = null;
    }

    private void showNext() {
        Bitmap next = decode(wallpapers.get(nextIndex++ % wallpapers.size()));
        if (frontVisible) {
            backLayer.setImageBitmap(next);
            drift(backLayer);
            frontLayer.animate().alpha(0f).setDuration(FADE_MILLIS)
                    .withEndAction(() -> release(frontLayer)).start();
        } else {
            frontLayer.setImageBitmap(next);
            drift(frontLayer);
            frontLayer.animate().alpha(1f).setDuration(FADE_MILLIS)
                    .withEndAction(() -> release(backLayer)).start();
        }
        frontVisible = !frontVisible;
        handler.postDelayed(advance, DISPLAY_MILLIS);
    }

    private static void drift(ImageView layer) {
        int width = layer.getResources().getDisplayMetrics().widthPixels;
        int height = layer.getResources().getDisplayMetrics().heightPixels;
        double angle = Math.random() * 2 * Math.PI;
        layer.setScaleX(1f);
        layer.setScaleY(1f);
        layer.setTranslationX(0f);
        layer.setTranslationY(0f);
        layer.animate().scaleX(DRIFT_SCALE).scaleY(DRIFT_SCALE)
                .translationX((float) Math.cos(angle) * width * DRIFT_PAN)
                .translationY((float) Math.sin(angle) * height * DRIFT_PAN)
                .setInterpolator(new AccelerateDecelerateInterpolator())
                .setDuration(DISPLAY_MILLIS + FADE_MILLIS).start();
    }

    Bitmap current() {
        Drawable drawable = (frontVisible ? frontLayer : backLayer).getDrawable();
        return drawable instanceof BitmapDrawable ? ((BitmapDrawable) drawable).getBitmap() : null;
    }

    private void release(ImageView layer) {
        Drawable drawable = layer.getDrawable();
        layer.setImageDrawable(null);
        if (drawable instanceof BitmapDrawable) {
            spare = ((BitmapDrawable) drawable).getBitmap();
        }
    }

    private Bitmap decode(String path) {
        BitmapFactory.Options options = new BitmapFactory.Options();
        options.inPreferredConfig = Bitmap.Config.RGB_565;
        options.inMutable = true;
        options.inBitmap = spare;
        spare = null;
        try (InputStream stream = assets.open(path)) {
            return BitmapFactory.decodeStream(stream, null, options);
        } catch (IllegalArgumentException e) {
            options.inBitmap = null;
            return decodeFresh(path, options);
        } catch (IOException e) {
            return null;
        }
    }

    private Bitmap decodeFresh(String path, BitmapFactory.Options options) {
        try (InputStream stream = assets.open(path)) {
            return BitmapFactory.decodeStream(stream, null, options);
        } catch (IOException e) {
            return null;
        }
    }
}
