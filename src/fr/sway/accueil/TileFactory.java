package fr.sway.accueil;

import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.Color;
import android.graphics.Typeface;
import android.graphics.drawable.GradientDrawable;
import android.util.TypedValue;
import android.view.Gravity;
import android.view.View;
import android.widget.FrameLayout;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.TextView;

final class TileFactory {
    private static final float FOCUSED_SCALE = 1.12f;
    private static final long FOCUS_ANIMATION_MILLIS = 160;

    private final Context context;
    private final Typeface titleFont;
    private final Typeface labelFont;
    private final int tileWidth;
    private final int tileHeight;
    private final int spacing;

    TileFactory(Context context, int tileWidth) {
        this.context = context;
        this.titleFont = Fonts.satoshi(context.getAssets(), 400);
        this.labelFont = Fonts.satoshi(context.getAssets(), 500);
        this.tileWidth = tileWidth;
        this.tileHeight = tileWidth * 9 / 16;
        this.spacing = tileWidth / 8;
    }

    View appTile(LauncherApp app, Runnable onClick) {
        FrameLayout card = card(Color.argb(90, 255, 255, 255));
        ImageView image = new ImageView(context);
        if (app.banner != null) {
            image.setScaleType(ImageView.ScaleType.CENTER_CROP);
            image.setImageDrawable(app.banner);
            card.addView(image, new FrameLayout.LayoutParams(tileWidth, tileHeight));
        } else {
            image.setScaleType(ImageView.ScaleType.FIT_CENTER);
            image.setImageDrawable(app.icon);
            int iconSize = tileHeight / 2;
            card.addView(image, new FrameLayout.LayoutParams(iconSize, iconSize, Gravity.CENTER));
        }
        return wrap(card, app.label, onClick);
    }

    View iconTile(Bitmap icon, String title, String caption, Runnable onClick) {
        FrameLayout card = card(Color.argb(140, 18, 18, 22));
        LinearLayout content = new LinearLayout(context);
        content.setOrientation(LinearLayout.HORIZONTAL);
        content.setGravity(Gravity.CENTER);

        int iconSize = tileHeight / 4;
        ImageView image = new ImageView(context);
        image.setImageBitmap(icon);
        LinearLayout.LayoutParams iconParams = new LinearLayout.LayoutParams(iconSize, iconSize);
        iconParams.rightMargin = iconSize / 2;
        content.addView(image, iconParams);

        TextView text = new TextView(context);
        text.setText(title);
        text.setTextColor(Color.WHITE);
        text.setTypeface(titleFont);
        text.setTextSize(TypedValue.COMPLEX_UNIT_PX, tileHeight * 0.2f);
        content.addView(text);

        card.addView(content, new FrameLayout.LayoutParams(tileWidth, tileHeight));
        return wrap(card, caption, onClick);
    }

    private FrameLayout card(int fillColor) {
        GradientDrawable background = new GradientDrawable();
        background.setColor(fillColor);
        background.setCornerRadius(tileHeight * 0.12f);
        FrameLayout card = new FrameLayout(context);
        card.setBackground(background);
        card.setClipToOutline(true);
        card.setElevation(spacing / 4f);
        return card;
    }

    private View wrap(FrameLayout card, String caption, Runnable onClick) {
        LinearLayout tile = new LinearLayout(context);
        tile.setOrientation(LinearLayout.VERTICAL);
        tile.setGravity(Gravity.CENTER_HORIZONTAL);
        tile.setFocusable(true);
        tile.setClipChildren(false);
        tile.setClipToPadding(false);
        tile.setPadding(spacing, spacing, spacing, 0);
        tile.addView(card, new LinearLayout.LayoutParams(tileWidth, tileHeight));

        TextView label = new TextView(context);
        label.setText(caption);
        label.setTextColor(Color.WHITE);
        label.setTypeface(labelFont);
        label.setTextSize(TypedValue.COMPLEX_UNIT_PX, tileHeight * 0.14f);
        label.setShadowLayer(8f, 0f, 2f, Color.BLACK);
        label.setSingleLine(true);
        label.setGravity(Gravity.CENTER);
        label.setAlpha(0f);
        LinearLayout.LayoutParams labelParams = new LinearLayout.LayoutParams(tileWidth, LinearLayout.LayoutParams.WRAP_CONTENT);
        labelParams.topMargin = spacing * 3 / 2;
        tile.addView(label, labelParams);

        tile.setOnFocusChangeListener((view, focused) -> {
            float scale = focused ? FOCUSED_SCALE : 1f;
            card.setPivotY(tileHeight);
            card.animate().scaleX(scale).scaleY(scale).translationZ(focused ? spacing * 2f : 0f)
                    .setInterpolator(Motion.SPRING).setDuration(Motion.SPRING_MILLIS).start();
            label.animate().alpha(focused ? 1f : 0f).setDuration(FOCUS_ANIMATION_MILLIS).start();
        });
        tile.setOnClickListener(view -> onClick.run());
        return tile;
    }
}
