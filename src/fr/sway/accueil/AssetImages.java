package fr.sway.accueil;

import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.widget.ImageView;

import java.io.IOException;
import java.io.InputStream;

final class AssetImages {
    private AssetImages() {
    }

    static Bitmap bitmap(Context context, String path) {
        try (InputStream stream = context.getAssets().open(path)) {
            return BitmapFactory.decodeStream(stream);
        } catch (IOException e) {
            return null;
        }
    }

    static ImageView view(Context context, String path) {
        ImageView image = new ImageView(context);
        image.setScaleType(ImageView.ScaleType.FIT_CENTER);
        image.setImageBitmap(bitmap(context, path));
        return image;
    }
}
