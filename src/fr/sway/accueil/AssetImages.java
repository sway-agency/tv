package fr.sway.accueil;

import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.widget.ImageView;

import java.io.IOException;
import java.io.InputStream;
import java.util.HashMap;
import java.util.Map;

final class AssetImages {
    private static final Map<String, Bitmap> CACHE = new HashMap<>();

    private AssetImages() {
    }

    static Bitmap bitmap(Context context, String path) {
        Bitmap cached = CACHE.get(path);
        if (cached == null) {
            cached = decode(context, path, Bitmap.Config.ARGB_8888);
            CACHE.put(path, cached);
        }
        return cached;
    }

    static Bitmap opaque(Context context, String path) {
        return decode(context, path, Bitmap.Config.RGB_565);
    }

    static ImageView view(Context context, String path) {
        ImageView image = new ImageView(context);
        image.setScaleType(ImageView.ScaleType.FIT_CENTER);
        image.setImageBitmap(bitmap(context, path));
        return image;
    }

    private static Bitmap decode(Context context, String path, Bitmap.Config config) {
        BitmapFactory.Options options = new BitmapFactory.Options();
        options.inPreferredConfig = config;
        try (InputStream stream = context.getAssets().open(path)) {
            return BitmapFactory.decodeStream(stream, null, options);
        } catch (IOException e) {
            return null;
        }
    }
}
